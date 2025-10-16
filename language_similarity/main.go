package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	similaritymatrix "language_similarity/similaritymatrix"
)

// i stole this function from the bible
func IndexLanguageFileMap(root string) (map[string]map[string]string, error) {
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

func buildOrthographicSimilarityMatrix() {
    fmt.Println("Loading trigram counts...")
    index, err := IndexLanguageFileMap("../bible_cleaning/corpus/by_verses")
    if err != nil {
        panic(err)
    }

    trigramCounts, err := similaritymatrix.BuildTrigramCounts(index)
    if err != nil {
        panic(err)
    }

    matrix := similaritymatrix.BuildJaccardSimilarityMatrix(trigramCounts)

    outPath := "similaritymatrix/orthographic_similarity_matrix.tsv"
    if err := similaritymatrix.SaveOrthographicMatrix(matrix, outPath); err != nil {
        panic(err)
    }
}

func buildPhoneticSimilarityMatrix() {
    index, err := IndexLanguageFileMap("../bible_cleaning/corpus/by_verses")
    if err != nil {
        panic(err)
    }

	trigramCounts, err := similaritymatrix.BuildTrigramCounts(index)
    if err != nil {
        panic(err)
    }
	phoneticSets := similaritymatrix.BuildPhoneticCountsFromTrigrams(trigramCounts)
    matrix := similaritymatrix.BuildPhoneticSimilarityMatrix(phoneticSets)

    outPath := "similaritymatrix/phonetic_similarity_matrix.tsv"
    if err := similaritymatrix.SavePhoneticMatrix(matrix, outPath); err != nil {
        panic(err)
    }
}

func main() {

	if len(os.Args) < 2 {
		panic("No argument provided")
	}

	switch os.Args[1] {
	case "orthographic":
		buildOrthographicSimilarityMatrix()
	case "phonetic":
		buildPhoneticSimilarityMatrix()
	default:
		panic("Non-exaustive switch-case or argument not found.")
	}
}
