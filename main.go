package main

import "net/http"

func main() {
	serveMux := http.NewServeMux()
	newServer := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	newServer.ListenAndServe()
}
