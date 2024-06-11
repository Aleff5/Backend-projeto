package database

import (
	"context"
	"projetov2/Backend-projeto/models"

	"go.mongodb.org/mongo-driver/mongo"
)

func InsertOne(usuario models.Usuario) (*mongo.InsertOneResult, error) {
	client := ConnectBd()
	collection := client.Database("ProjetoLTP2").Collection("Usuarios")

	return collection.InsertOne(context.Background(), usuario)

}
