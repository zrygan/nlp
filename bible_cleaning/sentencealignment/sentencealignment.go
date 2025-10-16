package sentencealignment

import (
	"fmt"
	"strings"
<<<<<<< HEAD

=======
	"unicode"

	"github.com/xrash/smetrics"
>>>>>>> 7ecbd5fd670f8dfc71a1b8d24ef2f53232acf56c
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
<<<<<<< HEAD
	Calculates the Dice coefficient similarity between two sets of n-grams.
	Returns a float64 value between 0 and 1, where 1 means identical sets and 0 means no overlap.
=======
Calculates the Dice coefficient similarity between two sets of n-grams by words.
Returns a float64 value between 0 and 1, where 1 means identical sets and 0 means no overlap.
>>>>>>> 7ecbd5fd670f8dfc71a1b8d24ef2f53232acf56c
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
<<<<<<< HEAD
	Does a sentence similarity based on a combination of Dice coefficient of n-grams and length ratio.
	Based on "A Fast, Flexible Model for Sentence Alignment" by Daniel M. Cer et al.
	https://aclanthology.org/W17-2511.pdf
*/
func SentenceSimilarity(sent1, sent2 string) float64 {
	tokens1 := Tokenize(sent1)
	tokens2 := Tokenize(sent2)
=======
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

func LengthRatioSimilarity(sent1, sent2 string) float64 {
	lenSrc := float64(len([]rune(sent1)))
	lenTgt := float64(len([]rune(sent2)))
	LenRatio := float64(min(lenSrc, lenTgt)) / float64(max(lenSrc, lenTgt))
	return LenRatio
}

func ProperNounSimilarity(a, b string) float64 {
	return smetrics.JaroWinkler(a, b, 0.7, 4)
}

func ProperNounOverlapScore(src, tgt string, cache *types.ProperNounCache) float64 {
	srcTokens := strings.Fields(src)
	tgtTokens := strings.Fields(tgt)

	var matches int
	var total int

	for _, s := range srcTokens {
		s = strings.Trim(s, ".,;:!?\"'")
		if _, ok := cache.Words[s]; !ok {
			continue
		}
		total++
		for _, t := range tgtTokens {
			t = strings.Trim(t, ".,;:!?\"'")
			if _, ok := cache.Words[t]; !ok {
				continue
			}
			if ProperNounSimilarity(s, t) > 0.85 {
				matches++
				break
			}
		}
	}

	if total == 0 {
		return 0.0
	}
	return float64(matches) / float64(total)
}

/*
Does a sentence similarity based on a combination of Dice coefficient of n-grams and length ratio.
Based on "A Fast, Flexible Model for Sentence Alignment" by Daniel M. Cer et al.
https://aclanthology.org/W17-2511.pdf
*/
func SentenceSimilarity(sent1, sent2 string, cache *types.ProperNounCache) float64 {
	sent1 = strings.TrimSpace(sent1)
	sent2 = strings.TrimSpace(sent2)
>>>>>>> 7ecbd5fd670f8dfc71a1b8d24ef2f53232acf56c

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

<<<<<<< HEAD
	ngrams1 := NGrams(tokens1, n)
	ngrams2 := NGrams(tokens2, n)
	
	NGramDiceSim := NGramDiceSimilarity(ngrams1, ngrams2)
	lenSrc := float64(len([]rune(sent1)))
	lenTgt := float64(len([]rune(sent2)))
	LenRatio	 :=  float64(min(lenSrc, lenTgt)) /
                 float64(max(lenSrc, lenTgt))

	return config.LENGTH_RATIO_BIAS * LenRatio + (config.DICE_SIMILARITY_THRESHOLD) * NGramDiceSim
=======
	// Compute character n-gram Dice similarity
	NGramDiceSim := CharNGramDiceSimilarity(sent1, sent2, n)
	LenRatio := LengthRatioSimilarity(sent1, sent2)
	PropSim := ProperNounOverlapScore(sent1, sent2, cache)
	return config.LENGTH_RATIO_SIMILARITY_BIAS*LenRatio + config.NGRAMS_DICE_SIMILARITY_BIAS*NGramDiceSim + config.PROPER_NOUNS_SIMILARITY_BIAS*PropSim
>>>>>>> 7ecbd5fd670f8dfc71a1b8d24ef2f53232acf56c
}

func AlignSentencesByGaleChurchDP(srcSents, tgtSents []string, verseID string, cache *types.ProperNounCache) []types.TextPair {
	m, n := len(srcSents), len(tgtSents)
	if m == 0 || n == 0 {
		return nil
	}

<<<<<<< HEAD
	// DP matrix and backpointers
=======
	const maxMerge = 5 // support up to 7 sentences merged per side

	// DP tables
>>>>>>> 7ecbd5fd670f8dfc71a1b8d24ef2f53232acf56c
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
<<<<<<< HEAD
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
=======
			for srcCount := 1; srcCount <= maxMerge && i-srcCount >= 0; srcCount++ {
				for tgtCount := 1; tgtCount <= maxMerge && j-tgtCount >= 0; tgtCount++ {

					// concatenate source/tgt group
					srcGroup := strings.Join(srcSents[i-srcCount:i], " ")
					tgtGroup := strings.Join(tgtSents[j-tgtCount:j], " ")

					score := dp[i-srcCount][j-tgtCount] + SentenceSimilarity(srcGroup, tgtGroup, cache)
					if score > dp[i][j] {
						dp[i][j] = score
						bt[i][j] = struct{ srcCount, tgtCount int }{srcCount, tgtCount}
					}
>>>>>>> 7ecbd5fd670f8dfc71a1b8d24ef2f53232acf56c
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
