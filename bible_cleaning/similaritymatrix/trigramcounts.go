package similaritymatrix

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getTrigrams(word string) []string {
	padded := "  " + strings.ToLower(word) + " "
	trigrams := make([]string, 0, len(padded)-2)
	for i := 0; i < len(padded)-2; i++ {
		trigrams = append(trigrams, padded[i:i+3])
	}
	return trigrams
}

// BuildTrigramCounts traverses the index and builds trigram frequencies per language.
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
                    for _, tri := range getTrigrams(word) {
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

func SaveTrigramCounts(counts map[string]map[string]int, outdir string) error {
    if err := os.MkdirAll(outdir, os.ModePerm); err != nil {
        return err
    }

    for lang, tris := range counts {
        outPath := fmt.Sprintf("%s/%s_trigrams.tsv", outdir, lang)
        f, err := os.Create(outPath)
        if err != nil {
            return err
        }
        defer f.Close()

        for tri, count := range tris {
            fmt.Fprintf(f, "%s\t%d\n", tri, count)
        }
        f.Close()
        fmt.Printf("Saved: %s\n", outPath)
    }

    return nil
}
