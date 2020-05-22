package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ksanta/word-stallion/dao"
	"github.com/ksanta/word-stallion/model"
	"github.com/ksanta/word-stallion/scraper"
	"os"
	"strconv"
)

var (
	limit    int
	wordsDao *dao.WordsDao
)

func init() {
	limit, _ = strconv.Atoi(os.Getenv("LIMIT"))
	wordsDao = dao.NewWordsDao(os.Getenv("WORDS_BUCKET"))
}

func handler() error {
	fmt.Println("Scraping", limit, "words")
	// Get a channel which pumps out word definitions
	myScraper := scraper.NewMeriamScraper(limit)
	wordsChan := myScraper.Scrape()

	// Initialise words with enough capacity
	words := make(model.Words, 0, limit)

	count := 0
	for word := range wordsChan {
		words = append(words, word)
		count++
		if count%100 == 0 {
			fmt.Println("Scraped", count)
		}
	}

	return wordsDao.SaveWords(words)
}

func main() {
	lambda.Start(handler)
}
