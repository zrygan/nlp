package similaritymatrix

import (
	"fmt"
	"os"
	"strings"

	"github.com/twuillemin/doublemetaphone/pkg/doublemetaphone"
)

// convert trigrams to phonetic frequency maps
func BuildPhoneticCountsFromTrigrams(trigramCounts map[string]map[string]int) map[string]map[string]int {
	phoneticCounts := make(map[string]map[string]int)

	for lang, tris := range trigramCounts {
		counts := make(map[string]int)
		for tri, freq := range tris {
			primary, alt := doublemetaphone.DoubleMetaphone(tri)
			if primary != "" {
				counts[primary] += freq
			}
			if alt != "" {
				counts[alt] += freq
			}
		}
		phoneticCounts[lang] = counts
	}

	return phoneticCounts
}

// build phonetic similarity matrix using frequency counts
func BuildPhoneticSimilarityMatrix(phoneticCounts map[string]map[string]int) map[string]map[string]float64 {
	matrix := make(map[string]map[string]float64)
	langsList := make([]string, 0, len(phoneticCounts))
	for lang := range phoneticCounts {
		langsList = append(langsList, lang)
	}

	for _, langA := range langsList {
		matrix[langA] = make(map[string]float64)
		for _, langB := range langsList {
			if langA == langB {
				matrix[langA][langB] = 1.0
			} else {
				sim := ComputeJaccardSimilarity(phoneticCounts[langA], phoneticCounts[langB])
				matrix[langA][langB] = sim
			}
		}
	}

	return matrix
}

// save matrix to a .tsv
func SavePhoneticMatrix(matrix map[string]map[string]float64, outPath string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	langs := make([]string, 0, len(matrix))
	for lang := range matrix {
		langs = append(langs, lang)
	}

	fmt.Fprintf(f, "lang\t%s\n", strings.Join(langs, "\t"))
	for _, langA := range langs {
		fmt.Fprintf(f, "%s", langA)
		for _, langB := range langs {
			fmt.Fprintf(f, "\t%.4f", matrix[langA][langB])
		}
		fmt.Fprintln(f)
	}

	fmt.Printf("Saved phonetic similarity matrix: %s\n", outPath)
	return nil
}
