package main

import (
	"net/http"
)

func healthyHandler(resWriter http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(http.StatusOK)
	bodyBytes := []byte("OK")
	resWriter.Write(bodyBytes)
}
