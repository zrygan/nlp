package parallelcorpus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"strings"
)

type ParallelCorpusEntry struct {
	SourceLang string   `json:"source_lang"`
	TargetLang string   `json:"target_lang"`
	Pairs      []TextPair       `json:"pairs"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type TextPair struct {
	SourceText string `json:"source_text"`
	TargetText string `json:"target_text"`
	ID 			 string   `json:"id"`
}

func GetKeys(bibles map[string]map[string]string) []string{
	result := make([]string, 0, len(bibles))
	for key := range bibles {
		result = append(result, key)
	}
	return result
}

func LoadParallelCorpus(sourcePath, targetPath, sourceLang, targetLang string) (*ParallelCorpusEntry, error) {
	corpusPath := "corpus/"

	sourceFile, err := os.Open(corpusPath+sourcePath)
	if err != nil {
			return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	targetFile, err := os.Open(corpusPath+targetPath)
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
						ID:     fmt.Sprintf("pair_%d", lineNum),
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

// SaveAsJSON saves the corpus as JSON
func (pc *ParallelCorpusEntry) SaveAsJSON(path string) error {

		path = filepath.Join("parallel_corpus", path);

		if err := os.MkdirAll("parallel_corpus", os.ModePerm); err != nil {
			return err
		}

    file, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(pc)
}

// Add adds a new text pair to the corpus
func (pc *ParallelCorpusEntry) Add(source, target, id string) {
    if id == "" {
        id = fmt.Sprintf("pair_%d", len(pc.Pairs))
    }
    pc.Pairs = append(pc.Pairs, TextPair{
        SourceText: source,
        TargetText: target,
        ID:     id,
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


// fileIndex: lang -> verseID -> filepath
func indexFiles(root string) (map[string]map[string]string, error) {
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
func buildCorpus(src, tgt string, index map[string]map[string]string, wg *sync.WaitGroup, out chan<- *ParallelCorpusEntry) {
	defer wg.Done()

	entry := &ParallelCorpusEntry{
		SourceLang: src,
		TargetLang: tgt,
	}

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
	}
	out <- entry
}

func GenerateParallelCorpus() ([]*ParallelCorpusEntry, error) {
	root := "corpus"
	
	index, err := indexFiles(root)
	if err != nil {
		return nil, err
	}

	langs := make([]string, 0, len(index))
	for k := range index {
		langs = append(langs, k)
	}

	var wg sync.WaitGroup
	out := make(chan *ParallelCorpusEntry)

	// launch workers for each unique pair
	for i := 0; i < len(langs); i++ {
		for j := i + 1; j < len(langs); j++ {
			wg.Add(1)
			go buildCorpus(langs[i], langs[j], index, &wg, out)
		}
	}

	// closer
	go func() {
		wg.Wait()
		close(out)
	}()

	// collect
	var results []*ParallelCorpusEntry
	for entry := range out {
		results = append(results, entry)
	}

	return results, nil
}
