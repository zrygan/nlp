package sentencealignment

import (
	"strings"
	"fmt"
	"github.com/zrygan.nlp/bible_cleaning/config"
	"github.com/zrygan.nlp/bible_cleaning/types"
)

/*
	Assumes text is cleaned.
	Lowercases and splits on whitespace.
*/
func Tokenize(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

/*
	Generates n-grams from a list of tokens.
	Returns a list of n-grams as strings, or an empty list if n is less than 1 or greater than the number of tokens.
*/
func NGrams(tokens []string, n int) []string {
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
	Calculates the Dice coefficient similarity between two sets of n-grams.
	Returns a float64 value between 0 and 1, where 1 means identical sets and 0 means no overlap.
*/
func NGramDiceSimilarity(ngrams1, ngrams2 []string) float64 {
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
	Does a sentence similarity based on a combination of Dice coefficient of n-grams and length ratio.
	Based on "A Fast, Flexible Model for Sentence Alignment" by Daniel M. Cer et al.
	https://aclanthology.org/W17-2511.pdf
*/
func SentenceSimilarity(sent1, sent2 string) float64 {
	tokens1 := Tokenize(sent1)
	tokens2 := Tokenize(sent2)

	if len(tokens1) == 0 || len(tokens2) == 0 {
		return 0.0
	}

	// Choose n based on the length of the shorter sentence
	n := 1
	if len(tokens1) >= 3 && len(tokens2) >= 3 {
		n = 3
	} else if len(tokens1) >= 2 && len(tokens2) >= 2 {
		n = 2
	}

	ngrams1 := NGrams(tokens1, n)
	ngrams2 := NGrams(tokens2, n)
	
	NGramDiceSim := NGramDiceSimilarity(ngrams1, ngrams2)
	lenSrc := float64(len([]rune(sent1)))
	lenTgt := float64(len([]rune(sent2)))
	LenRatio	 :=  float64(min(lenSrc, lenTgt)) /
                 float64(max(lenSrc, lenTgt))

	return config.LENGTH_RATIO_BIAS * LenRatio + config.DICE_SIMILARITY_THRESHOLD * NGramDiceSim
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