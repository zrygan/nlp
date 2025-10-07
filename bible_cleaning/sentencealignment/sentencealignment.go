package sentencealignment

import (
	"strings"
	"fmt"
	"github.com/zrygan.nlp/bible_cleaning/config"
	"github.com/zrygan.nlp/bible_cleaning/types"
	"unicode"
)

/*
	Assumes text is cleaned.
	Lowercases and splits on whitespace.
*/
func WordTokenize(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

/*
	Generates n-grams from a list of tokens.
	Returns a list of n-grams as strings, or an empty list if n is less than 1 or greater than the number of tokens.
*/
func WordNGrams(tokens []string, n int) []string {
	var ngrams []string
	
	if n <= 0 || len(tokens) < n {
		return ngrams
	}
	for i := 0; i <= len(tokens)-n; i++ {
		ngram := strings.Join(tokens[i:i+n], " ")
		ngrams = append(ngrams, ngram)
	}
	return ngrams
}

/*
	Calculates the Dice coefficient similarity between two sets of n-grams by words.
	Returns a float64 value between 0 and 1, where 1 means identical sets and 0 means no overlap.
*/
func NGramWordDiceSimilarity(ngrams1, ngrams2 []string) float64 {
	set1 := make(map[string]struct{})
	set2 := make(map[string]struct{})

	for _, ngram := range ngrams1 {
		set1[ngram] = struct{}{}
	}
	for _, ngram := range ngrams2 {
		set2[ngram] = struct{}{}
	}

	intersectionSize := 0
	for ngram := range set1 {
		if _, exists := set2[ngram]; exists {
			intersectionSize++
		}
	}

	if len(set1)+len(set2) == 0 {
		return 0.0
	}

	return (2.0 * float64(intersectionSize)) / float64(len(set1)+len(set2))
}


/*
	Assumes text is cleaned.
	Lowercases and splits on char.
*/
func CharTokenize(text string) []string {
	text = strings.ToLower(text)
	var tokens []string
	for _, r := range text {
		if !unicode.IsSpace(r) { // skip spaces
			tokens = append(tokens, string(r))
		}
	}
	return tokens
}

/*
	Generates n-grams from a list of tokens.
	Returns a list of n-grams as strings, or an empty list if n is less than 1 or greater than the number of tokens.
*/
func CharNGrams(tokens []string, n int) []string {
	if n <= 0 || len(tokens) < n {
		return nil
	}
	ngrams := make([]string, 0, len(tokens)-n+1)
	for i := 0; i <= len(tokens)-n; i++ {
		ngrams = append(ngrams, strings.Join(tokens[i:i+n], ""))
	}
	return ngrams
}

/*
	Calculates the Dice coefficient similarity between two sets of n-grams.
	Returns a float64 value between 0 and 1, where 1 means identical sets and 0 means no overlap.
*/
func CharNGramDiceSimilarity(text1, text2 string, n int) float64 {
	tokens1 := CharTokenize(text1)
	tokens2 := CharTokenize(text2)

	ngrams1 := CharNGrams(tokens1, n)
	ngrams2 := CharNGrams(tokens2, n)

	set1 := make(map[string]struct{}, len(ngrams1))
	set2 := make(map[string]struct{}, len(ngrams2))

	for _, ng := range ngrams1 {
		set1[ng] = struct{}{}
	}
	for _, ng := range ngrams2 {
		set2[ng] = struct{}{}
	}

	intersection := 0
	for ng := range set1 {
		if _, exists := set2[ng]; exists {
			intersection++
		}
	}

	if len(set1)+len(set2) == 0 {
		return 0.0
	}
	return (2.0 * float64(intersection)) / float64(len(set1)+len(set2))
}

func LengthRatioSimilarity(sent1, sent2 string) float64{
	lenSrc := float64(len([]rune(sent1)))
	lenTgt := float64(len([]rune(sent2)))
	LenRatio := float64(min(lenSrc, lenTgt)) / float64(max(lenSrc, lenTgt))
	return LenRatio
}


/*
	Does a sentence similarity based on a combination of Dice coefficient of n-grams and length ratio.
	Based on "A Fast, Flexible Model for Sentence Alignment" by Daniel M. Cer et al.
	https://aclanthology.org/W17-2511.pdf
*/
func SentenceSimilarity(sent1, sent2 string) float64 {
	sent1 = strings.TrimSpace(sent1)
	sent2 = strings.TrimSpace(sent2)

	if len(sent1) == 0 || len(sent2) == 0 {
		return 0.0
	}

	len1 := len([]rune(sent1))
	len2 := len([]rune(sent2))

	// Dynamically choose n based on shortest length
	n := 3
	if len1 < 3 || len2 < 3 {
		if len1 < 2 || len2 < 2 {
			n = 1 // fallback for very short strings
		} else {
			n = 2
		}
	}

	// Compute character n-gram Dice similarity
	NGramDiceSim := CharNGramDiceSimilarity(sent1, sent2, n)
	LenRatio := LengthRatioSimilarity(sent1, sent2)
	return config.LENGTH_RATIO_BIAS*LenRatio + config.DICE_SIMILARITY_THRESHOLD*NGramDiceSim
}

func AlignSentencesByGaleChurchDP(srcSents, tgtSents []string, verseID string) []types.TextPair {
	m, n := len(srcSents), len(tgtSents)
	if m == 0 || n == 0 {
		return nil
	}

	const maxMerge = 10 // support up to 10 sentences merged per side

	// DP tables
	dp := make([][]float64, m+1)
	bt := make([][]struct{ srcCount, tgtCount int }, m+1)
	for i := range dp {
		dp[i] = make([]float64, n+1)
		bt[i] = make([]struct{ srcCount, tgtCount int }, n+1)
		for j := range dp[i] {
			dp[i][j] = -1e9
		}
	}
	dp[0][0] = 0

	// Transition loop
	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			for srcCount := 1; srcCount <= maxMerge && i-srcCount >= 0; srcCount++ {
				for tgtCount := 1; tgtCount <= maxMerge && j-tgtCount >= 0; tgtCount++ {

					// concatenate source/tgt group
					srcGroup := strings.Join(srcSents[i-srcCount:i], " ")
					tgtGroup := strings.Join(tgtSents[j-tgtCount:j], " ")

					score := dp[i-srcCount][j-tgtCount] + SentenceSimilarity(srcGroup, tgtGroup)
					if score > dp[i][j] {
						dp[i][j] = score
						bt[i][j] = struct{ srcCount, tgtCount int }{srcCount, tgtCount}
					}
				}
			}
		}
	}

	// Backtrack
	var pairs []types.TextPair
	i, j := m, n
	count := 1
	for i > 0 && j > 0 {
		step := bt[i][j]
		if step.srcCount == 0 || step.tgtCount == 0 {
			break
		}

		srcGroup := strings.Join(srcSents[i-step.srcCount:i], " ")
		tgtGroup := strings.Join(tgtSents[j-step.tgtCount:j], " ")

		pairs = append([]types.TextPair{{
			ID:         fmt.Sprintf("%s_%03d", verseID, count),
			SourceText: srcGroup,
			TargetText: tgtGroup,
		}}, pairs...)

		count++
		i -= step.srcCount
		j -= step.tgtCount
	}

	return pairs
}