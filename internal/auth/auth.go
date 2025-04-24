package auth

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	x, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Printf("error generating hash of password: %v", err)
		return "", err
	}
	return string(x), nil
}

func CheckPasswordHash(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	uniqueID := fmt.Sprintf("%x", randomBytes)

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
		ID:        uniqueID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %v", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authorizationRequest := headers.Get("Authorization")
	if authorizationRequest == "" {
		return "", fmt.Errorf("No authorization header provided")
	}
	if strings.HasPrefix(authorizationRequest, "Bearer ") {
		userToken := strings.Trim(authorizationRequest[len("Bearer "):], " ")
		if len(userToken) == 0 {
			return "", fmt.Errorf("No bearer token found")
		}
		return userToken, nil
	}
	return "", fmt.Errorf("Bearer token incorrectly formatted")
}
