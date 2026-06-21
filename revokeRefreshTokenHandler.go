package main

import (
	"net/http"

	"github.com/thingammabob/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't retrieve token from request", err)
		return
	}
	err = cfg.queries.RevokeRefreshToken(r.Context(), tok)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
