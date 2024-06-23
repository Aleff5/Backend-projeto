package database

import (
	"context"
	"projetov2/Backend-projeto/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAllUsers(client *mongo.Client, ctx context.Context) ([]bson.M, error) {
	collection := client.Database("ProjetoLTP2").Collection("Usuarios")

	projection := bson.D{
		{"email", 1},
		{"username", 1},
	}

	cur, err := collection.Find(ctx, bson.D{}, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []bson.M
	for cur.Next(ctx) {
		var result bson.M
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func FindAllImages() ([]bson.M, error) {
	client := ConnectBd()
	collection := client.Database("ProjetoLTP2").Collection("Imagens")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []bson.M
	for cur.Next(ctx) {
		var result bson.M
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func FindUrl() ([]string, error) {
	client := ConnectBd()
	collection := client.Database("ProjetoLTP2").Collection("Imagens")

	// Find all documents in the collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var urls []string

	// Iterate through the cursor and extract FileUrl
	for cursor.Next(context.Background()) {
		var img models.Imagem
		if err := cursor.Decode(&img); err != nil {
			return nil, err
		}
		urls = append(urls, img.FileUrl)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}
