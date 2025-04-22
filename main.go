package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mharkness1/http_server_intro/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
	PLATFORM       string
}

func main() {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error opening db connection")
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()

	svrStruct := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	apiCfg := apiConfig{
		DB:       dbQueries,
		PLATFORM: platform,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)

	err = svrStruct.ListenAndServe()
	if err != nil {
		fmt.Printf("error occured: %v", err)
	}
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	type validateRequest struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := validateRequest{}
	err := decoder.Decode(&chirp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)

		type errorResponse struct {
			Error string `json:"error"`
		}

		res := errorResponse{
			Error: "Something went wrong decoding.",
		}

		jsonBody, err := json.Marshal(res)
		if err != nil {
			log.Printf("error encoding error reply: %v", err)
			return
		}

		w.Write(jsonBody)
		return
	}

	if len(chirp.Body) < 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		type validResponse struct {
			Body string `json:"cleaned_body"`
		}

		res := validResponse{
			Body: cleanChirp(chirp.Body),
		}

		jsonBody, err := json.Marshal(res)
		if err != nil {
			log.Printf("error encoding error reply: %v", err)
			return
		}

		w.Write(jsonBody)
		return

	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)

		type errorResponse struct {
			Error string `json:"error"`
		}

		res := errorResponse{
			Error: "Chirp is too long",
		}

		jsonBody, err := json.Marshal(res)
		if err != nil {
			log.Printf("error encoding error reply: %v", err)
			return
		}

		w.Write(jsonBody)
		return
	}

}
