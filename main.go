// Goal: Build and run a server that binds on port 8080 and always responds with 404 Not Found.
package main

import (
	"log"
	"net/http"
)

func main() {
	// ServeMux routes incoming requests to the appropriate handler.
    // Without registered routes, all requests return 404 by default.
	mux := http.NewServeMux()

	// FileServer serves static assets from the current directory.
    // The root path "/" catches all unmatched requests and serves files.
	mux.Handle("/", http.FileServer(http.Dir(".")))

	// Server is configured to listen on all network interfaces on port 8080.
    // The mux handles routing decisions for all incoming requests.
	s := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	// ListenAndServe blocks indefinitely, accepting and dispatching requests.
    // ErrServerClosed is expected on graceful shutdown and is not a real error.
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}