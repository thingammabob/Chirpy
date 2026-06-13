package main

import (
	"database/sql"
	"fmt"
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
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Unable to establish connection with postgres database.")
		return
	}
	newConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		queries:        database.New(db),
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
