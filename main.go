package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/toshi0607/unzipper/s3"
	"github.com/toshi0607/unzipper/zip"
)

const (
	artifactPath = "/tmp/artifact/"
	zipPath      = artifactPath + "zipped/"
	unzipPath    = artifactPath + "unzipped/"
	tempZip      = "temp.zip"
	dirPerm      = 0777
)

var (
	now              string
	zipContentPath   string
	unzipContentPath string
	destBucket       string
)

func init() {
	now = strconv.Itoa(int(time.Now().UnixNano()))
	zipContentPath = zipPath + now + "/"
	unzipContentPath = unzipPath + now + "/"
	destBucket = os.Getenv("DEST_BUCKET")
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
		Region: aws.String(endpoints.ApNortheast1RegionID)}),
	)

	downloader := s3.NewDownloader(sess, bucket, key, zipContentPath+tempZip)
	downloadedZipPath, err := downloader.Download()
	if err != nil {
		log.Fatal(err)
	}

	if err := zip.Unzip(downloadedZipPath, unzipContentPath); err != nil {
		log.Fatal(err)
	}

	uploader := s3.NewUploader(sess, unzipPath, destBucket)
	if err := uploader.Upload(); err != nil {
		log.Fatal(err)
	}

	log.Printf("%s unzipped to S3 bucket: %s", downloadedZipPath, destBucket)

	return nil
}

func prepareDirectory() error {
	if _, err := os.Stat(artifactPath); err == nil {
		if err := os.RemoveAll(artifactPath); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(zipContentPath, dirPerm); err != nil {
		return err
	}
	if err := os.MkdirAll(unzipContentPath, dirPerm); err != nil {
		return err
	}

	return nil
}
