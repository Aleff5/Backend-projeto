package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Usuario struct {
	Email    string `json:"email" bson:"email"`	
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	// Id primitive.ObjectID `bson:"_id,omitempty"`
}

type Imagem struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FilePath     string             `json:"filepath" bson:"filepath"`
	Img          string             `json:"img" bson:"img"`
	Descricaoimg string             `json:"descricao" bson:"descricao"`
}
