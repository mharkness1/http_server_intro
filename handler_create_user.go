package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mharkness1/http_server_intro/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func mapDatabaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	userEmail := createUserRequest{}
	err := decoder.Decode(&userEmail)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to decode user info", err)
		return
	}

	databaseUser, err := cfg.DB.CreateUser(r.Context(), userEmail.Email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}
	user := mapDatabaseUserToUser(databaseUser)
	respondJSON(w, http.StatusCreated, user)
}
