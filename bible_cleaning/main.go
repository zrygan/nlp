package main

import (
	"fmt"
	"os"
	"regexp"
	"sync"
)

func corpus() {
	// 1189 is the total number of chapters in the English Bible
	total := 80

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

	var wg sync.WaitGroup
	var mu sync.Mutex
	for language, root := range bibles {
		count := 1

		// parallelize this
		wg.Add(1)
		go func(lang, url string, count *int) {
			defer wg.Done()
			res := ScrapeAndParse(
				url,
				lang,
				"corpus/"+lang,
				cleaningRes,
				make(map[string]bool),
				count,
				total,
			)
			mu.Lock()
			corpusSizes[lang] = res
			mu.Unlock()
		}(language, root, &count)
	}

	wg.Wait()

	sum := 0
	for language, corpusSize := range corpusSizes {
		sum += corpusSize
		fmt.Println(language, " : ", corpusSize)
	}

	fmt.Println("Sig", " : ", sum)
}

func main() {
	if len(os.Args) < 2 {
		panic("No argument provided")
	}

	switch os.Args[1] {
	case "corpus":
		corpus()
	default:
		panic("Non-exaustive switch-case or argument not found")
	}
}
