package main

import (
	"net/http"
)

//handlerGetAllChirps returns an array of all Chirps from the database
func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	//Retrieves chirps from the database
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
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