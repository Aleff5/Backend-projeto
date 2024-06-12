package models

type Usuario struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Admin    bool   `json:"admin" bson:"admin"`
	// Id primitive.ObjectID `bson:"_id,omitempty"`
}
