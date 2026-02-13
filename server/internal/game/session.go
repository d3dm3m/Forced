package game

import (
	"context"
	"fmt"
	"time"
	"biohorror/internal/db"
	"github.com/google/uuid"
)

const TokenTTL = 60 * time.Second

// GenerateTransferToken creates a one-time use token for session handoff
func GenerateTransferToken(username string) (string, error) {
	if db.RedisClient == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	token := uuid.New().String()
	key := fmt.Sprintf("token:%s", token)

	ctx := context.Background()

	// Store username as the value, with TTL
	err := db.RedisClient.Set(ctx, key, username, TokenTTL).Err()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// ValidateTransferToken checks if a token is valid and returns the associated username.
// It deletes the token immediately (One-Time Use).
func ValidateTransferToken(token string) (string, error) {
	if db.RedisClient == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("token:%s", token)
	ctx := context.Background()

	// Get value
	username, err := db.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("invalid or expired token")
	}

	// Delete token (One-Time Use)
	db.RedisClient.Del(ctx, key)

	return username, nil
}
