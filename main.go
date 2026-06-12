package main

import (
	"encoding/json"
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

func validateHandler(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	aChirp := chirp{}
	err := decoder.Decode(&aChirp)
	if err != nil {
		respondWithError(w, 500, "Something went wrong decoding the request.")
		return
	}
	if len(aChirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respondWithJSON(w, 200, struct {
		Valid bool `json:"valid"`
	}{
		Valid: true,
	})

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type error struct {
		Error string `json:"error"`
	}
	w.WriteHeader(code)
	newError := error{
		Error: msg,
	}
	dat, err := json.Marshal(newError)
	if err != nil {
		fmt.Printf("Error marshalling error: %s", err)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(msg))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)

}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	w.WriteHeader(code)

	dat, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling response: %s", err)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Error in giving response"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
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
	serveMux.HandleFunc("POST /api/validate_chirp", validateHandler)
	serveMux.HandleFunc("POST /admin/reset", newConfig.resetServerhits)
	newServer.ListenAndServe()
}
