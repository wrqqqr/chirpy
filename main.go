package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

	const headersContent = "text/plain; charset=utf-8"
	const headersKey = "Content-Type"
	const stsCode = 200
	bdyText := []byte("OK")

	w.Header().Add(headersKey, headersContent)
	w.WriteHeader(stsCode)
	w.Write(bdyText)
}

func main() {
	const port = "8080"
	const fileServerPath = "/app/*"
	const filepathRoot = "."
	const trafficPath = "/healthz"

	mux := http.NewServeMux()

	mux.Handle(fileServerPath, http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	mux.HandleFunc(trafficPath, handler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
