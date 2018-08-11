package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {

	events := events.S3Event{
		Records: []events.S3EventRecord{
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{Name: "zipped-artifact"},
					Object: events.S3Object{Key: "test.zip"},
				},
			},
		},
	}

	handler(events)
}
