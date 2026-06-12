package main

import (
	"net/http"
)

func (cfg *apiConfig) resetServerhits(resWriter http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	bodyBytes := []byte("Server hit counter has been reset to 0.")
	resWriter.Write(bodyBytes)
}
