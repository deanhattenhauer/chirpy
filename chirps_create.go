package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/deanhattenhauer/chirpy/internal/database"
	"github.com/google/uuid"
)

// Chirp is the API representation of a Chirp.
// Mapped from database.Chirp to control JSON key names via struct tags.
// Keeping this separate from the database model allows the API shape
// to evolve independently of the database schema.
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

// handlerCreateChirp creates a new chirp.
// Returns 201 Created with the full chirp object on success.
func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	// parameters defines the expected shape of the request body.
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	// Decode the request body — returns 500 if JSON is malformed or wrong types.
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	
	// Verify chirp is valid before mapping
	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	// Persist the chirp to the database.
	// This decouples the JSON response shape from the internal database model.
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Body: cleaned, UserID: params.UserID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:     chirp.Body,
		UserID: chirp.UserID,
	})

}
func validateChirp(body string) (string, error){
	// Named constant avoids magic numbers and makes the limit easy to change.
	const maxChirpLength = 140
	
	// Enforce the character limit before processing further.
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	// Map lookup is O(1) — more efficient than slice search for word filtering.
	// Passed to getCleanedBody so the bad word list can vary per call if needed.
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
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