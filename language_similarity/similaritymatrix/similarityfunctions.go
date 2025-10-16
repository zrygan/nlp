package similaritymatrix

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// load all words from the corpus
func LoadCorpusWords(index map[string]map[string]string) map[string][]string {
	langs := make(map[string][]string)

	for lang, verseMap := range index {
		for _, path := range verseMap {
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			text := string(content)
			words := strings.Fields(text)
			langs[lang] = append(langs[lang], words...)
		}
	}

	return langs
}


// gets trigrams of a word and returns it in an array
func GetTrigrams(word string) []string {
	padded := "  " + strings.ToLower(word) + " "
	trigrams := make([]string, 0, len(padded)-2)
	for i := 0; i < len(padded)-2; i++ {
		trigrams = append(trigrams, padded[i:i+3])
	}
	return trigrams
}

// traverses the index and builds trigram frequencies per language
func BuildTrigramCounts(index map[string]map[string]string) (map[string]map[string]int, error) {
    trigramCounts := make(map[string]map[string]int)

    for lang, fileMap := range index {
        fmt.Printf("Processing language: %s (%d files)\n", lang, len(fileMap))
        trigramCounts[lang] = make(map[string]int)

        for _, filePath := range fileMap {
            file, err := os.Open(filePath)
            if err != nil {
                return nil, fmt.Errorf("failed to open %s: %v", filePath, err)
            }

            scanner := bufio.NewScanner(file)
            for scanner.Scan() {
                line := strings.TrimSpace(scanner.Text())
                if line == "" {
                    continue
                }

                words := strings.Fields(line)
                for _, word := range words {
                    for _, tri := range GetTrigrams(word) {
                        trigramCounts[lang][tri]++
                    }
                }
            }

            file.Close()
        }

        fmt.Printf("Done %s: %d unique trigrams\n", lang, len(trigramCounts[lang]))
    }

    return trigramCounts, nil
}

// computes the Jaccard similarity
func ComputeJaccardSimilarity(a, b map[string]int) float64 {
    union := make(map[string]struct{})
    intersection := 0

    for k := range a {
        union[k] = struct{}{}
        if _, exists := b[k]; exists {
            intersection++
        }
    }

    for k := range b {
        union[k] = struct{}{}
    }

    if len(union) == 0 {
        return 0.0
    }

    return float64(intersection) / float64(len(union))
}