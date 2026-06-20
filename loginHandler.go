package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/thingammabob/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		RequestDuration int    `json:"expires_in_seconds"`
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
	dur := 1 * time.Hour
	if creds.RequestDuration != 0 {
		dur = time.Duration(creds.RequestDuration) * time.Second
	}
	jwt, err := auth.MakeJWT(user.ID, cfg.tokenSecret, dur)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error in creating JWT", err)
		return
	}

	type response struct {
		User
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: jwt,
	})

}
