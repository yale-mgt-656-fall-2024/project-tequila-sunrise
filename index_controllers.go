package main

import (
	"net/http"
	"time"
    "strconv"
    "github.com/go-chi/chi/v5"
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

func eventDetailController(w http.ResponseWriter, r *http.Request) {
    // Extract the event ID from the URL
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid event ID", http.StatusBadRequest)
        return
    }

    // Retrieve the event by ID
    event, found := getEventByID(id)
    if !found {
        http.NotFound(w, r)
        return
    }

    // Prepare the context data
    type eventDetailContextData struct {
        Event Event
    }

    contextData := eventDetailContextData{
        Event: event,
    }

    // Render the template
    tmpl["event_detail"].Execute(w, contextData)
}

func rsvpController(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Extract the event ID from the URL
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid event ID", http.StatusBadRequest)
        return
    }

    // Parse form data
    err = r.ParseForm()
    if err != nil {
        http.Error(w, "Error parsing form data", http.StatusBadRequest)
        return
    }

    email := r.FormValue("email")
    if email == "" {
        http.Error(w, "Email is required", http.StatusBadRequest)
        return
    }

    // Add attendee to the event
    err = addAttendee(id, email)
    if err != nil {
        http.Error(w, "Error adding attendee: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Redirect back to the event detail page
    http.Redirect(w, r, "/events/"+idStr, http.StatusSeeOther)
}
