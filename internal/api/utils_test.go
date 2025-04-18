package api

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

// CleanBadWords

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

func TestGetAuthorizationField(t *testing.T) {

	expected := "Bearer SomeValue"
	newRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	newRequest.Header.Set("Authorization", expected)
	value, err := GetAuthorizationField(newRequest.Header)
	if value != "SomeValue" {
		t.Fatalf("Extracted value %sis not expected %s", value, expected)
	}

	if err != nil {
		t.Fatal("Error extracting Field from Header")
	}
}

func TestGetAuthorizationFieldFromInvalidHeader(t *testing.T) {

	newRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	newRequest.Header.Set("Authorization", "Bearer")
	_, err := GetAuthorizationField(newRequest.Header)

	if err != nil {
		t.Log("Invalid header as expected")
	}
}

func TestCleanBadWords(t *testing.T) {

	testCases := []struct {
		input    string
		expected string
	}{
		{"This is a kerfuffle", "This is a ****"},
		{"This is a kerfuffle and a sharbert", "This is a **** and a ****"},
		{"This is a kerfuffle and a sharbert and a fornax", "This is a **** and a **** and a ****"},
		{"This is a clean sentence", "This is a clean sentence"},
	}

	for _, sentence := range testCases {
		cleaned := CleanBadWords(sentence.input)
		if cleaned != sentence.expected {
			t.Fatalf("Expected cleaned sentence to be '%s', got '%s'", sentence.expected, cleaned)
		}
	}

}
