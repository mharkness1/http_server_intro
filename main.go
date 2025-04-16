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
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	svrStruct := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	var apiCfg apiConfig

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)

	err := svrStruct.ListenAndServe()
	if err != nil {
		fmt.Printf("error occured: %v", err)
	}

}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	counter := cfg.fileserverHits.Load()
	html := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", counter)
	w.Write([]byte(html))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	fmt.Fprintf(w, "Counter reset to 0")
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	type validateRequest struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	body := validateRequest{}
	err := decoder.Decode(&body)
	if err != nil {
		fmt.Fprintf(w, "err decoding json: %v", err)
	}

}
