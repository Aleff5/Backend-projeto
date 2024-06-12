package database

import (
	"context"
	"projetov2/Backend-projeto/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindOne(filter primitive.D) (*models.Usuario, error) {
	client := ConnectBd()
	collection := client.Database("ProjetoLTP2").Collection("Usuarios")

	var usuario models.Usuario
	err := collection.FindOne(context.Background(), filter).Decode(&usuario)
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}
