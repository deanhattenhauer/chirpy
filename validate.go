package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// handlerValidateChirp validates incoming chirp content against Chirpy's rules.
// Chirps must be 140 characters or fewer — the same limit as early Twitter.
func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	
	// parameters maps the incoming JSON request body to a Go struct.
	// Struct tags ensure the JSON key "body" maps to the Body field.
	type parameters struct {
		Body string `json:"body"`
	}

	// errorVals is the shape of all error responses from this endpoint.
	type errorVals struct {
    	Error string `json:"error"`
	}

	// returnVals is the shape of a successful validation response.
	type returnVals struct {
        Valid bool `json:"valid"`
    }

	// Decode the request body — returns an error if JSON is malformed or wrong types.
	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)

    if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	// Enforce Chirpy's 140 character limit before processing further.
	// Returns 400 Bad Request with a descriptive error message.
	if len(params.Body) > 140 {
		errBody := errorVals{Error: "Chirp is too long"}
		dat, err := json.Marshal(errBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	// Chirp is valid — respond with 200 and a confirmation payload.
    respBody := returnVals{
    Valid: true,
	}

    dat, err := json.Marshal(respBody)
	if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}