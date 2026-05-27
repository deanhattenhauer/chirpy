package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondWithError writes a JSON error response with the given status code and message.
// 5XX errors are logged as server-side failures — 4XX errors are expected client mistakes.
// Pass a non-nil err to log the underlying cause for server errors.
func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

// respondWithJSON marshals the payload to JSON and writes it to the response.
// All JSON responses flow through here to ensure consistent headers and encoding.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}