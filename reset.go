// Handlers for resetting server state — used for testing and monitoring resets.
package main

import "net/http"

// handlerReset clears the fileserver hit counter back to zero.
// Useful for resetting metrics between test runs or monitoring windows.
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}