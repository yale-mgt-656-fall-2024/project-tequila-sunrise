package main

import (
    "net/http"
)

func registerFormController(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    tmpl["register"].Execute(w, nil)
}

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

func loginFormController(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    tmpl["login"].Execute(w, nil)
}

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
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Check password
    if !checkPasswordHash(password, user.Password) {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Set session cookie (we'll cover this next)
    session, _ := store.Get(r, "session")
    session.Values["user_id"] = user.ID.Hex()
    session.Save(r, w)

    // Redirect to the home page
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

