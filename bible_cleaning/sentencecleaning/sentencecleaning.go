package sentencecleaning
import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"


)

func extractSentences(text string) []string {
	reSentence := regexp.MustCompile(`(G|Gng|Bb)\.|\s*([\?!.])\s*`)
	reNormalize := regexp.MustCompile(`\s+`)

	text = reNormalize.ReplaceAllString(text, " ")
	matches := reSentence.FindAllStringSubmatchIndex(text, -1)

	var sentences []string
	lastEnd := 0

	for _, match := range matches {
		fullEnd := match[1]

		switch {
		case match[4] != -1:
			sentence := strings.TrimSpace(text[lastEnd:match[5]])
			if sentence != "" {
				sentences = append(sentences, sentence)
			}
			lastEnd = fullEnd

		case match[2] != -1:
			lastEnd = fullEnd
		}
	}

	remaining := strings.TrimSpace(text[lastEnd:])
	if remaining != "" {
		sentences = mergeOrAppend(sentences, remaining)
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

// splitSentences splits the existing corpus into sentences
func SplitCorpusBySentence(root string, outRoot string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		return processFile(path, root, outRoot)
	})
}

// processFile reads a file, splits into sentences, and writes to output location.
func processFile(path, root, outRoot string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	sentences := extractSentences(string(data))

	outPath, err := makeOutputPath(path, root, outRoot)
	if err != nil {
		return err
	}

	return writeSentences(outPath, sentences)
}


// makeOutputPath maps input file path -> output file path.
func makeOutputPath(path, root, outRoot string) (string, error) {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return "", err
	}
	return filepath.Join(outRoot, relPath), nil
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
	for _, s := range sentences {
		if _, err := f.WriteString(s + "\n"); err != nil {
			return err
		}
	}
	return nil
}
