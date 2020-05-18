package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ksanta/word-stallion/scraper"
	"os"
	"strconv"
)

var (
	bucketName string
	limit      int
	s3service  *s3manager.Uploader
)

func init() {
	bucketName = os.Getenv("WORDS_BUCKET")
	limit, _ = strconv.Atoi(os.Getenv("LIMIT"))
	mySession := session.Must(session.NewSession())
	s3service = s3manager.NewUploader(mySession)
}

func handler() error {
	fmt.Println("Scraping", limit, "words to", bucketName)
	// Get a channel which pumps out word definitions
	myScraper := scraper.NewMeriamScraper(limit)
	wordsChan := myScraper.Scrape()

	// buffer collects the bytes
	buffer := &bytes.Buffer{}
	// csvWriter writes bytes in CSV format
	csvWriter := csv.NewWriter(buffer)

	// Stream these into a byte array
	count := 0
	for word := range wordsChan {
		err := csvWriter.Write(word.ToStringSlice())
		if err != nil {
			return fmt.Errorf("error writing csv: %w", err)
		}
		count++
		if count%100 == 0 {
			fmt.Println("Scraped", count)
		}
	}
	csvWriter.Flush()

	// Save the byte array to S3
	putObjectInput := &s3manager.UploadInput{
		Body:   buffer,
		Bucket: aws.String(bucketName),
		Key:    aws.String("words.txt"),
	}

	_, err := s3service.Upload(putObjectInput)
	return err
}

func main() {
	lambda.Start(handler)
}
