package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/toshi0607/unzipper/zip"
	"golang.org/x/sync/errgroup"
)

const (
	artifactPath = "tmp/artifact/"
	zipPath      = artifactPath + "zipped/"
	unzipPath    = artifactPath + "unzipped/"
	tempZip      = "temp.zip"
)

var (
	now              string
	zipContentPath   string
	unzipContentPath string
)

func init() {
	now = strconv.Itoa(int(time.Now().UnixNano()))
	zipContentPath = zipPath + now + "/"
	unzipContentPath = unzipPath + now + "/"
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, s3Event events.S3Event) error {
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		log.Printf("AwsRequestID: %s", lc.AwsRequestID)
	}

	bucket := s3Event.Records[0].S3.Bucket.Name
	key := s3Event.Records[0].S3.Object.Key

	log.Printf("bucket: %s ,key: %s", bucket, key)

	if err := prepareDirectory(); err != nil {
		log.Fatal(err)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")}),
	)
	downloader := s3manager.NewDownloader(sess)

	file, err := os.Create(zipContentPath + tempZip)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Downloaded", file.Name(), numBytes, "bytes")

	err = zip.Unzip(zipContentPath+tempZip, unzipContentPath)
	if err != nil {
		log.Fatal(err)
	}

	uploader := s3manager.NewUploader(sess)
	eg := errgroup.Group{}

	err = filepath.Walk(unzipPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		eg.Go(func() error {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			key := strings.Replace(file.Name(), unzipPath, "", 1)
			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String("unzipped-artifact"),
				Key:    aws.String(key),
				Body:   file,
			})
			if err != nil {
				return err
			}
			return nil
		})
		return nil
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func prepareDirectory() error {
	if _, err := os.Stat(artifactPath); err == nil {
		if err := os.RemoveAll(artifactPath); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(zipContentPath, 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(unzipContentPath, 0777); err != nil {
		return err
	}

	return nil
}
