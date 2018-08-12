package s3

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/sync/errgroup"
)

type Uploader struct {
	manager   s3manager.Uploader
	src, dest string
}

func NewUploader(s *session.Session, src, dest string) *Uploader {
	return &Uploader{
		manager: *s3manager.NewUploader(s),
		src:     src,
		dest:    dest,
	}
}

func (u Uploader) Upload() error {
	eg := errgroup.Group{}

	err := filepath.Walk(u.src, func(path string, info os.FileInfo, err error) error {
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

			key := strings.Replace(file.Name(), u.src, "", 1)
			_, err = u.manager.Upload(&s3manager.UploadInput{
				Bucket: aws.String(u.dest),
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
