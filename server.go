package main

import (
    "log"
    "net/http"
)

func main() {
    // Create the router
    r := createRoutes()

    // Get the port from the environment variable or default to 8080
    port := getEnv("PORT", "8080")

    // Start the server
    log.Println("Server is running on port " + port)
    err := http.ListenAndServe(":"+port, r)
    if err != nil {
        log.Fatal(err)
    }
}
