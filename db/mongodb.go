package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



//Return true if successful to connect to Mongo DB.
func OpenConnectionToMongo() (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		return &mongo.Client{}, fmt.Errorf("Could not establish connection to MongoDB")
	}

	return client, nil
}

func CloseConnectionToMongo(client *mongo.Client) {
	client.Disconnect(context.TODO())
}