// Chirpy — a social network API server
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/deanhattenhauer/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// apiConfig holds shared server state accessible across all request handlers.
// Using a struct allows state to be injected into handlers without global variables.
type apiConfig struct {
	// atomic.Int32 ensures safe concurrent access across multiple goroutines.
	// Each incoming HTTP request runs in its own goroutine, so a regular int would race.
	fileserverHits atomic.Int32
	// dbQueries provides type-safe access to all database operations via SQLC.
	dbQueries *database.Queries
	// platform controls environment-specific behavior like the admin reset endpoint.
	platform string
	//jwt secret for token creation
	jwtSecret string
}

func main() {
	// Load environment variables from .env file before reading any config.
	// Must be called before os.Getenv to ensure variables are available.
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// SQLC-generated query wrapper provides type-safe database access.
	dbQueries := database.New(db)

	// Platform controls environment-specific behavior.
	// Set to "dev" locally to enable dangerous endpoints like /admin/reset.
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")

	// Centralizing configuration avoids magic strings scattered through the codebase.
	const filepathRoot = "."
	const port = "8080"

	// apiCfg is the single source of truth for shared server state.
	// Passed to handlers as a pointer receiver so all handlers share the same instance.
	apiCfg := apiConfig{
		dbQueries: dbQueries,
		platform:  platform,
		jwtSecret: jwtSecret,
	}

	// ServeMux routes incoming requests to the appropriate handler.
	// Without registered routes, all requests return 404 by default.
	mux := http.NewServeMux()

	// Static assets are served under /app/ to avoid conflicts with API routes.
	// StripPrefix removes /app from the request path before the fileserver sees it,
	// so the fileserver resolves paths relative to the project root as expected.
	// Wrapped with middleware to count each fileserver request.
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	// Readiness endpoint registered as a named function to keep main focused
	// on wiring and allow the handler to grow independently.
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Chirp endpoints
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}" ,apiCfg.handlerGetSingleChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	// Token endpoints
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)

	// User endpoints
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)
	
	// Login endpoints
	mux.HandleFunc("POST /api/login" ,apiCfg.handlerUserLogin)

	// Admin endpoints are restricted by platform — dangerous in production.
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// The mux is injected as the handler so all routing decisions
	// flow through a single, centrally managed router.
	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Logged before blocking so the operator knows the server is ready.
	// Code after ListenAndServe only executes on shutdown or error.
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	// ErrServerClosed is not a real error — it signals a clean shutdown.
	// Any other error indicates an unexpected failure and is fatal.
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}