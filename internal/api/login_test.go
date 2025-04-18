package api

import (
	"log"
	"testing"

	"github.com/google/uuid"
)

func TestJWTCreationAndValidation(t *testing.T) {
	// Setup
	userID := uuid.New()
	secret := "my-secret-key"

	// Test creating a JWT
	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}
	if token == "" {
		t.Fatal("Expected token to be non-empty")
	}

	// Test validating a valid JWT
	extractedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate valid JWT: %v", err)
	}
	if extractedID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, extractedID)
	}

	log.Print(userID == extractedID)
}
