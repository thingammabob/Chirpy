package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhookHandler(resWriter http.ResponseWriter, r *http.Request) {
	type polkaRequest struct {
		Event string `json:"event"`
		Data  struct {
			User_ID string `json:"user_id"`
		} `json:"data"`
	}
	req := polkaRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(resWriter, http.StatusBadRequest, "Couldn't decode request.", err)
		return
	}
	if req.Event != "user.upgraded" {
		respondWithError(resWriter, http.StatusNoContent, "Event is not user.upgraded", err)
		return
	}
	uuid, err := uuid.Parse(req.Data.User_ID)
	if err != nil {
		respondWithError(resWriter, http.StatusBadRequest, "Invalid user id.", err)
		return
	}
	_, err = cfg.queries.UpgradeUser(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(resWriter, http.StatusNotFound, "User not found.", err)
			return
		}
		respondWithError(resWriter, http.StatusInternalServerError, "Couldn't upgrade user.", err)
		return
	}
	respondWithJSON(resWriter, http.StatusNoContent, nil)

}
