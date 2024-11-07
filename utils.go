package main

import "os"

// getEnv retrieves environment variables or returns a fallback value
func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
