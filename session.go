package main

import (
	"fmt"
	"os"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(getEnv("SESSION_KEY", "default-session-key")))

func init() {
	// Set the session options for local env
	if isLocal() {
		store.Options = &sessions.Options{
			Path:     "/",       // Cookie applies to all paths
			MaxAge:   86400 * 7, // 1 week expiration
			HttpOnly: true,      // Prevent JavaScript access
			Secure:   false,     // Set to true if using HTTPS
		}
	}
}

func isLocal() bool {
	env := os.Getenv("ENV") // Use a common ENV variable to distinguish environments
	fmt.Printf("Environment is %v\n", env)
	return env == "" || env == "local"
}
