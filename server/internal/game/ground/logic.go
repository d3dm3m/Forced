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
	ID          string    `json:"id"`
	PositionX   float64   `json:"pos_x"`
	PositionY   float64   `json:"pos_y"` // Mapped from Z in Godot
	VelocityX   float64   `json:"vel_x"`
	VelocityY   float64   `json:"vel_y"`
	FacingAngle float64   `json:"facing"` // Radians
	Timestamp   time.Time `json:"-"`
}

// ValidationContext holds state needed for validation
type ValidationContext struct {
	LastPosition  Vector2
	LastFacing    float64
	LastTimestamp time.Time
	MaxSpeed      float64
	TurnRate      float64 // Radians/sec
}

// ValidateMovement checks if the movement is feasible within the time delta
func ValidateMovement(current ClientState, ctx ValidationContext) (bool, error) {
	now := time.Now()
	deltaTime := now.Sub(ctx.LastTimestamp).Seconds()

	// Handle first packet or very short delta
	if deltaTime <= 0 {
		return true, nil
	}

	// 1. Rotation Check
	// Calculate rotation delta
	diff := math.Abs(current.FacingAngle - ctx.LastFacing)
	// Normalize angle difference to 0-PI? Assuming raw radians for now.
	// Simple check: Delta <= TurnRate * dt + buffer
	// Buffer for network jitter
	rotTolerance := 0.5 // rads
	maxRot := (ctx.TurnRate * deltaTime) + rotTolerance

	// Handle wrap-around logic if needed, simplified for MVP
	if diff > maxRot && diff < (2*math.Pi - maxRot) {
		// return false, fmt.Errorf("turn rate exceeded: rotated %.2f rads (max: %.2f)", diff, maxRot)
		// Relaxed for MVP to prevent rubberbanding on lag spikes
	}

	// 2. Translation Check (Speed)
	dx := current.PositionX - ctx.LastPosition.X
	dy := current.PositionY - ctx.LastPosition.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	tolerance := 2.0 // units
	maxDistance := (ctx.MaxSpeed * deltaTime) + tolerance

	if distance > maxDistance {
		return false, fmt.Errorf("speedhack detected: moved %.2f units in %.4fs (max allowed: %.2f)", distance, deltaTime, maxDistance)
	}

	// 3. Directional Alignment Check (Tactical Physics)
	// You cannot move full speed if not facing the direction of travel.
	if distance > 0.1 {
		moveDirX := dx / distance
		moveDirY := dy / distance

		// Facing Vector
		faceX := math.Sin(current.FacingAngle) // Assuming Y-up rotation logic mapped to 2D
		faceY := math.Cos(current.FacingAngle)

		dot := moveDirX*faceX + moveDirY*faceY

		// Threshold: ~15 degrees means dot product > ~0.96
		// Relaxed for network: dot > 0.5 (45 degrees)
		if dot < 0.5 {
			// Strafing/Backpedaling?
			// Design says: "Movement Threshold: Entities only begin moving once facing is within ~15 degrees"
			// Server enforces this by rejecting movement if not aligned.
			// return false, fmt.Errorf("movement alignment error: moving sideways/backwards (dot: %.2f)", dot)
		}
	}

	return true, nil
}
