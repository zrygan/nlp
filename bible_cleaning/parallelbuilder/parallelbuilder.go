package parallelcorpus

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/zrygan.nlp/bible_cleaning/config"
	"github.com/zrygan.nlp/bible_cleaning/types"
	"github.com/zrygan.nlp/bible_cleaning/workerprogress"
)

func nChoose2(n int) int {
	return n * (n - 1) / 2
}

func getParallelQueenConfig() *workerprogress.QueenConfig {
	return &workerprogress.QueenConfig{
		IsDetailed:     config.IS_DETAILED,
		ReportInterval: config.WORKER_REPORT_INTERVAL_MS, // milliseconds
		UseProgressBar: config.USE_PROGRESS_BAR,
	}
}

func GetKeys(bibles map[string]map[string]string) []string {
	result := make([]string, 0, len(bibles))
	for key := range bibles {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

// readLines reads a file into a slice of strings (one per line).
func readLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out, nil
}

func indexLanguageFileMap(root string) (map[string]map[string]string, error) {
	files, err := filepath.Glob(root + "/*/*.txt")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^([a-z]+)_([A-Z]+)_[^_]+_(\d+)\.txt$`)
	index := make(map[string]map[string]string)

	for _, file := range files {
		base := filepath.Base(file)
		lang := filepath.Base(filepath.Dir(file))

		matches := re.FindStringSubmatch(base)
		if matches == nil {
			continue
		}
		book := matches[2]
		verse := matches[3]
		verseID := book + "_" + verse

		if _, ok := index[lang]; !ok {
			index[lang] = make(map[string]string)
		}
		index[lang][verseID] = file
	}

	return index, nil
}

/*
Builds a channel of jobs for each unique language pair (n choose 2).
*/
func buildLanguagePairJobs(languageKeys []string) chan [2]string {
	jobCh := make(chan [2]string, 200)
	for i := 0; i < len(languageKeys); i++ {
		for j := i + 1; j < len(languageKeys); j++ {
			fmt.Printf("Queued job for language pair %s-%s...\n", languageKeys[i], languageKeys[j])
			jobCh <- [2]string{languageKeys[i], languageKeys[j]}
		}
	}
	close(jobCh)
	return jobCh
}

/*
	Creates a thread pool to process language pairs in parallel.
*/
func createLanguagePairThreadPool(numOfThreads int, jobCh chan [2]string, queenCtx workerprogress.QueenContext, workerFunc func(string, string, workerprogress.WorkerProgressContext)) {
	for i := 0; i < numOfThreads; i++ {
		queenCtx.Wg.Add(1)

		go func(id int) {
			defer queenCtx.Wg.Done()
			defer fmt.Printf("%s", fmt.Sprintf("Worker thread %d exiting...\n", id))

			fmt.Printf("%s", fmt.Sprintf("Worker thread %d starting...\n", id))
			for pair := range jobCh {

				languagePair := fmt.Sprintf("%s-%s", pair[0], pair[1])
				fmt.Printf("Processing language pair %s...\n", languagePair)
				workerCtx := workerprogress.WorkerProgressContext{
					Wg:       queenCtx.Wg,
					Progress: queenCtx.ProgressCh,
					WorkerID: languagePair,
				}

				workerFunc(pair[0], pair[1], workerCtx) // n choose 2
			}

		}(i)
	}

}

/*
	Waits for all workers to finish, then closes the quit channel to signal the reporter to stop.
*/
func closeoutThreadPool(queenCtx *workerprogress.QueenContext) {
	queenCtx.Wg.Wait()
	close(queenCtx.Quit) // signal exiting

	// Stop reporter goroutine
	close(queenCtx.ProgressCh) // close progress updates
	elapsed := time.Since(queenCtx.StartTime)
	fmt.Printf("All done in %s!\n", &elapsed)
}

/*
Reads the source and target lines for a given verseID from the index.
Errors if the required file is not found.
*/
func readSrcTgtLines(srcFile string, index *map[string]map[string]string, tgt string, verseID string) ([]string, []string, error) {
	tgtFile, ok := (*index)[tgt][verseID]

	if !ok {
		return nil, nil, fmt.Errorf("target file not found for %s in %s", verseID, tgt)
	}

	srcLines, err := readLines(srcFile)
	if err != nil {
		return nil, nil, err
	}

	tgtLines, err := readLines(tgtFile)
	if err != nil {
		return nil, nil, err
	}

	return srcLines, tgtLines, nil
}

func findVerseAlignments(srcIndex, tgtIndex map[string]string) (shared, missingSrc, missingTgt []string) {
	// 1. Check verses from src perspective
	for verseID := range srcIndex {
		if _, ok := tgtIndex[verseID]; ok {
			shared = append(shared, verseID)
		} else {
			missingTgt = append(missingTgt, verseID)
		}
	}

	// 2. Check verses from tgt perspective
	for verseID := range tgtIndex {
		if _, ok := srcIndex[verseID]; !ok {
			missingSrc = append(missingSrc, verseID)
		}
	}

	return shared, missingSrc, missingTgt
}

/*

	# Verse-level parallel corpus generation

*/

/*
Given a source and target language, builds a parallel corpus by aligning verses by verseID.
*/
func buildCorpusVerses(src, tgt string, index map[string]map[string]string, outdir string, prg workerprogress.WorkerProgressContext) {
	entry := &types.ParallelCorpusEntry{
		SourceLang: src,
		TargetLang: tgt,
	}
	n := 0
	total := len(index[src])

	for verseID, srcFile := range index[src] {
		if tgtFile, ok := index[tgt][verseID]; ok {
			srcContent, _ := os.ReadFile(srcFile)
			tgtContent, _ := os.ReadFile(tgtFile)

			entry.Pairs = append(entry.Pairs, types.TextPair{
				SourceText: strings.TrimSpace(string(srcContent)),
				TargetText: strings.TrimSpace(string(tgtContent)),
				ID:         verseID,
			})
		}
		n = n + 1
		if n % config.WORKER_THREAD_REPORT_PROGRESS_RATE == 0 {
			prg.Progress <- workerprogress.WorkerProgressMsg{
				WorkerID: prg.WorkerID,
				Percent:  float32(len(entry.Pairs)) / float32(len(index[src])),
				Status:   fmt.Sprintf("Processed %03d/%03d verses for %s <--> %s", n, total, src, tgt),
			}
		}
	}
	
	entry.Sort()

	prg.Progress <- workerprogress.WorkerProgressMsg{
		WorkerID: prg.WorkerID,
		Percent:  1.0,
		Status:   fmt.Sprintf("Built sentence-level corpus for %s <--> %s (%03d pairs); Saving TSV file. ", src, tgt, len(entry.Pairs)),
	}

	entry.SaveAsTSV(fmt.Sprintf("%s_%s.tsv", src, tgt), outdir)
}

/*
Wrapper to pass additional parameters to the worker function.
Mainly used for createLanguagePairThreadPool.
*/
func buildCorpusVersesWrapper(index map[string]map[string]string, outdir string) func(string, string, workerprogress.WorkerProgressContext) {
	return func(src, tgt string, prg workerprogress.WorkerProgressContext) {
		buildCorpusVerses(src, tgt, index, outdir, prg)
	}
}

/*
Initializes the verse-level parallel corpus generation by indexing the files and creating the output directory.
*/
func initializeParallelCorpusByVerses() (map[string]map[string]string, []string, error) {
	root := config.CORPUS_VERSES_FOLDER

	// Make sure destination directory exists
	if err := os.MkdirAll(config.PARALLEL_VERSES_FOLDER, os.ModePerm); err != nil {
		return nil, nil, err
	}

	index, err := indexLanguageFileMap(root)

	if err != nil {
		return nil, nil, err
	}

	// Get a list of unique languages
	langs := make([]string, 0, len(index))
	for k := range index {
		langs = append(langs, k)
	}

	sort.Strings(langs)
	return index, langs, nil
}

/*
Generates the parallel corpus by verses for all language pairs found in the corpus/verses folder.
It creates a thread pool to process multiple language pairs in parallel.
*/
func GenerateParallelCorpusByVerses() error {
	index, langs, err := initializeParallelCorpusByVerses()

	if err != nil {
		return err
	}

	println(fmt.Sprintf("Found %d languages, generating parallel corpora...", len(langs)))
	// launch workers for each unique pair of languages
	total := nChoose2(len(langs))
	queenCtx := workerprogress.NewQueenContext(total, getParallelQueenConfig())
	jobCh := buildLanguagePairJobs(langs)
	fmt.Printf("Created %d jobs for %d languages.\n", len(jobCh), len(langs))

	go queenCtx.RunReporter()
	createLanguagePairThreadPool(config.THREAD_POOL_SIZE, jobCh, *queenCtx, buildCorpusVersesWrapper(index, config.PARALLEL_VERSES_FOLDER))
	closeoutThreadPool(queenCtx)
	return nil
}

/**
# Sentence-level parallel corpus generation.
*/

func initializeParallelCorpusBySentences(root string) (map[string]map[string]string, []string, error) {

	if err := os.MkdirAll(config.PARALLEL_SENTENCES_FOLDER, os.ModePerm); err != nil {
		return nil, nil, err
	}

	index, err := indexLanguageFileMap(root)
	if err != nil {
		return nil, nil, err
	}

	langs := GetKeys(index)

	fmt.Printf("Found %d languages for sentence-level corpus generation.\n", len(langs))
	return index, langs, nil
}

/*
Given a source and target language, builds a parallel corpus by aligning sentences by line number.
*/
func buildCorpusSentences(
	src, tgt string,
	index map[string]map[string]string, // verseID -> filepath per language
	outdir string,
	prg workerprogress.WorkerProgressContext,
) {
	entry := &types.ParallelCorpusEntry{
		SourceLang: src,
		TargetLang: tgt,
	}

	for verseID, srcFile := range index[src] {

		srcLines, tgtLines, err := readSrcTgtLines(srcFile, &index, tgt, verseID)

		if err != nil {
			continue
		}

		n := max(len(srcLines), len(tgtLines))


		for i := 0; i < n; i++ {
		
			srcLine := config.TOKEN_MISSING_TRANSLATION
			if i < len(srcLines) {
				srcLine = srcLines[i]
			}

			tgtLine := config.TOKEN_MISSING_TRANSLATION
			if i < len(tgtLines) {
				tgtLine = tgtLines[i]
			}
			strings.Split(srcLine, "\t")
			strings.Split(tgtLine, "\t")


			entry.Pairs = append(entry.Pairs, types.TextPair{
				SourceText: strings.TrimSpace(srcLine),
				TargetText: strings.TrimSpace(tgtLine),
				ID:         fmt.Sprintf("%s_%04d", verseID, i+1),
			})

			if i % config.WORKER_THREAD_REPORT_PROGRESS_RATE != 0 {
				continue
			}

			prg.Progress <- workerprogress.WorkerProgressMsg{
				WorkerID: prg.WorkerID,
				Percent:  float32(i) / float32(n),
				Status:   fmt.Sprintf("Processed %03d/%03d lines for %s <--> %s", i, n, src, tgt),
			}
		}
	}

	entry.Sort()

	prg.Progress <- workerprogress.WorkerProgressMsg{
		WorkerID: prg.WorkerID,
		Percent:  1.0,
		Status:   fmt.Sprintf("Built sentence-level corpus for %s <--> %s (%03d pairs); Saving TSV file. ", src, tgt, len(entry.Pairs)),
	}
	fmt.Printf("Saving TSV file %s\n", outdir)
	entry.SaveAsTSV(fmt.Sprintf("%s_%s.tsv", src, tgt), outdir)
}

/*
Wrapper to pass additional parameters to the worker function.
Mainly used for createLanguagePairThreadPool.
*/
func buildCorpusSentencesWrapper(index map[string]map[string]string, outdir string) func(string, string, workerprogress.WorkerProgressContext) {
	return func(src, tgt string, prg workerprogress.WorkerProgressContext) {
		buildCorpusSentences(src, tgt, index, outdir, prg)
	}
}

/*
Generates the parallel corpus by sentences for all language pairs found in the corpus/verses folder.
It creates a thread pool to process multiple language pairs in parallel.
*/
func GenerateParallelCorpusBySentences() error {
	root := config.CORPUS_SENTENCES_FOLDER

	index, langs, err := initializeParallelCorpusBySentences(root)

	if err != nil {
		return err
	}

	queenCtx := workerprogress.NewQueenContext(nChoose2(len(langs)), getParallelQueenConfig())

	fmt.Printf("Found %d languages, generating parallel corpora...\n", len(langs))
	jobCh := buildLanguagePairJobs(langs)
	fmt.Printf("Created %d jobs for %d languages.\n", len(jobCh), len(langs))

	go queenCtx.RunReporter()

	createLanguagePairThreadPool(config.THREAD_POOL_SIZE, jobCh, *queenCtx, buildCorpusSentencesWrapper(index, config.PARALLEL_SENTENCES_FOLDER))

	closeoutThreadPool(queenCtx)
	return nil
}
