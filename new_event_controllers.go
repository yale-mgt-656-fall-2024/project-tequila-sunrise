package main

import (
	"net/http"
	"time"
)

// Renders the form to create a new event
func newEventFormController(w http.ResponseWriter, r *http.Request) {
	// Extract `IsLoggedIn` from the context
	isLoggedIn := r.Context().Value("isLoggedIn").(bool)

	// Retrieve the flash message from the session
	session, _ := store.Get(r, "session")
	message, _ := session.Values["flash"].(string)
	delete(session.Values, "flash") // Clear the flash message after retrieving it
	session.Save(r, w)

	// Pass `IsLoggedIn` and `Message` to the template
	err := tmpl["new_event"].Execute(w, map[string]interface{}{
		"IsLoggedIn": isLoggedIn,
		"Message":    message,
	})
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

// Handles the form submission to create a new event
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

	// Set a flash message for successful event creation
	session, _ := store.Get(r, "session")
	session.Values["flash"] = "Successfully created new event: " + title
	session.Save(r, w)

	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
