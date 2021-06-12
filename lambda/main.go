package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	//"github.com/nfnt/resize"
)

type Record struct {
	EventSource    string
	EventSourceArn string
	AWSRegion      string
	S3             events.S3Entity
	SQS            events.SQSMessage
	SNS            events.SNSEntity
}

type Event struct {
	Records []Record
}

func show(ctx context.Context, s3Event events.S3Event) (events.S3Event, error) {

	originalBucket := s3Event.Records[0].S3.Bucket.Name
    //thumbnailBucket :=
    imageKey := s3Event.Records[0].S3.Object.Key 
    //maxWidth
    //maxHeight

	log.Print(originalBucket)
    log.Print(imageKey)

    sess, err := session.NewSession()
    if err != nil {
		log.Print(err)
	}
	s3uploader := s3manager.NewUploader(sess)
    s3downloader := s3manager.NewDownloader(sess)

	//get image from event object
    item := "/tmp/{}"

    file, err := os.Create(item)
    if err != nil {
        log.Print("Unable to open file %q, %v", item, err)
    }

    numBytes, err := s3downloader.Download(file,
        &s3.GetObjectInput{
            Bucket: aws.String(originalBucket),
            Key:    aws.String(imageKey),
        })
    if err != nil {
        log.Print("unable to download")
    }

    log.Print(numBytes)

	//downsize image

	//upload to thumbnail bucket
    _, err = s3uploader.Upload(&s3manager.UploadInput{
        Bucket: "thumbnails-lambda-bucket7",
        Key:    imageKey,
        Body:   file,
        ACL:    aws.String("public-read"),
    }

    return s3Event, nil
}

func main() {
	lambda.Start(show)
}
