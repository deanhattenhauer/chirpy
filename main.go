package main

import (
	"io"
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
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))

	// Readiness endpoint — used by external systems to verify the server is alive.
	// Returns 200 OK with a plain text body to confirm the server is accepting traffic.
	// Can be extended to return 503 if dependencies like the database are unavailable.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {

		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
		w.WriteHeader(http.StatusOK)

		io.WriteString(w, "OK")
	})

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