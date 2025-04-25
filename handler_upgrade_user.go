package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) upgradeUserHandler(w http.ResponseWriter, r *http.Request) {
	type PoklaRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	userUpgradeRequest := PoklaRequest{}
	err := decoder.Decode(&userUpgradeRequest)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error decoding webhook", err)
		return
	}

	if userUpgradeRequest.Event != "user.upgraded" {
		respondError(w, http.StatusNoContent, "incorrect event type", nil)
	}

	userUuid, err := uuid.Parse(userUpgradeRequest.Data.UserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to extract user id", err)
	}
	_, err = cfg.DB.UpgradeUserChirpyRed(r.Context(), userUuid)
	if err != nil {
		respondError(w, http.StatusNotFound, "failed to upgrade user's status", err)
	}
	w.WriteHeader(204)
}
