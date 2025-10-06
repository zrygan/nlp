package sentencecleaning

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// extractSentences splits the verse text into sentences
func extractSentencesWithVerse(text, verseID string) []string {
	reSentence := regexp.MustCompile(`([.?!])\s+`)
	reNormalize := regexp.MustCompile(`\s+`)

	text = reNormalize.ReplaceAllString(text, " ")
	matches := reSentence.FindAllStringIndex(text, -1)

	var sentences []string
	lastEnd := 0

	for _, match := range matches {
		end := match[1]
		s := strings.TrimSpace(text[lastEnd:end])
		if s != "" {
			sentences = append(sentences, fmt.Sprintf("%s\t%s", verseID, s))
		}
		lastEnd = end
	}

	// remaining tail
	remaining := strings.TrimSpace(text[lastEnd:])
	if remaining != "" {
		sentences = append(sentences, fmt.Sprintf("%s\t%s", verseID, remaining))
	}

	return sentences
}


func mergeOrAppend(sentences []string, remaining string) []string {
	lastSentenceIndex := len(sentences) - 1
	if lastSentenceIndex >= 0 && !strings.HasSuffix(sentences[lastSentenceIndex], ".") &&
		!strings.HasSuffix(sentences[lastSentenceIndex], "?") &&
		!strings.HasSuffix(sentences[lastSentenceIndex], "!") {
		sentences[lastSentenceIndex] += " " + remaining
	} else {
		sentences = append(sentences, remaining)
	}
	return sentences
}

// SplitCorpusBySentence walks through files and processes each one
func SplitCorpusBySentence(root, outRoot string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		return processFile(path, root, outRoot)
	})
}
func processFile(path, root, outRoot string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fmt.Println("Processing: ", path)

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	parts := strings.Split(base, "_")
	if len(parts) < 4 {
		return fmt.Errorf("unexpected filename format: \"%s\"", base)
	}

	
	var sentences []string

	for i, verseText := range lines {
		verseNum := fmt.Sprintf("%03d", i+1)
		verseID := fmt.Sprintf("%s", verseNum)
		verseSentences := extractSentencesWithVerse(verseText, verseID)
		sentences = append(sentences, verseSentences...)
	}

	outPath, err := makeOutputPath(path, root, outRoot)
	if err != nil {
		return err
	}

	fmt.Printf("Extracted %d sentences from %s → %s\n", len(sentences), path, outPath)
	return writeSentences(outPath, sentences)
}


// writeSentences ensures the folder exists and saves sentences line by line. 
func writeSentences(outPath string, sentences []string) error { 
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil { 
		return err 
	} 
	f, err := os.Create(outPath) 
	if err != nil { 
		return err 
	} 
	defer f.Close() 
	println("Writing to:", outPath) 

	_, err = f.WriteString("verse\tcontent\n") 
	if err != nil { 
		return fmt.Errorf("failed to write header to TSV file: %w", err) 
	} 

	for _, s := range sentences { 
		if _, err := f.WriteString(s + "\n"); err != nil { return err } 
	} 
	return nil 
}


func makeOutputPath(path, root, outRoot string) (string, error) {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return "", err
	}
	return filepath.Join(outRoot, relPath), nil
}
