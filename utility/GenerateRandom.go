package utility

import (
	"math/rand"
	database "projetov2/Backend-projeto/Database"

	"go.mongodb.org/mongo-driver/bson"
)

func GenerateRandomImage() (string, error) {
	listaDeImagens, _ := database.GetFilenames()
	randomIndex := rand.Intn(len(listaDeImagens))

	filter := bson.D{
		{Key: "filename", Value: listaDeImagens[randomIndex]},
	}

	resultadoBusca, erroNaBusca := database.FindOneImage(filter)

	if erroNaBusca != nil {
		return "", erroNaBusca
	}
	return resultadoBusca.FileUrl, nil

	// input := &s3.GetObjectInput{
	// 	Bucket: aws.String("projeto-ltp2"),
	// 	Key:    aws.String(listaDeImagens[randomIndex]),
	// }
	// client := awsfunctions.Set()

	// objeto, err := client.GetObject(context.Background(), input)
	// if err != nil {
	// 	return nil, err

	// }
	// return objeto.Body, nil
}
