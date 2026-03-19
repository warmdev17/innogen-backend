package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndParseAccessToken(t *testing.T) {
	// Setup env vars for testing
	os.Setenv("JWT_SECRET", "test_secret_key")
	defer os.Unsetenv("JWT_SECRET")

	userID := uint(1)
	role := "student"

	// Test Generation
	tokenString, err := GenerateAccessToken(userID, role)
	if err != nil {
		t.Fatalf("Failed to generate access token: %v", err)
	}
	if tokenString == "" {
		t.Error("GenerateAccessToken returned empty string")
	}

	// Test Parsing
	parsedToken, err := ParseAccessToken(tokenString)
	if err != nil {
		t.Fatalf("Failed to parse access token: %v", err)
	}

	claims, err := GetTokenClaims(parsedToken)
	if err != nil {
		t.Fatalf("Failed to extract claims: %v", err)
	}

	// Validate claims
	if uint(claims["user_id"].(float64)) != userID {
		t.Errorf("Expected user_id %d, got %v", userID, claims["user_id"])
	}

	if claims["role"] != role {
		t.Errorf("Expected role %s, got %v", role, claims["role"])
	}

	if claims["type"] != "access" {
		t.Errorf("Expected type 'access', got %v", claims["type"])
	}

	// Validate expiration time is in the future
	exp := int64(claims["exp"].(float64))
	if exp <= time.Now().Unix() {
		t.Error("Token expiration time should be in the future")
	}
}

func TestGenerateAndParseRefreshToken(t *testing.T) {
	os.Setenv("JWT_REFRESH_SECRET", "test_refresh_secret")
	defer os.Unsetenv("JWT_REFRESH_SECRET")

	userID := uint(2)
	role := "admin"

	// Test Generation
	tokenString, err := GenerateRefreshToken(userID, role)
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	// Test Parsing
	parsedToken, err := ParseRefreshToken(tokenString)
	if err != nil {
		t.Fatalf("Failed to parse refresh token: %v", err)
	}

	claims, err := GetTokenClaims(parsedToken)
	if err != nil {
		t.Fatalf("Failed to get token claims: %v", err)
	}

	if claims["type"] != "refresh" {
		t.Errorf("Expected type 'refresh', got %v", claims["type"])
	}
}

func TestParseAccessTokenWithInvalidSignature(t *testing.T) {
	os.Setenv("JWT_SECRET", "correct_secret")
	defer os.Unsetenv("JWT_SECRET")

	// Create a token signed with wrong secret
	claims := jwt.MapClaims{"type": "access"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("wrong_secret"))

	_, err := ParseAccessToken(tokenString)
	if err == nil {
		t.Error("Parsing should fail when signature is invalid")
	}
}
