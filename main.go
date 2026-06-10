package main

import (
	"net/http"
)

func healthyHandler(resWriter http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	resWriter.WriteHeader(200)
	bodyBytes := []byte("OK")
	resWriter.Write(bodyBytes)
}

func main() {
	serveMux := http.NewServeMux()
	newServer := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	serveMux.HandleFunc("/healthz/", healthyHandler)
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	newServer.ListenAndServe()
}
