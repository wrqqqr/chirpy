package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) countRequestsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	//conc := "Hits: " + strconv.Itoa(cfg.fileserverHits)
	fmt.Fprintf(w, `<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`,
		cfg.fileserverHits)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}

func main() {
	const filepathRoot = "."
	const port = "8080"
	cfg := apiConfig{}

	mux := http.NewServeMux()
	mdwHandler := cfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", http.StripPrefix("/app", mdwHandler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.countRequestsHandler)
	mux.HandleFunc("GET /api/reset", cfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
