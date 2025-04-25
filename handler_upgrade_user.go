package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mharkness1/http_server_intro/internal/auth"
)

func (cfg *apiConfig) upgradeUserHandler(w http.ResponseWriter, r *http.Request) {
	type PoklaRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKeyRequest, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "No api key found", err)
		return
	}
	if apiKeyRequest != cfg.POLKA_KEY {
		respondError(w, http.StatusUnauthorized, "Api key doesn't match", nil)
	}

	decoder := json.NewDecoder(r.Body)
	userUpgradeRequest := PoklaRequest{}
	err = decoder.Decode(&userUpgradeRequest)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error decoding webhook", err)
		return
	}

	if userUpgradeRequest.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userUuid, err := uuid.Parse(userUpgradeRequest.Data.UserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to extract user id", err)
		return
	}
	_, err = cfg.DB.UpgradeUserChirpyRed(r.Context(), userUuid)
	if err != nil {
		respondError(w, http.StatusNotFound, "failed to upgrade user's status", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
