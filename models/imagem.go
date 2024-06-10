package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Imagem struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FilePath     string             `json:"filepath" bson:"filepath"`
	Img          string             `json:"img" bson:"img"`
	Descricaoimg string             `json:"descricao" bson:"descricao"`
}
