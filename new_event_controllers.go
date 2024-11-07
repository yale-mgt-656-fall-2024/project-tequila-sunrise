package main

import (
	"net/http"
	"time"
)

// Renders the form to create a new event
func newEventFormController(w http.ResponseWriter, r *http.Request) {
	tmpl["new_event"].Execute(w, nil)
}

// Handles the form submission to create a new event
func createNewEventController(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	location := r.FormValue("location")
	image := r.FormValue("image")
	dateStr := r.FormValue("date")

	date, err := time.Parse("2006-01-02T15:04", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	event := Event{
		Title:    title,
		Location: location,
		Image:    image,
		Date:     date,
	}

	err = addEvent(event)
	if err != nil {
		http.Error(w, "Error adding event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
