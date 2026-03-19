package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
)

// generateSecureToken creates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// hashToken hashes a token string using SHA256
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// CreateRefreshToken creates a new refresh token for a user
func CreateRefreshToken(userID uint) (string, error) {
	// Generate random token string
	rawToken, err := generateSecureToken(64)
	if err != nil {
		return "", err
	}

	// Hash the token
	hashedToken := hashToken(rawToken)

	// Create token record
	refreshToken := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	// Save to database
	if err := database.DB.Create(&refreshToken).Error; err != nil {
		return "", err
	}

	return rawToken, nil
}

// VerifyRefreshToken checks if a refresh token is valid
func VerifyRefreshToken(rawToken string, userID uint) (*models.RefreshToken, error) {
	hashedToken := hashToken(rawToken)

	var refreshToken models.RefreshToken
	if err := database.DB.Where("user_id = ? AND token_hash = ? AND revoked = ? AND expires_at > ?",
		userID, hashedToken, false, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// RevokeRefreshToken marks a refresh token as revoked
func RevokeRefreshToken(rawToken string, userID uint) error {
	refreshToken, err := VerifyRefreshToken(rawToken, userID)
	if err != nil {
		return err
	}

	refreshToken.Revoked = true
	return database.DB.Save(&refreshToken).Error
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func RevokeAllUserTokens(userID uint) error {
	return database.DB.Model(&models.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}

// RotateRefreshToken creates a new refresh token and revokes the old one
func RotateRefreshToken(oldToken string, userID uint) (string, error) {
	// Verify the old token exists and is valid
	_, err := VerifyRefreshToken(oldToken, userID)
	if err != nil {
		return "", err
	}

	// Revoke all tokens for the user (more secure - prevents token reuse)
	if err := RevokeAllUserTokens(userID); err != nil {
		return "", err
	}

	// Create a new token
	return CreateRefreshToken(userID)
}