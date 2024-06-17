package utility

import (
	"context"
	"io"
	"math/rand"
	awsfunctions "projetov2/Backend-projeto/AwsFunctions"
	database "projetov2/Backend-projeto/Database"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GenerateRandomImage() (io.ReadCloser, error) {
	listaDeImagens, _ := database.GetFilenames()
	randomIndex := rand.Intn(len(listaDeImagens))

	input := &s3.GetObjectInput{
		Bucket: aws.String("projeto-ltp2"),
		Key:    aws.String(listaDeImagens[randomIndex]),
	}
	client := awsfunctions.Set()

	objeto, err := client.GetObject(context.Background(), input)
	if err != nil {
		return nil, err

	}
	return objeto.Body, nil
}
