package main

import (
	"fmt"

	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	lambda.Start(handler)
}

func handler(s3Event events.S3Event) error {
	fmt.Print("lambda called")

	bucket := s3Event.Records[0].S3.Bucket.Name
	key := s3Event.Records[0].S3.Object.Key

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")}),
	)
	downloader := s3manager.NewDownloader(sess)

	file, _ := os.Create("downloaded_file.zip")
	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	return nil
}
