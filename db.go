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
    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb+srv://zpeterli523:<db_password>@cluster0.o2beiyt.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)

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
    dbClient = client
}

// Get a handle for your collection
func getCollection(collectionName string) *mongo.Collection {
    return dbClient.Database("eventbrite_clone").Collection(collectionName)
}
