package main

import (
	"net/http"
	"time"
)

func indexController(w http.ResponseWriter, r *http.Request) {
	type indexContextData struct {
		Events     []Event
		Today      time.Time
		User       *User
		Search     string
		IsLoggedIn bool
		Message    string
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
	IsLoggedIn := false
	ctxUser := r.Context().Value("user")
	if ctxUser != nil {
		if u, ok := ctxUser.(User); ok {
			user = &u
			IsLoggedIn = true // User is logged in
		}
	}

	// Retrieve the flash message from the session
	session, _ := store.Get(r, "session")
	message, _ := session.Values["flash"].(string)
	delete(session.Values, "flash") // Clear the flash message after retrieving it
	session.Save(r, w)

	contextData := indexContextData{
		Events:     events,
		Today:      time.Now(),
		User:       user,
		Search:     searchTerm,
		IsLoggedIn: IsLoggedIn,
		Message:    message,
	}

	// Execute the template with the context data
	err = tmpl["index"].Execute(w, contextData)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

func aboutController(w http.ResponseWriter, r *http.Request) {
	// Retrieve the IsLoggedIn flag from the context
	IsLoggedIn := false
	ctxUser := r.Context().Value("user")
	if ctxUser != nil {
		if _, ok := ctxUser.(User); ok {
			IsLoggedIn = true
		}
	}

	// Retrieve the flash message from the session
	session, _ := store.Get(r, "session")
	message, _ := session.Values["flash"].(string)
	delete(session.Values, "flash") // Clear the flash message after retrieving it
	session.Save(r, w)

	// Define data for the About page
	data := map[string]interface{}{
		"Title":      "About Us",
		"IsLoggedIn": IsLoggedIn,
		"Message":    message,
	}

	// Parse and execute the About template
	err := tmpl["about"].Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}
