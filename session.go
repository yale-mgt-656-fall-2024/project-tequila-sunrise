package main

import (
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(getEnv("SESSION_KEY", "default-session-key")))
