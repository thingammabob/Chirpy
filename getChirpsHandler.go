package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/thingammabob/chirpy/internal/database"
)

func (cfg *apiConfig) getChirpsHandler(resWriter http.ResponseWriter, r *http.Request) {
	author_id := r.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error
	if author_id == "" {
		chirps, err = cfg.queries.GetChirps(r.Context())
	} else {
		uuid, errParse := uuid.Parse(author_id)
		if errParse != nil {
			respondWithError(resWriter, http.StatusBadRequest, "Couldn't parse user id.", errParse)
			return
		}
		chirps, err = cfg.queries.GetChirpByAuthor(r.Context(), uuid)
	}

	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "Couldn't retrieve chirps.", err)
		return
	}
	if len(chirps) == 0 {
		respondWithError(resWriter, http.StatusNotFound, "No chirps available.", err)
		return
	}

	myChirps := []Chirp{}
	for _, chirp := range chirps {
		myChirps = append(myChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(resWriter, http.StatusOK, myChirps)

}

func (cfg *apiConfig) getAChirpHandler(resWriter http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(resWriter, http.StatusBadRequest, "Invalid id format", err)
		return
	}
	chirp, err := cfg.queries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(resWriter, http.StatusNotFound, "Chirp not found.", err)
		return
	}

	respondWithJSON(resWriter, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}
