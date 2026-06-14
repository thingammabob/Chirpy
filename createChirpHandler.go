package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/thingammabob/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler(resWriter http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type response struct {
		Chirp
	}
	aChirp := chirpRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&aChirp)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "Unable to decode chirp", err)
		return
	}
	if len(aChirp.Body) > 140 {
		respondWithError(resWriter, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	cleaned_body := cleanProfanity(aChirp.Body)
	newChirp, err := cfg.queries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: aChirp.UserId,
	})
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "Unable to create chirp", err)
		return
	}

	respondWithJSON(resWriter, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        newChirp.ID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body:      newChirp.Body,
			UserID:    newChirp.UserID,
		},
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
