package database

import (
	"context"
	"projetov2/Backend-projeto/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOneUser(filter primitive.D) (*models.Usuario, error) {
	client := ConnectBd()
	collection := client.Database("ProjetoLTP2").Collection("Usuarios")

	var usuario models.Usuario
	err := collection.FindOne(context.Background(), filter).Decode(&usuario)
	if err != nil {

		return nil, err
	}
	return &usuario, nil
}

func GetFilenames() ([]string, error) {
	client := ConnectBd()
	collection := client.Database("ProjetoLTP2").Collection("Imagens")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Define a projeção para retornar apenas o campo "filename"
	projection := bson.D{{"filename", 1}, {"_id", 0}}
	opts := options.Find().SetProjection(projection)

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var filenames []string
	for cursor.Next(ctx) {
		var result struct {
			Filename string `bson:"filename"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		filenames = append(filenames, result.Filename)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return filenames, nil
}
