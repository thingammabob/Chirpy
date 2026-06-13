package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) resetServerhits(resWriter http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		resWriter.WriteHeader(http.StatusForbidden)
		return
	}
	err := cfg.queries.DeleteAllUsers(req.Context())
	if err != nil {
		fmt.Println(err)
		respondWithError(resWriter, http.StatusInternalServerError, "Unable to delete all users")
		return
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	bodyBytes := []byte("Reset successful!")
	resWriter.Write(bodyBytes)
}
