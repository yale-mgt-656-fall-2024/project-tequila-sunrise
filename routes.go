func createRoutes() chi.Router {
	r := chi.NewRouter()

	// Apply middleware
	r.Use(authMiddleware)

	// Public routes
	r.Get("/", indexController)
	r.Get("/login", loginFormController)
	r.Post("/login", loginUserController)
	r.Get("/register", registerFormController)
	r.Post("/register", registerUserController)

	// Static files
	addStaticFileServer(r, "/static/", "staticfiles")

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authRequired)
		r.Get("/events/new", newEventFormController)
		r.Post("/events/new", createNewEventController)
		r.Get("/events/{id}", eventDetailController)
		r.Post("/events/{id}/rsvp", rsvpController)
	})

	return r
}
