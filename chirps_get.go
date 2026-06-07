package main

import (
	"database/sql"
	"net/http"

	"github.com/deanhattenhauer/chirpy/internal/database"
	"github.com/google/uuid"
)

//handlerGetAllChirps returns an array of all Chirps from the database
func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	// Check for optional author_id
	authorID := r.URL.Query().Get("author_id")

	var chirps []database.Chirp
	var err error
	
	// If authorID is not empty, parse it as UUID and call GetAllChirpsSingleUser
	// If authorID is empty, call GetAllChirps
	if authorID != "" {
		parsedID, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", parseErr)
			return
		}
		chirps, err = cfg.dbQueries.GetAllChirpsSingleUser(r.Context(), parsedID)
	} else {
    	chirps, err = cfg.dbQueries.GetAllChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't not get author id", err)
		return
	}

	newSlice := []Chirp{}
	
	//Map database.Chirp to API Chirp struct to ensure correct JSON key casing via struct tags.
	for _, chirp := range chirps {
		newSlice = append(newSlice, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:     chirp.Body,
		UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, newSlice)
}

//handlerGetSingleChirp return a single chirp by ID
func (cfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {

	//Get ID from path
	id := r.PathValue("chirpID")

	//Convert path from string to uuid
	path, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't not parse string", err)
		return
	}
	
	//Retrieves chirps from the database
	chirpID, err := cfg.dbQueries.GetSingleChirp(r.Context(), path)
	if err == sql.ErrNoRows {
    	respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	} else if err != nil {
    	respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirpID.ID,
		CreatedAt: chirpID.CreatedAt,
		UpdatedAt: chirpID.UpdatedAt,
		Body:     chirpID.Body,
		UserID: chirpID.UserID,
		})
}