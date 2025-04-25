package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mharkness1/http_server_intro/internal/auth"
	"github.com/mharkness1/http_server_intro/internal/database"
)

type User struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Email           string    `json:"email"`
	Hashed_Password string    `json:"hashed_password"`
	IsChirpyRed     bool      `json:"is_chirpy_red"`
}

func mapDatabaseUserToUserResponse(dbUser database.User) User {
	return User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	userRequest := createUserRequest{}
	err := decoder.Decode(&userRequest)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to decode user info", err)
		return
	}

	hashedPassword, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to encode password", err)
		return
	}

	databaseUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          userRequest.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}
	user := mapDatabaseUserToUserResponse(databaseUser)
	respondJSON(w, http.StatusCreated, user)
	return
}

func (cfg *apiConfig) updateUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	type updateUserInfoRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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

	decoder := json.NewDecoder(r.Body)
	userRequest := updateUserInfoRequest{}
	err = decoder.Decode(&userRequest)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to decode user info", err)
		return
	}

	hashedPassword, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed hashing password", err)
		return
	}

	updatedUser, err := cfg.DB.UpdateUserInfo(r.Context(), database.UpdateUserInfoParams{
		ID:             userID,
		Email:          userRequest.Email,
		HashedPassword: hashedPassword,
	})

	respondJSON(w, http.StatusOK, mapDatabaseUserToUserResponse(updatedUser))
	return
}
