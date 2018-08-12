package s3

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Downloader struct {
	manager           s3manager.Downloader
	bucket, key, dest string
}

func NewDownloader(s *session.Session, bucket, key, dest string) *Downloader {
	return &Downloader{
		manager: *s3manager.NewDownloader(s),
		bucket:  bucket,
		key:     key,
		dest:    dest,
	}
}

func (d Downloader) Download() (string, error) {
	file, err := os.Create(d.dest)
	if err != nil {
		return "", err
	}
	defer file.Close()

	numBytes, err := d.manager.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(d.bucket),
			Key:    aws.String(d.key),
		})

	if err != nil {
		return "", err
	}
	log.Println("Downloaded", file.Name(), numBytes, "bytes")

	return file.Name(), nil
}
