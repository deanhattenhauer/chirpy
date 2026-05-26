package main

import (
	"encoding/json"
	"net/http"
)

// handlerChirpsValidate validates incoming chirp content against Chirpy's character limit.
// Chirps must be 140 characters or fewer — returns 400 if exceeded, 200 if valid.
func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	
	// parameters maps the incoming JSON request body to a Go struct.
	// Struct tags ensure the JSON key "body" maps to the Body field.
	type parameters struct {
		Body string `json:"body"`
	}

	// returnVals is the shape of a successful validation response.
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	// Decode the request body — invalid JSON or wrong types return a 500.
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Named constant avoids magic numbers and makes the limit easy to change.
	const maxChirpLength = 140

	// Enforce the character limit before processing further.
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// Chirp is valid — respond with confirmation payload.
	respondWithJSON(w, http.StatusOK, returnVals{
		Valid: true,
	})
}
