package game

import (
	"fmt"
	"github.com/google/uuid"
)

// RedisGatekeeper handles cross-server travel and system routing.
// For Sprint 18 MVP, this is a stubbed implementation backed by hardcoded values.
type RedisGatekeeper struct {
	// In future, this would hold a Redis client to query the 'gatekeeper:routes' hash.
}

func NewRedisGatekeeper() *RedisGatekeeper {
	return &RedisGatekeeper{}
}

// GetServerForSystem returns the websocket/udp URL for the server hosting the given system.
// Currently hardcoded to localhost for the Twin-Core MVP.
func (g *RedisGatekeeper) GetServerForSystem(systemID string) (string, error) {
	// Mock Routing Table
	routes := map[string]string{
		"Sol-0":      "ws://localhost:8080/ws",
		"Eden's Rot": "ws://localhost:8081/ws", // Hypothetical Shard 2
	}

	if url, ok := routes[systemID]; ok {
		return url, nil
	}

	// Fallback for MVP
	return "ws://localhost:8080/ws", nil
}

// RegisterUserTransfer initiates a handoff sequence.
// It generates a one-time token that the destination server will validate.
func (g *RedisGatekeeper) RegisterUserTransfer(userID string, targetServerID string) (string, error) {
	// 1. Validate User State (e.g. not in combat)
	// (Skipped for MVP)

	// 2. Generate Token
	// In a real implementation, this would use `game.GenerateTransferToken` logic
	// but store it in a globally accessible Redis key like `transfer:{token}`
	// that the Target Server can read.

	// For this stub, we simulate success.
	token := uuid.New().String()

	// Log the intent (Real impl would write to Redis)
	fmt.Printf("Gatekeeper: Registered transfer for %s to %s. Token: %s\n", userID, targetServerID, token)

	return token, nil
}
