package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/warmdev17/innogen-backend/internal/database"
)

const OTPExpiration = 5 * time.Minute

// GenerateOTP generates a random 6-digit OTP.
func GenerateOTP() string {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "000000"
	}
	return fmt.Sprintf("%06d", n.Int64())
}

// StoreOTP stores the OTP in Redis.
func StoreOTP(email, otp string) error {
	ctx := context.Background()
	key := fmt.Sprintf("otp:%s", email)
	return database.Rdb.Set(ctx, key, otp, OTPExpiration).Err()
}

// VerifyOTP verifies the OTP from Redis.
func VerifyOTP(email, otp string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("otp:%s", email)

	val, err := database.Rdb.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	if val == otp {
		// Delete OTP after successful verification
		database.Rdb.Del(ctx, key)
		return true
	}

	return false
}
