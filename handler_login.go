package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mharkness1/http_server_intro/internal/auth"
	"github.com/mharkness1/http_server_intro/internal/database"
)

type UserLogin struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Email           string    `json:"email"`
	Hashed_Password string    `json:"hashed_password"`
	Token           string    `json:"token"`
	RefreshToken    string    `json:"refresh_token"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	login := loginRequest{}
	err := decoder.Decode(&login)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to decode login info", err)
	}

	user, err := cfg.DB.Login(r.Context(), login.Email)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Incorrect email or password", err)
	}

	err = auth.CheckPasswordHash(user.HashedPassword, login.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Incorrect email or password", err)
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.SECRET, time.Hour)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate JWT", err)
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error generating refresh token", err)
	}
	_, err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to store refresh token", err)
	}

	respondJSON(w, http.StatusOK, mapToLoginResponse(user, jwt, refreshToken))
}

func mapToLoginResponse(dbUser database.User, jwt string, refreshToken string) UserLogin {
	return UserLogin{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        jwt,
		RefreshToken: refreshToken,
	}
}
