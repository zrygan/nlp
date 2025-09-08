package main

import (
	"regexp"
	"sync"

	"github.com/zrygan.nlp/bible_cleaning/scraper"
)

func main() {
	total := 10

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
		"bik": "https://www.bible.com/versions/890-mbbbik92-marahay-na-bareta-biblia",
		"war": "https://www.bible.com/bible/2198/GEN.1.MBBSAM",
		"pam": "https://www.bible.com/versions/1141-pmpv-ing-mayap-a-balita-biblia",
		"pag": "https://www.bible.com/bible/2194/GEN.1.MBBPAN83",
		"tiu": "https://www.bible.com/bible/2812/MAT.1.YBT",
		"cbk": "https://www.bible.com/versions/1129-cbknt-el-nuevo-testamento",
		"prf": "https://www.bible.com/bible/438/MAT.1.PRF",
		"tsg": "https://www.bible.com/versions/1319-tsg-kitab-inji",
		"rol": "https://www.bible.com/bible/2244/MAT.1.BKR",
		"msb": "https://www.bible.com/bible/1222/MAT.1.MSB",
		"krj": "https://www.bible.com/bible/1489/MAT.1.KRJNT",
		"tao": "https://www.bible.com/bible/2364/MAT.1.SNT",
	}
	var wg sync.WaitGroup
	for language, root := range bibles {
		count := 0
		wg.Add(1)
		go func(lang, url string, count *int) {
			defer wg.Done()
			scraper.ScrapeAndParse(root, language, "corpus/"+language, cleaningRes, make(map[string]bool), count, total)
		}(language, root, &count)
	}

	wg.Wait()
}
