package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/deanhattenhauer/chirpy/internal/auth"
	"github.com/deanhattenhauer/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	// parameters defines the expected shape of the request body.
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	// response defines expected shape of the returned response body
	type response struct {
    	User
    	Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	// Create JWT token and catch error
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to make JWT", err)
		return
	}

	// Create refreshToken for response
	refreshToken := auth.MakeRefreshToken()
	_ , err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: refreshToken, CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, ExpiresAt:  time.Now().Add(60 * 24 * time.Hour)})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create refresh token", err)
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
		RefreshToken: refreshToken,
		})
}
