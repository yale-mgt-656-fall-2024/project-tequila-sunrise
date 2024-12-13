package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func init() {
	dbClient = connectToMongo("mongodb+srv://makerstate:ALHXB%5FtKQN27YZC@cluster0.o2beiyt.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
}

// connectToMongo establishes a MongoDB connection
func connectToMongo(uri string) *mongo.Client {
	// Define ServerAPIOptions
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// Set client options
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	return client
}				

// getCollection returns a reference to a MongoDB collection
func getCollection(collectionName string) *mongo.Collection {
	return dbClient.Database("eventbrite_clone").Collection(collectionName)
}
