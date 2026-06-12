package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	newConfig := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	serveMux := http.NewServeMux()
	newServer := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	serveMux.Handle("/app/", newConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /admin/metrics", newConfig.serverhitsHandler)
	serveMux.HandleFunc("POST /admin/reset", newConfig.resetServerhits)
	serveMux.HandleFunc("GET /api/healthz", healthyHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", validateHandler)
	newServer.ListenAndServe()
}
