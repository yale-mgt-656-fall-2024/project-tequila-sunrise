package main

import (
	"net/http"
	"time"
)

func indexController(w http.ResponseWriter, r *http.Request) {
	type indexContextData struct {
		Events []Event
		Today  time.Time
		User   *User
	}

	// Get all events from the database
	events, err := getAllEvents()
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
	}

	// Execute the template with the context data
	err = tmpl["index"].Execute(w, contextData)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}
