package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type body struct {
	Body string `json:"body"`
}

type errStruct struct {
	Error string `json:"error"`
}

type validStruct struct {
	Valid bool `json:"valid"`
}

type cleanOutput struct {
	Clean string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	es := errStruct{
		Error: msg,
	}
	w.WriteHeader(code)
	w.Write([]byte(es.Error))
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	dat, err := json.Marshal(payload)

	if err != nil {
		respondWithError(w, 400, "error when marshaling")
	}
	w.Write(dat)

}

func validateHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	b := body{}
	err := decoder.Decode(&b)

	defer r.Body.Close()

	fmt.Println(b)

	if err != nil {

		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(b.Body) > 140 {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	v := cleanOutput{
		Clean: badWordReplacement(b.Body),
	}

	//dat, err := json.Unmarshal(v)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		respondWithError(w, 500, "Error marshalling JSON")
		return
	}

	respondWithJSON(w, 200, v)

}

func badWordReplacement(str string) string {
	badWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	replacer := "****"
	k := strings.Split(str, " ")
	for index, word := range k {
		_, ok := badWords[strings.ToLower(word)]
		if ok {
			k[index] = replacer
		}
	}
	return strings.Join(k, " ")
}
