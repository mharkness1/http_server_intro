package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mharkness1/http_server_intro/internal/auth"
	"github.com/mharkness1/http_server_intro/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type createChirpRequest struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "token extraction failed", err)
		return
	}
	userFromToken, err := auth.ValidateJWT(token, cfg.SECRET)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "token failed validation", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirp := createChirpRequest{}
	err = decoder.Decode(&chirp)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to decode chirp submission", err)
		return
	}

	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		respondError(w, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	cleanedChirp := cleanChirp(chirp.Body)

	dbChirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedChirp,
		UserID: uuid.NullUUID{
			UUID:  userFromToken,
			Valid: true,
		},
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create chirp", err)
		return
	}
	respondJSON(w, http.StatusCreated, mapDatabaseChirpToChirp(dbChirp))
	return
}

func cleanChirp(msg string) string {
	badwords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(msg, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badwords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleanChirp := strings.Join(words, " ")
	return cleanChirp
}

func mapDatabaseChirpToChirp(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID.UUID,
	}
}
