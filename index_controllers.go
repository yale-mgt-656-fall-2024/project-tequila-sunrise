package main

import (
	"net/http"
	"time"
)

type EventWithCreator struct {
	Event        Event
	CreatorEmail string
}

func indexController(w http.ResponseWriter, r *http.Request) {
    type indexContextData struct {
        Events []EventWithCreator
        Today  time.Time
        User   *User
        Filter string
    }

    // Retrieve the user from the context
    var user *User
    ctxUser := r.Context().Value("user")
    if ctxUser != nil {
        if u, ok := ctxUser.(User); ok {
            user = &u
        }
    }

    // Get the filter parameter
    filter := r.URL.Query().Get("filter")
    if filter == "" {
        filter = "all"
    }

    // Get events based on filter
    var events []Event
    var err error
    if filter == "my" && user != nil {
        events, err = getEventsByUser(user.ID)
    } else {
        events, err = getAllEvents()
    }

    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // For each event, fetch the creator's email
    eventsWithCreators := make([]EventWithCreator, 0, len(events))
    for _, event := range events {
        creatorEmail := "Unknown"
        if !event.CreatedBy.IsZero() {
            creator, err := getUserByID(event.CreatedBy)
            if err == nil {
                creatorEmail = creator.Email
            }
        }
        eventsWithCreators = append(eventsWithCreators, EventWithCreator{
            Event:        event,
            CreatorEmail: creatorEmail,
        })
    }

    contextData := indexContextData{
        Events: eventsWithCreators,
        Today:  time.Now(),
        User:   user,
        Filter: filter,
    }

    // Execute the template with the context data
    err = tmpl["index"].Execute(w, contextData)
    if err != nil {
        http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
    }
}
