package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}
	extractedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}
	if extractedID != userID {
		t.Errorf("Extracted user ID doesn't match original: got %v, want %v", extractedID, userID)
	}
}

func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	// Create a token that expires immediately
	token, err := MakeJWT(userID, secret, -1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create expired JWT: %v", err)
	}

	// Validation should fail with an error about expiration
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestInvalidSecret(t *testing.T) {
	userID := uuid.New()
	secret := "correct-secret"
	wrongSecret := "wrong-secret"

	// Create a token with the correct secret
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Try to validate with wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("Expected error when validating with wrong secret, got nil")
	}
}

func TestMalformedToken(t *testing.T) {
	// Test with completely invalid token
	_, err := ValidateJWT("not-a-valid-token", "any-secret")
	if err == nil {
		t.Error("Expected error for malformed token, got nil")
	}

	// Test with tampered token
	userID := uuid.New()
	secret := "test-secret"
	token, _ := MakeJWT(userID, secret, time.Hour)
	tamperedToken := token + "tampered"

	_, err = ValidateJWT(tamperedToken, secret)
	if err == nil {
		t.Error("Expected error for tampered token, got nil")
	}
}

func TestTokenWithCorrectFormat(t *testing.T) {
	// Generate a UUID for testing
	userID := uuid.New()
	secret := "test-secret"

	// Create token
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Verify token format (should be 3 parts separated by periods)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Expected JWT with 3 parts, got %d parts", len(parts))
	}
}
func TestMultipleTokensForSameUser(t *testing.T) {
	// Test that multiple tokens for the same user are unique
	userID := uuid.New()
	secret := "test-secret"

	token1, _ := MakeJWT(userID, secret, time.Hour)
	time.Sleep(10 * time.Millisecond)
	token2, _ := MakeJWT(userID, secret, time.Hour)

	t.Logf("Token1: %s", token1)
	t.Logf("Token2: %s", token2)

	if token1 == token2 {
		t.Error("Expected different tokens for same user created at different times")
	}

	// But they should both validate to the same user ID
	id1, err1 := ValidateJWT(token1, secret)
	id2, err2 := ValidateJWT(token2, secret)

	if err1 != nil || err2 != nil {
		t.Fatalf("Validation errors: %v, %v", err1, err2)
	}

	if id1 != userID || id2 != userID {
		t.Error("Extracted IDs don't match original user ID")
	}
}
