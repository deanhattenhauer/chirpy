package main

import (
	"log"
	"net/http"
)

func main() {
	// Centralizing configuration avoids magic strings scattered through the codebase.
	const filepathRoot = "."
	const port = "8080"

	// ServeMux routes incoming requests to the appropriate handler.
    // Without registered routes, all requests return 404 by default.
	mux := http.NewServeMux()

	// Static assets are served under /app/ to avoid conflicts with API routes.
	// StripPrefix removes /app from the request path before the fileserver sees it,
	// so the fileserver resolves paths relative to the project root as expected.
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	// Readiness endpoint registered as a named function to keep main focused
	// on wiring and allow the handler to grow independently.
	mux.HandleFunc("/healthz", handlerReadiness)

	// Server is configured to listen on all network interfaces on port 8080.
    // The mux handles routing decisions for all incoming requests.
	s := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	// Logged before blocking so the operator knows the server is ready.
	// Code after ListenAndServe only executes on shutdown or error.
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	// ListenAndServe blocks indefinitely, accepting and dispatching requests.
    // ErrServerClosed is expected on graceful shutdown and is not a real error.
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

// handlerReadiness is the health check endpoint for external systems.
// Returns 200 OK to confirm the server is alive and accepting traffic.
// Can be extended to check dependencies like the database before responding.
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}