package main

import (
	"net/http"
	"time"

	"github.com/mharkness1/http_server_intro/internal/auth"
)

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	type refreshTokenResponse struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	currentTime := time.Now()
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Missing refresh token", err)
		return
	}
	databaseRefreshToken, err := cfg.DB.CheckRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "No refresh token found", err)
		return
	}
	if currentTime.After(databaseRefreshToken.ExpiresAt) {
		respondError(w, http.StatusUnauthorized, "Refresh token has expired", nil)
		return
	}
	if databaseRefreshToken.RevokedAt.Valid {
		respondError(w, http.StatusUnauthorized, "Refresh token has been revoked", nil)
		return
	}

	accessToken, err := auth.MakeJWT(databaseRefreshToken.UserID, cfg.SECRET, time.Hour)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate access token", err)
		return
	}

	respondJSON(w, http.StatusOK, refreshTokenResponse{
		Token: accessToken,
	})
	return
}

func (cfg *apiConfig) revokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Missing refresh token", err)
		return
	}
	err = cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to revoke token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}
