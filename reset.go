// Handlers for resetting server state — used for testing and monitoring resets.
package main

import "net/http"

// handlerReset clears all users from the database and resets the hit counter.
// Restricted to dev environments only — returns 403 in production to prevent
// accidental data loss. This endpoint should never be exposed in a live system.
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// Platform guard ensures this destructive endpoint is inaccessible outside dev.
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	// Reset the in-memory request counter alongside the database
	// so both metrics and data are returned to a clean state together.
	cfg.fileserverHits.Store(0)

	// Delete all users — used to reset state between test runs.
	err := cfg.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete users", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset successful"))
}