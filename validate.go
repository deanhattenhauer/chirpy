package main

import (
	"encoding/json"
	"net/http"
	"strings"
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
    	CleanedBody string `json:"cleaned_body"`
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

	// Map lookup is O(1) — more efficient than slice search for word filtering.
	// Passed to getCleanedBody so the bad word list can vary per call if needed.
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned,
	})
}

// getCleanedBody replaces profane words with asterisks while preserving case
// of surrounding words. Comparison is case-insensitive but punctuation is respected —
// "sharbert!" is not considered a match.
func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")

	// Modify words in place by index — more memory efficient than building a new slice.
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
