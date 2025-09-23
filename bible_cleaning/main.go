package main

import (
	"fmt"
	"os"
	"regexp"
	"sync"

	parallelcorpus "github.com/zrygan.nlp/bible_cleaning/parallelbuilder"
	"github.com/zrygan.nlp/bible_cleaning/scraper"
	"github.com/zrygan.nlp/bible_cleaning/types"
)

// initialize sets up the initial parameters for the webscraping process
func initialize() (int, map[string]string, map[string]int) {
	chapterLimit := 80

	bibles := map[string]string{
		//  ISO 639: Root URL
		"tgl": "https://www.bible.com/bible/2195/GEN.1.ABTAG01",
		"ceb": "https://www.bible.com/bible/562/GEN.1.RCPV",
		"ilo": "https://www.bible.com/bible/782/GEN.1.RIPV",
		"jil": "https://www.bible.com/bible/2190/GEN.1.MBBHIL12",
		"bik": "https://www.bible.com/bible/890/GEN.1.MBBBIK92",
		"war": "https://www.bible.com/bible/2198/GEN.1.MBBSAM",
		"pam": "https://www.bible.com/bible/1141/GEN.1.PMPV",
		"pag": "https://www.bible.com/bible/2194/GEN.1.MBBPAN83",

		"tiu": "https://www.bible.com/bible/2812/MAT.1.YBT",
		"cbk": "https://www.bible.com/bible/1129/MAT.1.CBKNT",
		"prf": "https://www.bible.com/bible/438/MAT.1.PRF",
		"tsg": "https://www.bible.com/bible/1319/MAT.1.TSG",
		"rol": "https://www.bible.com/bible/2244/MAT.1.BKR",
		"msb": "https://www.bible.com/bible/1222/MAT.1.MSB",
		"krj": "https://www.bible.com/bible/1489/MAT.1.KRJNT",
		"tao": "https://www.bible.com/bible/2364/MAT.1.SNT",
	}

	corpusSizes := map[string]int{
		"tgl": 0,
		"ceb": 0,
		"ilo": 0,
		"jil": 0,
		"bik": 0,
		"war": 0,
		"pam": 0,
		"pag": 0,
		"tiu": 0,
		"cbk": 0,
		"prf": 0,
		"tsg": 0,
		"rol": 0,
		"msb": 0,
		"krj": 0,
		"tao": 0,
	}

	return chapterLimit, bibles, corpusSizes
}

// webscrapeBibles handles concurrent webscraping of multiple bibles
func webscrapeBibles(
	bibleURLs map[string]string,
	corpusSizes map[string]int,
	cleaningConfig []types.FindReplaceTuple[*regexp.Regexp],
	chapterLimit int,
) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for language, root := range bibleURLs {
		chapterCount := 1

		wg.Add(1)

		classification := types.LanguageClass{Language: language, OutputDir: "corpus/" + language}

		go func(lang *types.LanguageClass, bibleURL string, chapterCount *int) {
			defer wg.Done()

			res := scraper.WebscrapeAndParse(bibleURL, lang, &cleaningConfig, make(map[string]bool), chapterCount, chapterLimit)
			//res := scraper.ConcurrentWebscrapeAndParse(bibleURL, lang, &cleaningConfig, chapterLimit, 5)

			// Critical Section: Update shared map
			mu.Lock()
			corpusSizes[lang.Language] = res
			mu.Unlock()

		}(&classification, root, &chapterCount)
	}

	wg.Wait()
}

// summarizeCorpus prints the corpus sizes per language and the total sum
func summarizeCorpus(corpusSizes map[string]int) {
	sum := 0

	for language, corpusSize := range corpusSizes {
		sum += corpusSize
		fmt.Println(language, " : ", corpusSize)
	}

	fmt.Println("Sig", " : ", sum)
}

func parallizeCorpus() {
	corpora, err := parallelcorpus.GenerateParallelCorpus()
	if err != nil {
		panic(err)
	}

	for _, c := range corpora {
		fmt.Printf("%s <--> %s (%d pairs)\n", c.SourceLang, c.TargetLang, len(c.Pairs))
		err = c.SaveAsJSON(fmt.Sprintf("%s_%s.json", c.SourceLang, c.TargetLang))
		fmt.Printf("%s\n", err)
	}
}

// getCorpus orchestrates the entire process of webscraping and corpus generation
func getCorpus() {
	// 1189 is the chapterLimit number of chapters in the English Bible
	chapterLimit, bibles, corpusSizes := initialize()

	cleaningTuples := types.TurnToRegexpsTuple([]types.FindReplaceTuple[string]{
		{
			Find:    `[^a-zA-Z0-9\s\.\,\;\:\!\?\'\"-]+`,
			Replace: "",
		}, {
			Find:    `[:\s]+$`,
			Replace: "",
		}, {
			Find:    `^[\d#:\s]+`,
			Replace: "",
		},
	})

	webscrapeBibles(bibles, corpusSizes, cleaningTuples, chapterLimit)

	summarizeCorpus(corpusSizes)

	parallizeCorpus()
}

func main() {

	if len(os.Args) < 2 {
		panic("No argument provided")
	}

	switch os.Args[1] {
	case "corpus":
		getCorpus()
	case "parallelize":
		parallizeCorpus()
	default:
		panic("Non-exaustive switch-case or argument not found.")
	}

}
