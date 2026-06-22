package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/thingammabob/chirpy/internal/auth"
)

func (cfg *apiConfig) deleteChirpHandler(resWriter http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(resWriter, http.StatusBadRequest, "Couldn't parse chirpID", err)
		return
	}
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "Couldn't retrieve access token.", err)
		return
	}
	uuid, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "Not allowed!", err)
		return
	}
	chirp, err := cfg.queries.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(resWriter, http.StatusNotFound, "Couldn't retrieve chirp", err)
		return
	}
	if chirp.UserID != uuid {
		respondWithError(resWriter, http.StatusForbidden, "User not allowed to delete this chirp.", err)
		return
	}

	err = cfg.queries.DeleteChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "Couldn't delete chirp.", err)
		return
	}
	respondWithJSON(resWriter, http.StatusNoContent, nil)

}
