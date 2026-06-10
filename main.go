package main

import "net/http"

func main() {
	serveMux := http.NewServeMux()
	newServer := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	serveMux.Handle("/", http.FileServer(http.Dir(".")))
	newServer.ListenAndServe()
}
