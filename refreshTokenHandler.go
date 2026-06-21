package main

import (
	"net/http"
	"time"

	"github.com/thingammabob/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshTokenHandler(resWriter http.ResponseWriter, r *http.Request) {
	tok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(resWriter, http.StatusBadRequest, "Couldn't retrieve refresh token from request", err)
		return
	}
	refTok, err := cfg.queries.RefreshRefreshToken(r.Context(), tok)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "Invalid refresh token!", err)
		return
	}
	if refTok.ExpiresAt.Before(time.Now()) {
		respondWithError(resWriter, http.StatusUnauthorized, "Refresh token has expired!", err)
		return
	}
	if refTok.RevokedAt.Valid {
		respondWithError(resWriter, http.StatusUnauthorized, "Refresh token has been revoked already!", err)
		return
	}
	jwt, err := auth.MakeJWT(refTok.UserID, cfg.tokenSecret)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "Couldn't create jwt", err)
		return
	}
	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(resWriter, http.StatusOK, response{
		Token: jwt,
	})

}
