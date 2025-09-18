package scraper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
	"github.com/zrygan.nlp/bible_cleaning/types"
	"github.com/gocolly/colly"
	"sync"
	"sync/atomic"
)

// ScrapeChapter scrapes a single chapter: returns verses and chapter name
func scrapeChapter(chapterURL string, cleaningConfig *[]types.FindReplaceTuple[*regexp.Regexp]) ([]string, string, error) {
	var verses []string
	var chapterName string

	c := colly.NewCollector()

	c.OnHTML("span[class^='ChapterContent_verse__']", func(e *colly.HTMLElement) {
		verseNumber := e.ChildText("span[class^='ChapterContent_label__']")
		verseNumber = regexp.MustCompile(`\D`).ReplaceAllString(verseNumber, "")

		var verseTexts []string
		e.DOM.Find("span[class^='ChapterContent_content__']").Each(func(_ int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				clean := text

				for _, tuple := range (*cleaningConfig) {
					clean = tuple.Find.ReplaceAllString(clean, tuple.Replace.String())
				}
				
				if clean != "" {
					verseTexts = append(verseTexts, clean)
				}
			}
		})

		if len(verseTexts) > 0 {
			// remove the verseNumber here.
			verses = append(verses, fmt.Sprintf("%s", strings.Join(verseTexts, " ")))
		}
	})

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		chapterName = strings.TrimSpace(e.Text)
	})

	err := c.Visit(chapterURL)
	if err != nil {
		return nil, "", err
	}

	return verses, chapterName, nil
}

// SaveChapter saves the verses to a file
func saveChapter(lang types.LanguageClass, chapterName string, verses []string) error {
	if chapterName == "" {
		chapterName = "chapter"
	}

	chapterClean := regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(chapterName, "_")
	filename := fmt.Sprintf("%s_%s.txt", lang.Language, chapterClean)
	filePath := filepath.Join(lang.OutputDir, filename)

	if err := os.MkdirAll(lang.OutputDir, os.ModePerm); err != nil {
		return err
	}

	err := os.WriteFile(filePath, []byte(strings.Join(verses, "\n")), 0644)
	if err != nil {
		return err
	}

	// fmt.Println("Saved:", filePath)
	return nil
}

// GetNextChapterURL finds the next chapter link
func getNextChapterURL(chapterURL string) (string, error) {
	var nextURL string
	c := colly.NewCollector()

	c.OnHTML("a[href^='/bible/']", func(e *colly.HTMLElement) {
		title := e.DOM.Find("svg title").Text()
		if title == "Next Chapter" {
			nextURL = "https://www.bible.com" + e.Attr("href")
		}
	})

	err := c.Visit(chapterURL)
	if err != nil {
		return "", err
	}

	return nextURL, nil
}

func countWordsFromURL(url string, langClass types.LanguageClass, cleaningConfig *[]types.FindReplaceTuple[*regexp.Regexp]) int {

	verses, chapterName, err := scrapeChapter(url, cleaningConfig)

	if err != nil {
		log.Println("Error scraping chapter:", err)
		return 0
	}

	if len(verses) > 0 {
		if err := saveChapter(langClass, chapterName, verses); err != nil {
			log.Println("Error saving chapter:", err)
		}
	}

	wordCount := 0
	for _, v := range verses {
		words := strings.Fields(v)
		wordCount += len(words)
	}

	return wordCount
}

// WebscrapeAndParse recursively scrapes chapters and returns total verses
func WebscrapeAndParse(
	websiteURL string,
	langClass *types.LanguageClass,
	cleaningConfig *[]types.FindReplaceTuple[*regexp.Regexp],
	visited map[string]bool,
	chapterCounter *int,
	maxCount int,
) int {
	// base/edge
	if visited[websiteURL] || *chapterCounter > maxCount {
		return 0
	}
	visited[websiteURL] = true

	
	wordCount := countWordsFromURL(websiteURL, *langClass, cleaningConfig)
	*chapterCounter++

	// Get next chapter
	nextURL, err := getNextChapterURL(websiteURL)
	if err != nil {
		log.Println("Error getting next chapter URL:", err)
		return wordCount
	}

	if nextURL != "" && !visited[nextURL] {
		// maybe we can parallelize this to mkae it faster?
		// recursive call
		wordCount += WebscrapeAndParse(
			nextURL,
			langClass, 
			cleaningConfig,
			visited,
			chapterCounter,
			maxCount,
		)
	}

	return wordCount
}

// Context for concurrent web scraping
type WebscrapeContext struct {
	visited       map[string]bool
	visitedMu     *sync.Mutex
	urlCh         chan string																	// Channel for URLs to process
	tasks         *sync.WaitGroup															// Wait group for tracking tasks
	langClass     *types.LanguageClass
	cleaningConfig *[]types.FindReplaceTuple[*regexp.Regexp]
	chapterCounter *atomic.Int64
	maxCount      int
	totalWordCount *int64
}


func prefetchStaringURLs(url string, depth int, ctx *WebscrapeContext) {
		if depth <= 0 {
			return
		}

		ctx.tasks.Add(1)
		ctx.urlCh <- url

		nextURL, err := getNextChapterURL(url)
		if err != nil {
			log.Println("[Prefetch] Error fetching next URL:", err)
			return
		}
		if nextURL != "" {
			prefetchStaringURLs(nextURL, depth-1, ctx)
		}
}


// BFSWebscrape is a worker function for concurrent web scraping
// Check [here](../docs/devs.md)
func BFSWebscrape(ctx *WebscrapeContext) {
	for url := range ctx.urlCh {

		// check visited
		(*ctx.visitedMu).Lock()
		if ctx.visited[url] {
			(*ctx.visitedMu).Unlock()
			ctx.tasks.Done()
			continue
		}
		ctx.visited[url] = true
		(*ctx.visitedMu).Unlock()


		// stop if max reached
		nextCount := ctx.chapterCounter.Add(1)
		if nextCount > int64(ctx.maxCount) {
			ctx.tasks.Done()
			continue
		}

		// process URL
		wordCount := countWordsFromURL(url, *ctx.langClass, ctx.cleaningConfig)
		atomic.AddInt64(ctx.totalWordCount, int64(wordCount))

		// fetch next
		nextURL, err := getNextChapterURL(url)
		if err != nil {
			log.Println("Error getting next chapter URL:", err)
		} else if nextURL != "" {
			(*ctx.visitedMu).Lock()
			if !ctx.visited[nextURL] {
				ctx.tasks.Add(1)
				ctx.urlCh <- nextURL
			}
			(*ctx.visitedMu).Unlock()
		}

		ctx.tasks.Done()
	}
}

func Timer() {
	start := time.Now()
	for {
		time.Sleep(10 * time.Millisecond)

		log.Printf("[Progress] Time: %v\n", time.Since(start))
	}
}

func ConcurrentWebscrapeAndParse(
	startURL string,
	langClass *types.LanguageClass,
	cleaningConfig *[]types.FindReplaceTuple[*regexp.Regexp],
	maxCount int,
	numWorkers int,
) int {
	var (
		visited       = make(map[string]bool)
		visitedMu     sync.Mutex
		totalWordCount int64
		chapterCounter atomic.Int64
	)

	urlCh := make(chan string, 100) // 100 string buffer
	var tasks sync.WaitGroup

	ctx := &WebscrapeContext{
			visited:        visited,
			visitedMu:      &visitedMu,
			urlCh:          urlCh,
			tasks:          &tasks,
			langClass:      langClass,
			cleaningConfig: cleaningConfig,
			chapterCounter: &chapterCounter,
			maxCount:       maxCount,
			totalWordCount: &totalWordCount,
	}

	// enqueue initial URL
	prefetchStaringURLs(startURL, numWorkers, ctx)

	// workers
	for i := 0; i < numWorkers; i++ {
		go BFSWebscrape(ctx)
	}

	// closer
	go func() {
		tasks.Wait()
		close(urlCh)
	}()

	// wait until all workers finish
	tasks.Wait()
	return int(totalWordCount)
}
