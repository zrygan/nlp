package scraper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// ScrapeChapter scrapes a single chapter: returns verses and chapter name
func scrapeChapter(chapterURL string, cleaningRes []*regexp.Regexp) ([]string, string, error) {
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
				for _, re := range cleaningRes {
					clean = re.ReplaceAllString(clean, "")
				}
				if clean != "" {
					verseTexts = append(verseTexts, clean)
				}
			}
		})

		if len(verseTexts) > 0 {
			verses = append(verses, fmt.Sprintf("%s: %s", verseNumber, strings.Join(verseTexts, " ")))
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
func saveChapter(outputDir, language, chapterName string, verses []string) error {
	if chapterName == "" {
		chapterName = "chapter"
	}
	chapterClean := regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(chapterName, "_")
	filename := fmt.Sprintf("%s_%s.txt", language, chapterClean)
	filePath := filepath.Join(outputDir, filename)

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	err := os.WriteFile(filePath, []byte(strings.Join(verses, "\n")), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Saved:", filePath)
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

// ScrapeAndParse recursively scrapes chapters and returns total verses
func ScrapeAndParse(
	url string,
	language string,
	outputDir string,
	cleaningRes []*regexp.Regexp,
	visited map[string]bool,
	chapterCounter *int,
	maxCount int,
) int {
	// base/edge
	if visited[url] || *chapterCounter > maxCount {
		return 0
	}
	visited[url] = true

	// Scrape current chapter
	verses, chapterName, err := scrapeChapter(url, cleaningRes)
	if err != nil {
		log.Println("Error scraping chapter:", err)
		return 0
	}

	if len(verses) > 0 {
		if err := saveChapter(outputDir, language, chapterName, verses); err != nil {
			log.Println("Error saving chapter:", err)
		}
	}

	totalVerses := len(verses)

	*chapterCounter++

	// Get next chapter
	nextURL, err := getNextChapterURL(url)
	if err != nil {
		log.Println("Error getting next chapter URL:", err)
		return totalVerses
	}

	if nextURL != "" && !visited[nextURL] {
		// maybe we can parallelize this to mkae it faster?
		// recursive call
		totalVerses += ScrapeAndParse(nextURL,
			language,
			outputDir, cleaningRes,
			visited,
			chapterCounter,
			maxCount,
		)
	}

	return totalVerses
}
