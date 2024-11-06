package main

import (
	"net/http"
	"time"
)

func indexController(w http.ResponseWriter, r *http.Request) {

	type indexContextData struct {
		Events []Event
		Today  time.Time
	}

	theEvents, err := getAllEvents()
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	contextData := indexContextData{
		Events: theEvents,
		Today:  time.Now(),
	}

	tmpl["index"].Execute(w, contextData)
}

func newEventFormController(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    tmpl["new_event"].Execute(w, nil)
}

func createNewEventController(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Error parsing form data", http.StatusBadRequest)
        return
    }

    title := r.FormValue("title")
    location := r.FormValue("location")
    image := r.FormValue("image")
    dateStr := r.FormValue("date")

    // Parse the date and time from the form
    date, err := time.Parse("2006-01-02T15:04", dateStr)
    if err != nil {
        http.Error(w, "Invalid date format", http.StatusBadRequest)
        return
    }

    // Create a new event
    event := Event{
        Title:    title,
        Location: location,
        Image:    image,
        Date:     date,
    }

    // Add the event to the list (using your existing addEvent function)
    addEvent(event)

    // Redirect to the index page or event detail page
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
