package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/deanhattenhauer/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	// parameters defines the expected shape of the request body.
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	// response defines expected shape of the returned response body
	type response struct {
    	User
    	Token string `json:"token"`
	}

	// Decode the request body — returns 500 if JSON is malformed or wrong types.
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Expiration logic - not specified = 1 hour, > 1 hour = 1 hour, otherwise whatever client says
	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
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

	// Create JWT token and catch error
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to make JWT", err)
		return
	}

	// Map database.User to the API User struct before responding.
	// This decouples the JSON response shape from the internal database model.
	respondWithJSON(w, http.StatusOK, response {
		User: User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		},
		Token: token,
		})
}
