package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) userCreateHandler(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Email string `json:"email"`
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
	user, err := cfg.queries.CreateUser(r.Context(), aUserRequest.Email)
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
