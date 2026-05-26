package main

import (
	"fmt"
	"net/http"
)

// handlerMetrics reports the number of fileserver requests since last reset.
// Exposes internal server telemetry for monitoring and product analytics.
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("text/html", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
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
