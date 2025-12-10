package ground

import (
	"fmt"
	"math"
	"time"
)

// Vector2 for simple logic
type Vector2 struct {
	X float64
	Y float64
}

// ClientState represents the state received from the client
type ClientState struct {
	ID        string    `json:"id"`
	PositionX float64   `json:"pos_x"`
	PositionY float64   `json:"pos_y"` // Mapped from Z in Godot
	VelocityX float64   `json:"vel_x"`
	VelocityY float64   `json:"vel_y"`
	Timestamp time.Time `json:"-"`
}

// ValidationContext holds state needed for validation
type ValidationContext struct {
	LastPosition  Vector2
	LastTimestamp time.Time
	MaxSpeed      float64
}

// ValidateMovement checks if the movement is feasible within the time delta
func ValidateMovement(current ClientState, ctx ValidationContext) (bool, error) {
	now := time.Now()
	deltaTime := now.Sub(ctx.LastTimestamp).Seconds()

	// Handle first packet or very short delta
	if deltaTime <= 0 {
		return true, nil
	}

	// Calculate Distance Traveled
	dx := current.PositionX - ctx.LastPosition.X
	dy := current.PositionY - ctx.LastPosition.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// Calculate Max Allowed Distance (Speed * Time)
	// Add a tolerance buffer (e.g., 10% or fixed units) for network jitter / lag compensation
	tolerance := 2.0 // units
	maxDistance := (ctx.MaxSpeed * deltaTime) + tolerance

	if distance > maxDistance {
		return false, fmt.Errorf("speedhack detected: moved %.2f units in %.4fs (max allowed: %.2f)", distance, deltaTime, maxDistance)
	}

	return true, nil
}
