package ground

import (
	"biohorror/internal/game/data"
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
	PlayerPos      Vector2
	PlayerStamina  int
	GlobalCooldown time.Time // When the last action finished
	AbilityDef     data.AbilityDefinition
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
	if ctx.PlayerStamina < ctx.AbilityDef.StaminaCost {
		return false, fmt.Errorf("insufficient stamina: has %d, needs %d", ctx.PlayerStamina, ctx.AbilityDef.StaminaCost)
	}

	// 3. Range Check
	// If ability target is self (e.g. Buff), skip range check or enforce 0 distance?
	// For MVP, treating Buffs as targeted on self location.
	if ctx.AbilityDef.Type == data.AbilityTypeBuff {
		// Range check is N/A or check if TargetPos is PlayerPos?
		// Skipping for now.
		return true, nil
	}

	dx := req.TargetPos.X - ctx.PlayerPos.X
	dy := req.TargetPos.Y - ctx.PlayerPos.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	// Validate against Max Range + Tolerance (hitbox size)
	// Currently Range is not on AbilityDefinition? It should be?
	// Assuming 0 range for now if not defined, will fix in registry.
	rangeLimit := ctx.AbilityDef.Range
	rangeTolerance := 1.0
	maxDist := rangeLimit + rangeTolerance

	if dist > maxDist {
		return false, fmt.Errorf("target out of range: %.2f > %.2f", dist, maxDist)
	}

	return true, nil
}
