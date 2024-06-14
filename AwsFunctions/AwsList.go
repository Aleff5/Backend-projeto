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
		Bucket: aws.String("my-bucket"),
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
