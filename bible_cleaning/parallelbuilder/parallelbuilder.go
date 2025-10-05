package parallelcorpus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"github.com/zrygan.nlp/bible_cleaning/workerprogress"
	"github.com/zrygan.nlp/bible_cleaning/config"
)

type ParallelCorpusEntry struct {
	SourceLang string            `json:"source_lang"`
	TargetLang string            `json:"target_lang"`
	Pairs      TextPairArray     `json:"pairs"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type TextPairArray []TextPair

type TextPair struct {
	SourceText string `json:"source_text"`
	TargetText string `json:"target_text"`
	ID         string `json:"id"`
}

func GetKeys(bibles map[string]map[string]string) []string {
	result := make([]string, 0, len(bibles))
	for key := range bibles {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func LoadParallelCorpus(sourcePath, targetPath, sourceLang, targetLang string) (*ParallelCorpusEntry, error) {
	corpusPath := config.SRC_PATH + "/"

	sourceFile, err := os.Open(corpusPath + sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	targetFile, err := os.Open(corpusPath + targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open target file: %w", err)
	}
	defer targetFile.Close()

	sourceScanner := bufio.NewScanner(sourceFile)
	targetScanner := bufio.NewScanner(targetFile)

	corpus := &ParallelCorpusEntry{
		SourceLang: sourceLang,
		TargetLang: targetLang,
		Pairs:      []TextPair{},
	}
	lineNum := 0
	for sourceScanner.Scan() && targetScanner.Scan() {
		lineNum++
		sourceLine := strings.TrimSpace(sourceScanner.Text())
		targetLine := strings.TrimSpace(targetScanner.Text())

		if sourceLine != "" && targetLine != "" {
			corpus.Pairs = append(corpus.Pairs, TextPair{
				SourceText: sourceLine,
				TargetText: targetLine,
				ID:         fmt.Sprintf("pair_%d", lineNum),
			})
		}
	}

	// Check for scanning errors
	if err := sourceScanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading source file: %w", err)
	}
	if err := targetScanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading target file: %w", err)
	}
	return corpus, nil
}

func (pc TextPairArray) Sort() {
	sort.Slice(pc, func(i, j int) bool {
		return pc[i].ID < pc[j].ID
	})
}

func (pc *ParallelCorpusEntry) Sort() {
	pc.Pairs.Sort()
}

// SaveAsJSON saves the corpus as JSON
func (pc *ParallelCorpusEntry) SaveAsJSON(path string) error {

	path = filepath.Join("parallel_corpus", path)

	if err := os.MkdirAll("parallel_corpus", os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(pc)
}

func (pc *ParallelCorpusEntry) SaveAsTSV(path string, outDir string) error {
	path = filepath.Join(outDir, path)

	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create TSV file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write header
	_, err = writer.WriteString("id\tsource_text\ttarget_text\n")

	if err != nil {
		return fmt.Errorf("failed to write header to TSV file: %w", err)
	}

	// Write each text pair
	for _, pair := range pc.Pairs {
		line := fmt.Sprintf("%s\t%s\t%s\n", pair.ID, transfromEscapeCharTSV(pair.SourceText), transfromEscapeCharTSV(pair.TargetText))
		_, err = writer.WriteString(line)
		if err != nil {
			return fmt.Errorf("failed to write line to TSV file: %w", err)
		}
	}

	return nil
}

// Transform special characters for TSV format
func transfromEscapeCharTSV(text string) string {
	text = strings.ReplaceAll(text, "\t", "<TAB>")
	text = strings.ReplaceAll(text, "\n", "<NEWLINE>")
	text = strings.ReplaceAll(text, "\r", "<RETURN>")
	return text
}

// Add adds a new text pair to the corpus
func (pc *ParallelCorpusEntry) Add(source, target, id string) {
	if id == "" {
		id = fmt.Sprintf("pair_%d", len(pc.Pairs))
	}
	pc.Pairs = append(pc.Pairs, TextPair{
		SourceText: source,
		TargetText: target,
		ID:         id,
	})
}

// Size returns the number of text pairs
func (pc *ParallelCorpusEntry) Size() int {
	return len(pc.Pairs)
}

// Filter filters pairs based on a predicate function
func (pc *ParallelCorpusEntry) Filter(predicate func(TextPair) bool) *ParallelCorpusEntry {
	filtered := &ParallelCorpusEntry{
		SourceLang: pc.SourceLang,
		TargetLang: pc.TargetLang,
		Metadata:   pc.Metadata,
		Pairs:      []TextPair{},
	}

	for _, pair := range pc.Pairs {
		if predicate(pair) {
			filtered.Pairs = append(filtered.Pairs, pair)
		}
	}

	return filtered
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


// worker that builds corpus for one language pair
func buildCorpusVerses(src, tgt string, index map[string]map[string]string, outdir string, prg workerprogress.WorkerProgressContext) {
	defer prg.Wg.Done()

	entry := &ParallelCorpusEntry{
		SourceLang: src,
		TargetLang: tgt,
	}
	n := 0;
	total := len(index[src])
	for verseID, srcFile := range index[src] {
		if tgtFile, ok := index[tgt][verseID]; ok {
			srcContent, _ := os.ReadFile(srcFile)
			tgtContent, _ := os.ReadFile(tgtFile)

			entry.Pairs = append(entry.Pairs, TextPair{
				SourceText: strings.TrimSpace(string(srcContent)),
				TargetText: strings.TrimSpace(string(tgtContent)),
				ID:         verseID,
			})
		}
		n = n + 1;
		if (n % 50 == 0)  {
			prg.Progress <- workerprogress.WorkerProgressMsg{
				WorkerID: prg.WorkerID,
				Percent:  float32(len(entry.Pairs)) / float32(len(index[src])),
				Status:   fmt.Sprintf("Processed %03d/%03d verses for %s <--> %s", n, total, src, tgt),
			}
		}
	}

	entry.Sort()

	prg.Progress <- workerprogress.WorkerProgressMsg {
		WorkerID: prg.WorkerID,
		Percent: 1.0,
		Status: fmt.Sprintf("Built sentence-level corpus for %s <--> %s (%03d pairs); Saving TSV file. ", src, tgt, len(entry.Pairs)),
	}

	entry.SaveAsTSV(fmt.Sprintf("%s_%s.tsv", src, tgt), outdir)
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

func GenerateParallelCorpusByVerses() error {
	root := config.CORPUS_VERSES_FOLDER

	// Make sure destination directory exists
	if err := os.MkdirAll(config.PARALLEL_VERSES_FOLDER, os.ModePerm); err != nil {
		return err
	}

	index, err := indexLanguageFileMap(root)

	if err != nil {
		return err
	}

	// Get a list of unique languages
	langs := make([]string, 0, len(index))
	for k := range index {
		langs = append(langs, k)
	}

	sort.Strings(langs)


	println(fmt.Sprintf("Found %d languages, generating parallel corpora...", len(langs)))
	// n choose 2

	// launch workers for each unique pair of languages
	total := int32(len(langs) * (len(langs) - 1) / 2)	
	workers := atomic.Int32{}
	start := time.Now()

	var queenWg sync.WaitGroup
	progressCh := make(chan workerprogress.WorkerProgressMsg, 100)

	queenCtx := workerprogress.QueenContext{
			DoneWorkers:   &workers,
			TotalWorkers:  int(total),
			StartTime:     start,
			Quit:          make(chan struct{}),
			Wg:            &queenWg,
			ProgressCh:    progressCh,
	}

	jobCh := make(chan [2]string)
	for i := 0; i < len(langs); i++ {
		for j := i + 1; j < len(langs); j++ {
			jobCh <- [2]string{langs[i], langs[j],}
		}
	}

	queenWg.Add(1)
	go queenCtx.RunReporter()
	numOfCores := 8;
	for i := 0 ; i < numOfCores; i++ {
		queenCtx.Wg.Add(1)
		go func(id int) {
			defer queenCtx.Wg.Done()
			defer fmt.Printf(fmt.Sprintf("Worker thread %d exiting...\n", id))

			for pair := range jobCh {
			
				workerID := fmt.Sprintf("%s-%s", pair[0], pair[1])

				workerCtx := workerprogress.WorkerProgressContext{
						Wg:       queenCtx.Wg,
						Progress: progressCh,
						WorkerID: workerID,
				}

				buildCorpusVerses(pair[0], pair[1], index, config.PARALLEL_VERSES_FOLDER, workerCtx) // n choose 2
			}
		}(i)
	}
	

	queenWg.Wait()
	close(queenCtx.Quit) // signal exiting
	
	// Stop reporter goroutine
	close(progressCh)   // close progress updates
	elapsed := time.Since(start)
	println(fmt.Sprintf("All done in %s!", &elapsed))
	return nil
}

func nChoose2(n int) int {
	return n * (n - 1) / 2
}

func getParallelQueenConfig() *workerprogress.QueenConfig {
	return &workerprogress.QueenConfig{
		IsDetailed:     false,
		ReportInterval: 1000, // milliseconds
		UseProgressBar: false,
	}
}

func initializePCBySentences(root string) (map[string]map[string]string, []string, error) {

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


// worker that builds corpus for one language pair at the sentence level
func buildCorpusSentences(
	src, tgt string,
	index map[string]map[string]string, // verseID -> filepath per language
	outdir string,
	prg workerprogress.WorkerProgressContext,
) {
	defer prg.Wg.Done()

	entry := &ParallelCorpusEntry{
		SourceLang: src,
		TargetLang: tgt,
	}

	for verseID, srcFile := range index[src] {

		tgtFile, ok := index[tgt][verseID]

		if !ok {
			continue
		}

		srcLines, err := readLines(srcFile)
		if err != nil {
			continue
		}

		tgtLines, err := readLines(tgtFile)
		if err != nil {
			continue
		}

		// align by line number
		n := min(len(srcLines), len(tgtLines))

		for i := 0; i < n; i++ {
			entry.Pairs = append(entry.Pairs, TextPair{
				SourceText: strings.TrimSpace(srcLines[i]),
				TargetText: strings.TrimSpace(tgtLines[i]),
				ID:         fmt.Sprintf("%s_%04d", verseID, i+1),
			})

			if i % 50 == 0 {
				prg.Progress <- workerprogress.WorkerProgressMsg{
					WorkerID: prg.WorkerID,
					Percent:  float32(i) / float32(n),
					Status:   fmt.Sprintf("Processed %03d/%03d lines for %s <--> %s", i, n, src, tgt),
				}
			}
		}
	}

	entry.Sort()

	prg.Progress <- workerprogress.WorkerProgressMsg{
		WorkerID: prg.WorkerID,
		Percent: 1.0,
		Status: fmt.Sprintf("Built sentence-level corpus for %s <--> %s (%03d pairs); Saving TSV file. ", src, tgt, len(entry.Pairs)),
	}
	
	entry.SaveAsTSV(fmt.Sprintf("%s_%s.tsv", src, tgt), outdir)
}

func GenerateParallelCorpusBySentences() error {
	index, langs, err := initializePCBySentences(config.CORPUS_SENTENCES_FOLDER)

	if err != nil {
		return err
	}

	queenCtx := workerprogress.NewQueenContext(nChoose2(len(langs)), getParallelQueenConfig())

	go queenCtx.RunReporter()

	// Launch workers
	for i := 0; i < len(langs); i++ {
		for j := i + 1; j < len(langs); j++ {
			workerCtx := queenCtx.CreateWorkerContext(fmt.Sprintf("%s-%s", langs[i], langs[j]))
			go buildCorpusSentences(langs[i], langs[j], index, config.PARALLEL_SENTENCES_FOLDER, workerCtx)
		}
	}

	queenCtx.Wg.Wait()

	close(queenCtx.Quit) 
	close(queenCtx.ProgressCh)   

	elapsed := time.Since(queenCtx.StartTime)

	fmt.Printf("\nAll done in %s!\n", elapsed.Truncate(time.Millisecond))
	return nil
}
