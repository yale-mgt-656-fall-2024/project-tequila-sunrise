package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// registerFormController renders the registration form
func registerFormController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is already logged in
	isLoggedIn := r.Context().Value("isLoggedIn").(bool)
	if isLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve the flash message from the session
	session, _ := store.Get(r, "session")
	message, _ := session.Values["flash"].(string)
	delete(session.Values, "flash") // Clear the flash message after retrieving it
	session.Save(r, w)

	// Render the registration form with the message
	err := tmpl["register"].Execute(w, map[string]interface{}{
		"Message":    message,
		"IsLoggedIn": isLoggedIn,
	})
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

// registerUserController processes the registration form
// registerUserController processes the registration form
func registerUserController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Check if the email is already registered
	_, err = getUserByEmail(email)
	if err == nil {
		// User already exists, redirect back to the register page with a message
		session, _ := store.Get(r, "session")
		session.Values["flash"] = "Email is already registered."
		session.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// Hash the password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := User{
		Email:    email,
		Password: hashedPassword,
	}

	// Add the user to the database
	err = addUser(user)
	if err != nil {
		http.Error(w, "Error adding user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set flash message for successful registration
	session, _ := store.Get(r, "session")
	session.Values["flash"] = "Account successfully created. Please log in."
	session.Save(r, w)

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// loginFormController renders the login form
func loginFormController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is already logged in
	isLoggedIn := r.Context().Value("isLoggedIn").(bool)
	if isLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Retrieve the flash message from the session
	session, _ := store.Get(r, "session")
	message, _ := session.Values["flash"].(string)
	delete(session.Values, "flash") // Clear the flash message after retrieving it
	session.Save(r, w)

	// Render the login form with the message
	err := tmpl["login"].Execute(w, map[string]interface{}{
		"Message":    message,
		"IsLoggedIn": isLoggedIn,
	})
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

// loginUserController processes the login form
func loginUserController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Get user by email
	user, err := getUserByEmail(email)
	if err != nil {
		// Email not found, set flash message and redirect to login
		session, _ := store.Get(r, "session")
		session.Values["flash"] = "Invalid email or password."
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Check password
	if !checkPasswordHash(password, user.Password) {
		// Password incorrect, set flash message and redirect to login
		session, _ := store.Get(r, "session")
		session.Values["flash"] = "Invalid email or password."
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Set session cookie for successful login
	session, _ := store.Get(r, "session")
	session.Values["user_id"] = user.ID.Hex()
	session.Values["flash"] = "Successfully logged in."
	session.Save(r, w)

	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HashPassword hashes the plain text password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares the hashed password with the plain text password
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// logoutController destroys the session and redirects to the home page
func logoutController(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	// Check if the user is logged in
	if _, ok := session.Values["user_id"]; !ok {
		// User is not logged in
		session.Values["flash"] = "User not logged in."
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// If user is logged in, log them out
	delete(session.Values, "user_id")
	session.Values["flash"] = "Successfully logged out."
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
