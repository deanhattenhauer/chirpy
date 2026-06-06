package main

import (
	"net/http"
	"time"

	"github.com/deanhattenhauer/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	
	// Get the token from the request headers
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to get token header", err)
		return
	}

	// Look up the user from the refresh token
	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}

	// Create new access token
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
    	Token string `json:"token"`
	}{Token: accessToken})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	
	// Get the token from the request headers
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to get token header", err)
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(),token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not refresh token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}