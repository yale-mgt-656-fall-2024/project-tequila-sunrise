package main

import (
	"net/http"
	"html/template"
	"time"
)

func indexController(w http.ResponseWriter, r *http.Request) {
	type indexContextData struct {
		Events []Event
		Today  time.Time
		User   *User
		Search string
	}

	// Get the search term from the query parameters
	searchTerm := r.URL.Query().Get("search")

	var events []Event
	var err error
	if searchTerm != "" {
		// Get events that match the search term
		events, err = searchEvents(searchTerm)
	} else {
		// Get all events from the database
		events, err = getAllEvents()
	}
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the user from the context
	var user *User
	ctxUser := r.Context().Value("user")
	if ctxUser != nil {
		if u, ok := ctxUser.(User); ok {
			user = &u
		}
	}

	contextData := indexContextData{
		Events: events,
		Today:  time.Now(),
		User:   user,
		Search: searchTerm,
	}

	// Execute the template with the context data
	err = tmpl["index"].Execute(w, contextData)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

func aboutController(w http.ResponseWriter, r *http.Request) {
	// Define data for the About page
	data := map[string]string{
		"Title": "About Us",
	}

	// Parse and execute the About template
	err := tmpl["about"].Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}