package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

// apiConfig holds shared server state accessible across all request handlers.
// Using a struct allows state to be injected into handlers without global variables.
type apiConfig struct {
	// atomic.Int32 ensures safe concurrent access across multiple goroutines.
	// Each incoming HTTP request runs in its own goroutine, so a regular int would race.
	fileserverHits atomic.Int32
}

func main() {
	// Centralizing configuration avoids magic strings scattered through the codebase.
	const filepathRoot = "."
	const port = "8080"
	
	// apiCfg is the single source of truth for shared server state.
	// Passed to handlers as a pointer receiver so all handlers share the same instance.
	apiCfg := apiConfig{}
	
	// ServeMux routes incoming requests to the appropriate handler.
    // Without registered routes, all requests return 404 by default.
	mux := http.NewServeMux()

	// Static assets are served under /app/ to avoid conflicts with API routes.
	// StripPrefix removes /app from the request path before the fileserver sees it,
	// so the fileserver resolves paths relative to the project root as expected.
	// Wrapped with middleware to count each fileserver request.
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// Readiness endpoint registered as a named function to keep main focused
	// on wiring and allow the handler to grow independently.
	mux.HandleFunc("/healthz", handlerReadiness)

	// Metrics and reset endpoints are methods on apiConfig to access shared state.
	// Only handlers that need state are bound to the config struct.
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)

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

// handlerMetrics reports the number of fileserver requests since last reset.
// Exposes internal server telemetry for monitoring and product analytics.
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

// middlewareMetricsInc wraps a handler to increment the request counter on each call.
// Middleware pattern allows cross-cutting concerns like metrics to be applied
// without modifying the underlying handler logic.
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}