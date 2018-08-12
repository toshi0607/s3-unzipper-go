package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func TestHandler(t *testing.T) {

	events := events.S3Event{
		Records: []events.S3EventRecord{
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{Name: "zipped-artifact"},
					Object: events.S3Object{Key: "sample.zip"},
				},
			},
		},
	}

	ctx := context.Background()
	lc := new(lambdacontext.LambdaContext)
	ctx = lambdacontext.NewContext(ctx, lc)

	err := handler(ctx, events)
	if err != nil {
		t.Error(err)
	}
}
