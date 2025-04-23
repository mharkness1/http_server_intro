package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChripsByCreatedAsc(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve chirps", err)
		return
	}

	var allChirps []Chirp
	for i, _ := range chirps {
		allChirps = append(allChirps, mapDatabaseChirpToChirp(chirps[i]))
	}

	respondJSON(w, http.StatusOK, allChirps)
}

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUuid, err := uuid.Parse(chirpID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to decode chirp ID", err)
		return
	}
	chirp, err := cfg.DB.GetChirpByID(r.Context(), chirpUuid)
	if err != nil {
		respondError(w, http.StatusNotFound, "No matching chirp", err)
	}
	respondJSON(w, http.StatusOK, mapDatabaseChirpToChirp(chirp))
}
