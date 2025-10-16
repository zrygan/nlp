package types

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"github.com/zrygan.nlp/bible_cleaning/config"
)

type FindReplaceTuple[T string | *regexp.Regexp] struct {
	Find    T
	Replace T
}

func TurnToRegexpsTuple(tuples []FindReplaceTuple[string]) []FindReplaceTuple[*regexp.Regexp] {
	var result []FindReplaceTuple[*regexp.Regexp]

	for _, t := range tuples {
		result = append(result, FindReplaceTuple[*regexp.Regexp]{
			Find:    regexp.MustCompile(t.Find),
			Replace: regexp.MustCompile(t.Replace),
		})
	}

	return result
}
type VerseRef struct {
	book string
	chapter string
	verse string
}
type LanguageClass struct {
	Language string
	OutputDir string
}


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
	Book 	   string `json:"book,omitempty"`
	Chapter    string `json:"chapter,omitempty"`
	Verse 	   string `json:"verse,omitempty"`
	Sentence   string `json:"sentence,omitempty"`
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
	_, err = writer.WriteString("id\tbook\tchapter\tverse\tsource_text\ttarget_text\n")

	if err != nil {
		return fmt.Errorf("failed to write header to TSV file: %w", err)
	}

	// Write each text pair
	for _, pair := range pc.Pairs {
		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\n", pair.ID, pair.Book, pair.Chapter, pair.Verse, TransfromEscapeCharTSV(pair.SourceText), TransfromEscapeCharTSV(pair.TargetText))
		_, err = writer.WriteString(line)
		if err != nil {
			return fmt.Errorf("failed to write line to TSV file: %w", err)
		}
	}

	return nil
}

func (pc *ParallelCorpusEntry) SaveAsTSVSentences(path string, outDir string) error {
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
	_, err = writer.WriteString("id\tbook\tchapter\tsentence_no\tsource_text\ttarget_text\n")

	if err != nil {
		return fmt.Errorf("failed to write header to TSV file: %w", err)
	}

	// Write each text pair
	for _, pair := range pc.Pairs {
		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\n", pair.ID, pair.Book, pair.Chapter, pair.Sentence, TransfromEscapeCharTSV(pair.SourceText), TransfromEscapeCharTSV(pair.TargetText))
		_, err = writer.WriteString(line)
		if err != nil {
			return fmt.Errorf("failed to write line to TSV file: %w", err)
		}
	}

	return nil
}

// Transform special characters for TSV format
func TransfromEscapeCharTSV(text string) string {
	text = strings.ReplaceAll(text, "\t", config.TOKEN_TAB)
	text = strings.ReplaceAll(text, "\n", config.TOKEN_NEWLINE)
	text = strings.ReplaceAll(text, "\r", config.TOKEN_RETURN)
	return text
}

func RemoveEscapeCharTSV(text string) string {
	text = strings.ReplaceAll(text, config.TOKEN_TAB, "\t")
	text = strings.ReplaceAll(text, config.TOKEN_NEWLINE, "\n")
	text = strings.ReplaceAll(text, config.TOKEN_RETURN, "\r")
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

type ProperNounCache struct {
	Words map[string]struct{}
}

func ExtractProperNouns(sentences []string) *ProperNounCache {
	cache := &ProperNounCache{Words: make(map[string]struct{})}
	for _, s := range sentences {
		tokens := strings.Fields (s)
		for i, token := range tokens {
			clean := strings.Trim(token, ".,;:!?\"'")
			if len(clean) == 0 {
				continue
			}

			if i > 0 && unicode.IsUpper(rune(clean[0])) {
				cache.Words[clean] = struct{}{}
			}
		}
	}
	return cache
}