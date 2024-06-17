package awsfunctions

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ListObjects(client *s3.Client, bucket string) ([]string, error) {
	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("projeto-ltp2"),
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var conteudo []string

	for _, object := range output.Contents {

		conteudo = append(conteudo, aws.ToString(object.Key))

		log.Println(conteudo)

	}

	return conteudo, nil
}

// func GetObject(client *s3.Client) (string, error) {

// }

func DeleteObject(client *s3.Client, bucket, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := client.DeleteObject(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}
