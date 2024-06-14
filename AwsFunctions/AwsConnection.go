package awsfunctions

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Set() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	return client
}

func AwsUpload(client *s3.Client, bucket, filename string) (string, error) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		log.Fatal(openErr)
		return "", openErr
	}
	defer file.Close()

	uploader := manager.NewUploader(client)

	result, errUpload := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String("image/jpeg"),
	})

	if errUpload != nil {
		log.Fatal(errUpload)
		return "", errUpload
	}

	urlUpload := result.Location

	return urlUpload, nil

}
