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
	"github.com/zrygan.nlp/bible_cleaning/sentencealignment"
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

func sharedKeys(a, b map[string][]string) []string {
	var shared []string
	for k := range a {
		if _, ok := b[k]; ok {
			shared = append(shared, k)
		}
	}
	sort.Strings(shared)
	return shared
}

func readVerseMap(path string) (map[string][]string, error) {
	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}

	if len(lines) > 0 && strings.HasPrefix(lines[0], "verse") {
		lines = lines[1:] // remove header
	}

	verses := make(map[string][]string)

	for _, line := range lines {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		verseID := strings.TrimSpace(parts[0])
		content := strings.TrimSpace(types.RemoveEscapeCharTSV(parts[1]))

		if content == "" {
			continue
		}

		verses[verseID] = append(verses[verseID], content)
	}

	return verses, nil
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
			book := strings.SplitN(verseID, "_", 2)[0]
			chapter := strings.SplitN(verseID, "_", 2)[1]

			srcContent, _ := os.ReadFile(srcFile)
			tgtContent, _ := os.ReadFile(tgtFile)

			srcLines := strings.Split(strings.TrimSpace(string(srcContent)), "\n")
			tgtLines := strings.Split(strings.TrimSpace(string(tgtContent)), "\n")

			for v := 0; v < len(srcLines) && v < len(tgtLines); v++ {
				verseNum := fmt.Sprintf("%03d", v+1)
				entry.Pairs = append(entry.Pairs, types.TextPair{
					SourceText: strings.TrimSpace(srcLines[v]),
					TargetText: strings.TrimSpace(tgtLines[v]),
					ID:         fmt.Sprintf("%s_%s_%s", book, chapter, verseNum),
					Book:       book,
					Chapter:    chapter,
					Verse:      verseNum,
				})
			}
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
// buildCorpusSentences aligns verse-level TSVs (verse\tcontent) between src and tgt languages.
// It performs safe sentence alignment per verse and accounts for missing or uneven sentence counts.

func buildCorpusSentences(
	src, tgt string,
	index map[string]map[string]string, // chapterName -> filepath per language
	outdir string,
	prg workerprogress.WorkerProgressContext,
) {
	entry := &types.ParallelCorpusEntry{
		SourceLang: src,
		TargetLang: tgt,
	}

	for chapterName, srcFile := range index[src] {

		// Find corresponding target file (same chapter)
		tgtFile, ok := index[tgt][chapterName]
		if !ok {
			fmt.Printf("Missing chapter %s in %s\n", chapterName, tgt)
			continue
		}

		srcVerses, err := readVerseMap(srcFile)
		if err != nil {
			fmt.Printf("Skipping chapter %s (%s): failed to read src: %v\n", chapterName, srcFile, err)
			continue
		}
		tgtVerses, err := readVerseMap(tgtFile)
		if err != nil {
			fmt.Printf("Skipping chapter %s (%s): failed to read tgt: %v\n", chapterName, tgtFile, err)
			continue
		}

		shared := sharedKeys(srcVerses, tgtVerses)

		for _, verseID := range shared {
			srcSentences := srcVerses[verseID]
			tgtSentences := tgtVerses[verseID]

			if len(srcSentences) == 0 || len(tgtSentences) == 0 {
				continue
			}

			pairs := sentencealignment.AlignSentencesByGaleChurchDP(srcSentences, tgtSentences, verseID)
			
			chapterParts := strings.SplitN(chapterName, "_", 2)
			if len(chapterParts) != 2 {
				fmt.Printf("Invalid chapter name format: %s\n", chapterName)
				continue
			}
			book := chapterParts[0]
			chapter := chapterParts[1]

			for i := range pairs {
				sentenceNum := fmt.Sprintf("%d", i+1)
				pairs[i].ID = fmt.Sprintf("%s_%s_%s", book, chapter, sentenceNum)
				pairs[i].Book = book
				pairs[i].Chapter = chapter
				pairs[i].Sentence = sentenceNum
			}
			
			entry.Pairs = append(entry.Pairs, pairs...)

			prg.Progress <- workerprogress.WorkerProgressMsg{
				WorkerID: prg.WorkerID,
				Status:   fmt.Sprintf("Aligned %s (%d pairs)", verseID, len(pairs)),
			}
		}
	}

	entry.Sort()

	outPath := fmt.Sprintf("%s_%s.tsv", src, tgt)
	fmt.Printf("Saving aligned corpus: %s/%s\n", outdir, outPath)
	entry.SaveAsTSVSentences(outPath, outdir)

	prg.Progress <- workerprogress.WorkerProgressMsg{
		WorkerID: prg.WorkerID,
		Percent:  1.0,
		Status:   fmt.Sprintf("Finished alignment for %s <--> %s (%03d pairs)", src, tgt, len(entry.Pairs)),
	}
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
