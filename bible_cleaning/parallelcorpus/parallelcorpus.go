package parallelcorpus

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
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

func generateUniquePairs(bibles map[string]string, f func(string, string) *ParallelCorpusEntry) []*ParallelCorpusEntry {
	n := len(bibles)
	result := []*ParallelCorpusEntry{}
	// Get language languageKeys
	languageKeys := make([]string, 0, n)
	for k := range bibles {
		languageKeys = append(languageKeys, k)
	}

	for i := 0 ; i < n; i++ {
		for j := i + 1; j < n; j++ {
				result = append(result, f(languageKeys[i], languageKeys[j]))
		}
	}

	return result
}

func GenerateParallelCorpus(bibles map[string]string) {
	
	
	appliedFilter := func(lang1, lang2 string) *ParallelCorpusEntry {
		
		return nil
	}

	generateUniquePairs(bibles, appliedFilter) 
}