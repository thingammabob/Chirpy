package main

import (
	"encoding/json"
	"fmt"
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

	aUserRequest := userRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&aUserRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("Something went wrong decoding the request: %w", err))
		return
	}
	user, err := cfg.queries.CreateUser(r.Context(), aUserRequest.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("Something went wrong creating the user: %w", err))
		return
	}
	myUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, 201, myUser)

}
