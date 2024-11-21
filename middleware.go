package main

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session")
		if err != nil {
			log.Printf("Session error: %v\n", err)
			next.ServeHTTP(w, r)
			return
		}

		userID, ok := session.Values["user_id"].(string)
		if ok && userID != "" {
			user, err := getUserByID(userID)
			if err == nil {
				// Attach user to context
				ctx := context.WithValue(r.Context(), "user", user)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func getUserByID(id string) (User, error) {
	collection := getCollection("users")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return User{}, err
	}
	var user User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	return user, err
}

func authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")
		log.Printf("AuthRequired - User in context: %+v\n", user)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
