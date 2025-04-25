package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mharkness1/http_server_intro/internal/auth"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "token extraction failed", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.SECRET)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "token failed validation", err)
		return
	}
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
	if chirp.UserID.UUID != userID {
		respondError(w, http.StatusForbidden, "Chirp creator and User ID don't match", nil)
		return
	}
	err = cfg.DB.DeleteChirpByID(r.Context(), chirpUuid)
	if err != nil {
		respondError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}
