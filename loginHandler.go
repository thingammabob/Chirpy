package main

import (
	"encoding/json"
	"net/http"

	"github.com/thingammabob/chirpy/internal/auth"
	"github.com/thingammabob/chirpy/internal/database"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	creds := credentials{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&creds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request.", err)
		return
	}
	user, err := cfg.queries.GetUserFromEmail(r.Context(), creds.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials.", err)
		return
	}

	matches, err := auth.CheckPaswordHash(creds.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error in processing credentials.", err)
		return
	}
	if matches == false {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error in creating JWT", err)
		return
	}
	refreshToken := auth.MakeRefreshToken()
	cfg.queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        jwt,
		RefreshToken: refreshToken,
	})

}
