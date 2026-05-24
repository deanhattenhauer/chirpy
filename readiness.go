package main

import "net/http"

// handlerReadiness is the health check endpoint for external systems.
// Returns 200 OK to confirm the server is alive and accepting traffic.
// Can be extended to check dependencies like the database before responding.
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
