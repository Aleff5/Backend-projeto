package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectBd() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://aleff:fy31AhbHnMrpwvaP@ltp2-project.j7n4lt1.mongodb.net/?retryWrites=true&w=majority&appName=LTP2-Project"))
	if err != nil {
		panic(err)
	}
	return client
}
