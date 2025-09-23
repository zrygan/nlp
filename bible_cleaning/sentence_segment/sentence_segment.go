package sentence_segment

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Naive regex to split sentences
var sentenceSplitter = regexp.MustCompile(`([^.!?]+[.!?])`)

// SplitSentences splits a verse into individual sentences
func SplitSentences(text string) []string {
	matches := sentenceSplitter.FindAllString(text, -1)
	var sentences []string
	for _, s := range matches {
		cleaned := strings.TrimSpace(s)
		if cleaned != "" {
			sentences = append(sentences, cleaned)
		}
	}
	return sentences
}

// processCorpus reads all files in corpus/ and segments verses into sentences
func processCorpus(inputDir, outputDir string) error {
	return filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".txt" {
			rel, _ := filepath.Rel(inputDir, path)
			outPath := filepath.Join(outputDir, rel)

			// Ensure output folder exists
			if err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
				return err
			}

			// Open input file
			inFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer inFile.Close()

			scanner := bufio.NewScanner(inFile)

			// Open output file
			outFile, err := os.Create(outPath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			writer := bufio.NewWriter(outFile)

			for scanner.Scan() {
				verse := scanner.Text()
				sentences := SplitSentences(verse)

				for _, s := range sentences {
					_, _ = writer.WriteString(s + "\n")
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}

			writer.Flush()
			fmt.Println("Processed:", outPath)
		}
		return nil
	})
}
