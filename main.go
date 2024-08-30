package main

import (
	"log"
	"net/http"
	"server/database"
)

type apiConfig struct {
	fileserverHits int
	nextId         int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	cfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	//databaseData := database.DB{}

	mux := http.NewServeMux()
	mdwHandler := cfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", http.StripPrefix("/app", mdwHandler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.countRequestsHandler)
	mux.HandleFunc("GET /api/reset", cfg.resetHandler)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsRetrieve)
	database.NewDB("database.json")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
