package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	type error struct {
		Error string `json:"error"`
	}
	newError := error{
		Error: msg,
	}
	respondWithJSON(w, code, newError)

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
