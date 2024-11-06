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

    // Route for displaying the registration form
    r.Get("/register", registerFormController)

    return r
}
