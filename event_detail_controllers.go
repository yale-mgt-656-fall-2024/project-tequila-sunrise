package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
)

// Debug function to list all events
func listAllEvents() []Event {
	collection := getCollection("events")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("Error fetching events:", err)
		return nil
	}
	defer cursor.Close(context.TODO())

	var events []Event
	if err = cursor.All(context.TODO(), &events); err != nil {
		fmt.Println("Error decoding events:", err)
		return nil
	}

	// Debug: Print all events in DB
	fmt.Println("All events in DB:")
	for _, e := range events {
		fmt.Printf("Event Details:\n%+v\n\n", e)
	}
	return events
}

// Displays the details of a specific event
func eventDetailController(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Clean the ID
	id = strings.TrimPrefix(id, "ObjectID(")
	id = strings.TrimSuffix(id, ")")
	id = strings.Trim(id, "\"")

	event, found := getEventByID(id)
	if !found {
		http.NotFound(w, r)
		return
	}

	// Retrieve the flash message from the session
	session, _ := store.Get(r, "session")
	message, _ := session.Values["flash"].(string)
	delete(session.Values, "flash") // Clear the flash message
	session.Save(r, w)

	// Extract IsLoggedIn from the context
	isLoggedIn := r.Context().Value("isLoggedIn").(bool)

	// Context data for the template
	contextData := struct {
		Event      Event
		IsLoggedIn bool
		Message    string
	}{
		Event:      event,
		IsLoggedIn: isLoggedIn,
		Message:    message,
	}

	// Execute the template with context data
	err := tmpl["event_detail"].Execute(w, contextData)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

// Handles RSVP submission
func rsvpController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the event ID from the URL
	id := chi.URLParam(r, "id")
	id = strings.TrimPrefix(id, "ObjectID(")
	id = strings.TrimSuffix(id, ")")
	id = strings.Trim(id, "\"") // Remove quotes

	// Check if the user is logged in
	ctxUser := r.Context().Value("user")
	if ctxUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, ok := ctxUser.(User)
	if !ok {
		http.Error(w, "Invalid user context", http.StatusInternalServerError)
		return
	}

	// Use the logged-in user's email
	email := user.Email

	// Add the user as an attendee
	registered, err := addAttendee(id, email)
	if err != nil {
		http.Error(w, "Error adding attendee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the message based on the RSVP status
	message := "Successfully registered for the event."
	if !registered {
		message = "You are already registered for this event."
	}

	// Store the flash message in the session
	session, _ := store.Get(r, "session")
	session.Values["flash"] = message
	session.Save(r, w)

	// Redirect to the event details page
	http.Redirect(w, r, "/events/"+id, http.StatusSeeOther)
}
