package types

import "regexp"

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

type LanguageClass struct {
	Language string
	OutputDir string
}