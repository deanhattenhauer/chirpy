package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUserToRed(w http.ResponseWriter, r *http.Request) {
	
	// parameters defines the expected shape of the request body.
	type parameters struct {
    	Event string `json:"event"`
    	Data  struct {
        	UserID uuid.UUID `json:"user_id"`
    	} `json:"data"`
	}

	// Decode the request body — returns 500 if JSON is malformed or wrong types.
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Check if the event is "user.upgraded" - if not, respond with 204 and return
	if params.Event != 	"user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Call UpdateUserToChirpyRed with the user ID from params, handle the error, and respond with 204 on success
	_, err = cfg.dbQueries.UpdateUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)


}
