package main

import (
	"encoding/json"
	"log"
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

	databaseUser, err := cfg.DB.CreateUser(r.Context(), userEmail.Email)
	if err != nil {
		w.WriteHeader(500)
	}
	user := mapDatabaseUserToUser(databaseUser)

	jsonBody, err := json.Marshal(user)
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(jsonBody)
}
