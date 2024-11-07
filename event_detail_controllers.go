package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Displays the details of a specific event
func eventDetailController(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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
