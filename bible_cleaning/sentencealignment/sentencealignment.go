package sentencealignment

import (
	"fmt"
	"strings"

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

	return config.LENGTH_RATIO_BIAS * LenRatio + (config.DICE_SIMILARITY_THRESHOLD) * NGramDiceSim
}

func AlignSentencesByGaleChurchDP(srcSents, tgtSents []string, verseID string) []types.TextPair {
	m, n := len(srcSents), len(tgtSents)
	if m == 0 || n == 0 {
		return nil
	}

	// DP matrix and backpointers
	dp := make([][]float64, m+1)
	bt := make([][]string, m+1)
	for i := range dp {
		dp[i] = make([]float64, n+1)
		bt[i] = make([]string, n+1)
		for j := range dp[i] {
			dp[i][j] = -1e9
			bt[i][j] = ""
		}
	}
	dp[0][0] = 0

	// Transition states
	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			if i > 0 && j > 0 {
				score := dp[i-1][j-1] + SentenceSimilarity(srcSents[i-1], tgtSents[j-1])
				if score > dp[i][j] {
					dp[i][j] = score
					bt[i][j] = "1:1"
				}
			}
			if i > 0 && j > 1 {
				t := tgtSents[j-2] + " " + tgtSents[j-1]
				score := dp[i-1][j-2] + SentenceSimilarity(srcSents[i-1], t)
				if score > dp[i][j] {
					dp[i][j] = score
					bt[i][j] = "1:2"
				}
			}
			if i > 1 && j > 0 {
				s := srcSents[i-2] + " " + srcSents[i-1]
				score := dp[i-2][j-1] + SentenceSimilarity(s, tgtSents[j-1])
				if score > dp[i][j] {
					dp[i][j] = score
					bt[i][j] = "2:1"
				}
			}
		}
	}
	idFormat := "%s_%03d"
	// Backtrack to reconstruct alignment
	var pairs []types.TextPair
	i, j := m, n
	for i > 0 && j > 0 {
		state := bt[i][j]
		if state == "1:1" {
			pairs = append([]types.TextPair{{
				ID:         fmt.Sprintf(idFormat, verseID, len(pairs)+1),
				SourceText: srcSents[i-1],
				TargetText: tgtSents[j-1],
			}}, pairs...)
			i--
			j--
		} else if state == "1:2" {
			pairs = append([]types.TextPair{{
				ID:         fmt.Sprintf(idFormat, verseID, len(pairs)+1),
				SourceText: srcSents[i-1],
				TargetText: tgtSents[j-2] + " " + tgtSents[j-1],
			}}, pairs...)
			i--
			j -= 2
		} else if state == "2:1" {
			pairs = append([]types.TextPair{{
				ID:         fmt.Sprintf(idFormat, verseID, len(pairs)+1),
				SourceText: srcSents[i-2] + " " + srcSents[i-1],
				TargetText: tgtSents[j-1],
			}}, pairs...)
			i -= 2
			j--
		} else {
			break // alignment path ends
		}
	}

	return pairs
}