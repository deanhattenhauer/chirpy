package main

import (
	"encoding/json"
	"net/http"

	"github.com/deanhattenhauer/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	// parameters defines the expected shape of the request body.
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	// Decode the request body — returns 500 if JSON is malformed or wrong types.
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Look up the user by email first to get their stored hash
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Compare password and hash
	passwordsMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !passwordsMatch {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Map database.User to the API User struct before responding.
	// This decouples the JSON response shape from the internal database model.
	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
