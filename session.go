package main

import (
	"os"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(getEnv("SESSION_KEY", "default-session-key")))

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
