package main

import (
	"net/http"
)

func (cfg *apiConfig) resetServerhits(resWriter http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		resWriter.WriteHeader(http.StatusForbidden)
		resWriter.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	err := cfg.queries.DeleteAllUsers(req.Context())
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "Unable to delete all users: "+err.Error(), err)
		return
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	bodyBytes := []byte("Reset successful!")
	resWriter.Write(bodyBytes)
}
