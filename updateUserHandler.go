package main

import (
	"encoding/json"
	"net/http"

	"github.com/thingammabob/chirpy/internal/auth"
	"github.com/thingammabob/chirpy/internal/database"
)

func (cfg *apiConfig) updateUserHandler(resWriter http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "Couldn't retreive access token from request", err)
		return
	}
	uuid, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "Invalid access token", err)
		return
	}
	type updateReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := updateReq{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "Couldn't decode request.", err)
		return
	}
	hashed_password, err := auth.HashPassword(req.Password)
	user, err := cfg.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             uuid,
		Email:          req.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "Couldn't update user.", err)
		return
	}
	type response struct {
		User
	}
	respondWithJSON(resWriter, http.StatusOK, User{
		ID:            user.ID,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Email:         user.Email,
		Is_Chirpy_Red: user.IsChirpyRed,
	})

}
