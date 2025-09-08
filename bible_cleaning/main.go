package main

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/zrygan.nlp/bible_cleaning/scraper"
)

func main() {
	total := 1189

	cleaningRes := []*regexp.Regexp{
		regexp.MustCompile(`[^a-zA-Z0-9\s\.\,\;\:\!\?\'\"-]+`),
		regexp.MustCompile(`[:\s]+$`),
		regexp.MustCompile(`^[\d#:\s]+`),
	}

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
	var wg sync.WaitGroup
	for language, root := range bibles {
		count := 1
		wg.Add(1)
		go func(lang, url string, count *int) {
			defer wg.Done()
			fmt.Println(
				language,
				scraper.ScrapeAndParse(
					root,
					language,
					"corpus/"+language,
					cleaningRes,
					make(map[string]bool),
					count,
					total,
				),
			)
		}(language, root, &count)
	}

	wg.Wait()
}
