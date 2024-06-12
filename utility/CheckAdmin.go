package utility

import (
	"projetov2/Backend-projeto/Database"
	"projetov2/Backend-projeto/models"

	"go.mongodb.org/mongo-driver/bson"
)

func CheckAdm(usuario models.Usuario) (bool, error) {

	filter := bson.D{
		{"email", usuario.Email},
		{"password", usuario.Password},
	}

	resultado, err := database.FindOne(filter)

	if err != nil {
		return false, err
	}

	if resultado.Admin != true {
		return false, nil
	}
	return true, nil

}
