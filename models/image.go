package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Imagem struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Filename    string             `json:"filename" bson:"filename"`
	FileUrl     string             `json:"fileurl" bson:"fileurl"`
	Description string             `json:"description" bson:"description"`
}
