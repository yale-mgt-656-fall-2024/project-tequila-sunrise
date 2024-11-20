package main

import (
	"net/http"
	"regexp"
    "unicode"
	"golang.org/x/crypto/bcrypt"
)

// isValidEmail checks if the email is in a valid format
func isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
    return re.MatchString(email)
}

// isValidPassword checks if the password meets the required criteria
func isValidPassword(password string) bool {
    if len(password) < 8 {
        return false
    }
    var hasLetter, hasNumber bool
    for _, c := range password {
        switch {
        case unicode.IsLetter(c):
            hasLetter = true
        case unicode.IsNumber(c):
            hasNumber = true
        }
    }
    return hasLetter && hasNumber
}

// registerFormController renders the registration form
func registerFormController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tmpl["register"].Execute(w, nil)
}

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

    // Validate email
    if !isValidEmail(email) {
        http.Error(w, "Invalid email address", http.StatusBadRequest)
        return
    }

    // Validate password
    if !isValidPassword(password) {
        http.Error(w, "Password must be at least 8 characters long and include both letters and numbers", http.StatusBadRequest)
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

    // Redirect to login page or auto-login
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// loginFormController renders the login form
func loginFormController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tmpl["login"].Execute(w, nil)
}

// loginUserController processes the login form
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

    // Validate email
    if !isValidEmail(email) {
        http.Error(w, "Invalid email address", http.StatusBadRequest)
        return
    }

    // Validate password
    if !isValidPassword(password) {
        http.Error(w, "Invalid password format", http.StatusBadRequest)
        return
    }

    // Get user by email
    user, err := getUserByEmail(email)
    if err != nil {
        // To prevent user enumeration, use a generic error message
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Check password
    if !checkPasswordHash(password, user.Password) {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Set session cookie
    session, _ := store.Get(r, "session")
    session.Values["user_id"] = user.ID.Hex()
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

// logoutUserController logs out the user by clearing the session
func logoutUserController(w http.ResponseWriter, r *http.Request) {
    // Get the session
    session, _ := store.Get(r, "session")

    // Remove user_id from session
    delete(session.Values, "user_id")

    // Set MaxAge to -1 to delete the session cookie
    session.Options.MaxAge = -1

    // Save the session
    err := session.Save(r, w)
    if err != nil {
        http.Error(w, "Error saving session", http.StatusInternalServerError)
        return
    }

    // Redirect to the home page or login page
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
