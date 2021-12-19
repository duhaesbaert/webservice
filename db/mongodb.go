package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)



//Return true if successful to connect to Mongo DB.
func OpenConnectionToMongo() bool	{
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		return false
	}

	defer CloseConnectionToMongo(client)

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
		return false
	}

	return true
}

func CloseConnectionToMongo(client *mongo.Client) {
	client.Disconnect(context.TODO())
}