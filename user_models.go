package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"` // This should be a hashed password
}

// getUserByEmail - returns the user with the specified email
func getUserByEmail(email string) (User, error) {
	collection := getCollection("users")
	var user User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	return user, err
}

// addUser - Adds a new user to the collection
func addUser(user User) error {
	collection := getCollection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	return err
}

// getUserByID - returns the user with the specified ID
func getUserByID(id primitive.ObjectID) (User, error) {
    collection := getCollection("users")
    var user User
    err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
    return user, err
}


