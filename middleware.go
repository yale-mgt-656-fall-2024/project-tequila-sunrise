package main

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		userIDStr, ok := session.Values["user_id"].(string)
		if ok {
			objectID, err := primitive.ObjectIDFromHex(userIDStr)
			if err == nil {
				user, err := getUserByID(objectID)
				if err == nil {
					// Attach user to context
					ctx := context.WithValue(r.Context(), "user", user)
					r = r.WithContext(ctx)
				} else {
					// Handle error retrieving user
				}
			} else {
				// Handle error converting userID to ObjectID
			}
		}
		next.ServeHTTP(w, r)
	})
}

func authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
