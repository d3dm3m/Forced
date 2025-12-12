package ground

import (
	"fmt"
	"math"
	"time"
)

// ActionRequest represents a player's attempt to use an ability
type ActionRequest struct {
	AbilityID   string    `json:"ability_id"`
	TargetID    string    `json:"target_id,omitempty"`
	TargetPos   Vector2   `json:"target_pos"`
	Timestamp   time.Time `json:"timestamp"`
}

// ActionContext holds state needed for validation
type ActionContext struct {
	PlayerPos     Vector2
	PlayerStamina int
	GlobalCooldown time.Time // When the last action finished
	AbilityStats  AbilityStats
}

// AbilityStats defines the validation parameters for an ability
type AbilityStats struct {
	Range       float64
	StaminaCost int
	Cooldown    time.Duration
	CastTime    time.Duration
}

// ValidateAction checks if an action is legal
func ValidateAction(req ActionRequest, ctx ActionContext) (bool, error) {
	// 1. Global Cooldown Check
	// If the timestamp of the request is BEFORE the GCD expires, fail.
	// Allow a small tolerance for network latency.
	tolerance := 200 * time.Millisecond
	if req.Timestamp.Add(tolerance).Before(ctx.GlobalCooldown) {
		return false, fmt.Errorf("action on cooldown: ready at %v, req at %v", ctx.GlobalCooldown, req.Timestamp)
	}

	// 2. Stamina Check
	if ctx.PlayerStamina < ctx.AbilityStats.StaminaCost {
		return false, fmt.Errorf("insufficient stamina: has %d, needs %d", ctx.PlayerStamina, ctx.AbilityStats.StaminaCost)
	}

	// 3. Range Check
	dx := req.TargetPos.X - ctx.PlayerPos.X
	dy := req.TargetPos.Y - ctx.PlayerPos.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	// Validate against Max Range + Tolerance (hitbox size)
	rangeTolerance := 1.0 // Unit size buffer
	maxDist := ctx.AbilityStats.Range + rangeTolerance

	if dist > maxDist {
		return false, fmt.Errorf("target out of range: %.2f > %.2f", dist, maxDist)
	}

	return true, nil
}
