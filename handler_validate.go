package main

import (
	"encoding/json"
	"net/http"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	aChirp := chirp{}
	err := decoder.Decode(&aChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong decoding the request.")
		return
	}
	if len(aChirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, 200, struct {
		Valid bool `json:"valid"`
	}{
		Valid: true,
	})

}
