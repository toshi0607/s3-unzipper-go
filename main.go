package main

import (
	"archive/zip"
	"fmt"
	"time"

	"os"

	"io"
	"path/filepath"

	"strconv"

	"log"

	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/sync/errgroup"
)

func main() {
	lambda.Start(handler)
}

func handler(s3Event events.S3Event) error {
	fmt.Print("lambda called")

	bucket := s3Event.Records[0].S3.Bucket.Name
	key := s3Event.Records[0].S3.Object.Key

	//now := time.Now().Format("2006-01-02-15-04-05")
	now := strconv.Itoa(int(time.Now().UnixNano()))
	tmpZipFolder := "tmp/artifact/zipped/"
	tmpUnzipFolder := "tmp/artifact/unzipped/"
	tmpZipPath := tmpZipFolder + now + "/"
	tmpUnzipPath := tmpUnzipFolder + now + "/"

	if _, err := os.Stat("tmp/artifact"); err == nil {
		if err := os.RemoveAll("tmp/artifact"); err != nil {
			log.Fatal(err)
		}
	}

	if err := os.MkdirAll(tmpZipPath, 0777); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(tmpUnzipPath, 0777); err != nil {
		log.Fatal(err)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")}),
	)
	downloader := s3manager.NewDownloader(sess)

	file, err := os.Create(tmpZipPath + "temp.zip")
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
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	err = Unzip(tmpZipPath+"temp.zip", tmpUnzipPath)
	if err != nil {
		log.Fatal(err)
	}

	uploader := s3manager.NewUploader(sess)
	eg := errgroup.Group{}

	err = filepath.Walk(tmpUnzipFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
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

			key := strings.Replace(file.Name(), tmpUnzipFolder, "", 1)
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

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
