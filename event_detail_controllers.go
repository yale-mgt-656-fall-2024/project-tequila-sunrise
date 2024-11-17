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

	contextData := struct {
		Event Event
	}{
		Event: event,
	}

	tmpl["event_detail"].Execute(w, contextData)
}

// Handles RSVP form submission
func rsvpController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	id = strings.TrimPrefix(id, "ObjectID(")
	id = strings.TrimSuffix(id, ")")
	id = strings.Trim(id, "\"") // Remove quotes

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	err = addAttendee(id, email)
	if err != nil {
		http.Error(w, "Error adding attendee: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/events/"+id, http.StatusSeeOther)
}
