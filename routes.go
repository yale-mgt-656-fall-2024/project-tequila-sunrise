package main

import (
    "github.com/go-chi/chi/v5"
)

func createRoutes() chi.Router {
    r := chi.NewRouter()
    r.Get("/", indexController)
    addStaticFileServer(r, "/static/", "staticfiles")

    // Routes for creating a new event
    r.Get("/events/new", newEventFormController)
    r.Post("/events/new", createNewEventController)

    // Route for event details
    r.Get("/events/{id}", eventDetailController)

    // Optional: Route for RSVP (if implementing RSVP functionality)
    r.Post("/events/{id}/rsvp", rsvpController)

    return r
}
