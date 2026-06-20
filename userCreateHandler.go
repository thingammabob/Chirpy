package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thingammabob/chirpy/internal/auth"
	"github.com/thingammabob/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) userCreateHandler(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	aUserRequest := userRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&aUserRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong decoding the request.", err)
		return
	}
	hashedPassword, err := auth.HashPassword(aUserRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong trying to process the password.", err)
		return
	}
	user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          aUserRequest.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong creating the user.", err)
		return
	}
	respondWithJSON(w, 201, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})

}
