package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Store(cfg.fileserverHits.Add(1))
		next.ServeHTTP(w, r)
	})
}

func healthyHandler(resWriter http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(http.StatusOK)
	bodyBytes := []byte("OK")
	resWriter.Write(bodyBytes)
}
func (cfg *apiConfig) serverhitsHandler(resWriter http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/html")
	resWriter.WriteHeader(http.StatusOK)
	htmlContent := fmt.Sprintf(`
		<html>
		  <body>
		    <h1>Welcome, Chirpy Admin</h1>
		    <p>Chirpy has been visited %d times!</p>
		  </body>
		</html>
		`, cfg.fileserverHits.Load())
	bodyBytes := []byte(htmlContent)
	resWriter.Write(bodyBytes)
}
func (cfg *apiConfig) resetServerhits(resWriter http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	bodyBytes := []byte("Server hit counter has been reset to 0.")
	resWriter.Write(bodyBytes)
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
	serveMux.HandleFunc("GET /api/healthz", healthyHandler)
	serveMux.HandleFunc("POST /admin/reset", newConfig.resetServerhits)
	newServer.ListenAndServe()
}
