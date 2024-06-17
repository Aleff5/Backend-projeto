package utility

import (
	database "projetov2/Backend-projeto/Database"
	"projetov2/Backend-projeto/models"

	"go.mongodb.org/mongo-driver/bson"
)

func CheckAdm(usuario models.Usuario) (bool, error) {

	filter := bson.D{
		{Key: "email", Value: usuario.Email},
		{Key: "password", Value: usuario.Password},
	}

	resultado, err := database.FindOneUser(filter)

	if err != nil {
		return false, err
	}

	if resultado.Admin != true {
		return false, nil
	}
	return true, nil

}
