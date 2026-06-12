package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

	cleaned_body := cleanProfanity(aChirp.Body)

	respondWithJSON(w, 200, struct {
		Cleaned_body string `json:"cleaned_body"`
	}{
		Cleaned_body: cleaned_body,
	})

}

func cleanProfanity(body string) string {
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(body, " ")
	for i, word := range words {
		if _, ok := profaneWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")

}
