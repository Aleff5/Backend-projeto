package database

import (
	"context"
	"projetov2/Backend-projeto/models"

	"go.mongodb.org/mongo-driver/mongo"
)

func InsertOneUser(usuario models.Usuario) (*mongo.InsertOneResult, error) {
	client := ConnectBd()
	collection := client.Database("ProjetoLTP2").Collection("Usuarios")

	return collection.InsertOne(context.Background(), usuario)

}

func InsertOneImage(imagem models.Imagem) (*mongo.InsertOneResult, error) {
	client := ConnectBd()

	collection := client.Database("ProjetoLTP2").Collection("Imagens")

	return collection.InsertOne(context.Background(), imagem)
}
