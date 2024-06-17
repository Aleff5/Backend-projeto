package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteImage(filter primitive.D) error {
	client := ConnectBd()
	collection := client.Database("ProjetoLTP2").Collection("Imagens")

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
