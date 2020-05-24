package dao

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ksanta/word-stallion/model"
	"io"
	"log"
)

type WordsDao struct {
	bucketName      string
	downloadService *s3manager.Downloader
	uploadService   *s3manager.Uploader
}

const KEY = "words.txt"

func NewWordsDao(bucketName string) *WordsDao {
	mySession := session.Must(session.NewSession())

	return &WordsDao{
		bucketName:      bucketName,
		downloadService: s3manager.NewDownloader(mySession),
		uploadService:   s3manager.NewUploader(mySession),
	}
}

func (wordsDao *WordsDao) SaveWords(words model.Words) error {
	// buffer collects the bytes
	buffer := &bytes.Buffer{}
	// csvWriter writes bytes in CSV format
	csvWriter := csv.NewWriter(buffer)

	// Stream these into a byte array
	for _, word := range words {
		err := csvWriter.Write(word.ToStringSlice())
		if err != nil {
			return fmt.Errorf("error writing csv: %w", err)
		}
	}
	csvWriter.Flush()

	// Save the byte array to S3
	putObjectInput := &s3manager.UploadInput{
		Body:   buffer,
		Bucket: aws.String(wordsDao.bucketName),
		Key:    aws.String(KEY),
	}

	_, err := wordsDao.uploadService.Upload(putObjectInput)
	return err
}

func (wordsDao *WordsDao) GetWords() (model.Words, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(wordsDao.bucketName),
		Key:    aws.String(KEY),
	}

	_, err := wordsDao.downloadService.Download(buf, getObjectInput)
	if err != nil {
		return nil, err
	}

	words := model.Words{}
	wordReader := csv.NewReader(bytes.NewBuffer(buf.Bytes()))
	for {
		record, err := wordReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		word := model.NewWord(record)
		words = append(words, word)
	}
	return words, nil
}
