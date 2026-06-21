package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thingammabob/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
	platform       string
	tokenSecret    string
}

func main() {
	const port = "8080"
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	tokenSecret := os.Getenv("secret")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
		return
	}
	newConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		queries:        database.New(db),
		platform:       platform,
		tokenSecret:    tokenSecret,
	}
	serveMux := http.NewServeMux()
	newServer := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	serveMux.Handle("/app/", newConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /admin/metrics", newConfig.serverhitsHandler)
	serveMux.HandleFunc("POST /admin/reset", newConfig.resetServerhits)
	serveMux.HandleFunc("GET /api/healthz", healthyHandler)
	serveMux.HandleFunc("POST /api/users", newConfig.userCreateHandler)
	serveMux.HandleFunc("POST /api/chirps", newConfig.createChirpHandler)
	serveMux.HandleFunc("GET /api/chirps", newConfig.getChirpsHandler)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", newConfig.getAChirpHandler)
	serveMux.HandleFunc("POST /api/login", newConfig.loginHandler)
	serveMux.HandleFunc("POST /api/refresh", newConfig.refreshTokenHandler)
	serveMux.HandleFunc("POST /api/revoke", newConfig.revokeRefreshTokenHandler)
	serveMux.HandleFunc("PUT /api/users", newConfig.updateUserHandler)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", newConfig.deleteChirpHandler)
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(newServer.ListenAndServe())
}
