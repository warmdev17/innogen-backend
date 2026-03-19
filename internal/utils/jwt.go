// Package utils
package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret returns the JWT secret for access tokens
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set")
	}
	return []byte(secret)
}

// getJWTRefreshSecret returns the JWT secret for refresh tokens
func getJWTRefreshSecret() []byte {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		panic("JWT_REFRESH_SECRET environment variable is not set")
	}
	return []byte(secret)
}

// GenerateAccessToken creates a short-lived access token (15 minutes)
func GenerateAccessToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    "access",
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// GenerateRefreshToken creates a JWT refresh token (30 days)
func GenerateRefreshToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    "refresh",
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTRefreshSecret())
}

// GenerateSecureToken creates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashToken hashes a token string using SHA256
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// ParseAccessToken parses and validates an access token
func ParseAccessToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}

		// Verify token type
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, jwt.ErrInvalidKey
		}

		if claims["type"] != "access" {
			return nil, jwt.ErrInvalidKey
		}

		return getJWTSecret(), nil
	})
}

// ParseRefreshToken parses and validates a refresh token
func ParseRefreshToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}

		// Verify token type
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, jwt.ErrInvalidKey
		}

		if claims["type"] != "refresh" {
			return nil, jwt.ErrInvalidKey
		}

		return getJWTRefreshSecret(), nil
	})
}

// ParseToken parses any JWT token (legacy support - defaults to access token)
func ParseToken(tokenStr string) (*jwt.Token, error) {
	return ParseAccessToken(tokenStr)
}

// GetTokenClaims extracts claims from a parsed token
func GetTokenClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}
	return claims, nil
}
