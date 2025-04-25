package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/mharkness1/http_server_intro/internal/database"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	var err error

	authorID := r.URL.Query().Get("author_id")
	sortBy := r.URL.Query().Get("sort")

	if authorID != "" {
		authorUuid, err := uuid.Parse(authorID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "author id cannot be used", err)
			return
		}
		chirps, err = cfg.DB.GetChirpsByAuthorByCreatedASC(r.Context(), uuid.NullUUID{UUID: authorUuid, Valid: true})
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Unable to retrieve chirps", err)
			return
		}
	} else {
		chirps, err = cfg.DB.GetChirpsByCreatedAsc(r.Context())
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to retrieve chirps", err)
			return
		}

	}

	var allChirps []Chirp
	for i := range chirps {
		allChirps = append(allChirps, mapDatabaseChirpToChirp(chirps[i]))
	}

	if sortBy == "desc" {
		sort.Slice(allChirps, func(i, j int) bool {
			return allChirps[j].CreatedAt.Before(allChirps[i].CreatedAt)
		})
	}

	respondJSON(w, http.StatusOK, allChirps)
	return
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
		return
	}
	respondJSON(w, http.StatusOK, mapDatabaseChirpToChirp(chirp))
	return
}
