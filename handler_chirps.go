package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
		Body    string `json:"body"`
		User_id string `json:"user_id"`
	}

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
