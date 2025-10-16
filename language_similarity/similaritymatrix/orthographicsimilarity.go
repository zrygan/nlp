package similaritymatrix

import (
	"fmt"
	"os"
	"strings"
)

// builds the similarity matrix using frequency-aware Jaccard
func BuildJaccardSimilarityMatrix(trigramCounts map[string]map[string]int) map[string]map[string]float64 {
    matrix := make(map[string]map[string]float64)

    langs := make([]string, 0, len(trigramCounts))
    for lang := range trigramCounts {
        langs = append(langs, lang)
    }

    for _, langA := range langs {
        matrix[langA] = make(map[string]float64)
        for _, langB := range langs {
            if langA == langB {
                matrix[langA][langB] = 1.0
                continue
            }
            sim := ComputeJaccardSimilarity(trigramCounts[langA], trigramCounts[langB])
            matrix[langA][langB] = sim
        }
    }

    return matrix
}


// saves the actual matrix to a .tsv file
func SaveOrthographicMatrix(matrix map[string]map[string]float64, outPath string) error {
    f, err := os.Create(outPath)
    if err != nil {
        return err
    }
    defer f.Close()

    // Extract language list
    langs := make([]string, 0, len(matrix))
    for lang := range matrix {
        langs = append(langs, lang)
    }

    // Write header row
    fmt.Fprintf(f, "lang\t%s\n", strings.Join(langs, "\t"))

    // Write matrix rows
    for _, langA := range langs {
        fmt.Fprintf(f, "%s", langA)
        for _, langB := range langs {
            fmt.Fprintf(f, "\t%.4f", matrix[langA][langB])
        }
        fmt.Fprintln(f)
    }

    fmt.Printf("Saved orthographic similarity matrix: %s\n", outPath)
    return nil
}
