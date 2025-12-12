package space

import (
	"fmt"
	"math"
	"time"
)

// SpaceActionRequest represents an attempt to fire a weapon or use a module
type SpaceActionRequest struct {
	ModuleID   string    `json:"module_id"`
	TargetID   string    `json:"target_id,omitempty"`
	TargetPos  Vector3   `json:"target_pos"` // Predicted position
	Timestamp  time.Time `json:"timestamp"`
}

// SpaceActionContext holds state needed for validation
type SpaceActionContext struct {
	ShipPos        Vector3
	ShipCapacitor  float64
	GlobalCooldown time.Time
	ModuleStats    ModuleStats
	TargetLocked   bool
}

// ModuleStats defines the parameters for a ship module
type ModuleStats struct {
	Range         float64
	CapacitorCost float64
	Cooldown      time.Duration
	RequiresLock  bool
}

// ValidateSpaceAction checks if a space action is legal
func ValidateSpaceAction(req SpaceActionRequest, ctx SpaceActionContext) (bool, error) {
	// 1. Global Cooldown / Module Cooldown Check
	tolerance := 200 * time.Millisecond
	if req.Timestamp.Add(tolerance).Before(ctx.GlobalCooldown) {
		return false, fmt.Errorf("module on cooldown")
	}

	// 2. Capacitor Check
	if ctx.ShipCapacitor < ctx.ModuleStats.CapacitorCost {
		return false, fmt.Errorf("insufficient capacitor: has %.1f, needs %.1f", ctx.ShipCapacitor, ctx.ModuleStats.CapacitorCost)
	}

	// 3. Lock Requirement
	if ctx.ModuleStats.RequiresLock && !ctx.TargetLocked {
		return false, fmt.Errorf("no target lock")
	}

	// 4. Range Check (3D)
	dx := req.TargetPos.X - ctx.ShipPos.X
	dy := req.TargetPos.Y - ctx.ShipPos.Y
	dz := req.TargetPos.Z - ctx.ShipPos.Z
	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

	rangeTolerance := 5.0 // Larger tolerance in space for lag/speed
	maxDist := ctx.ModuleStats.Range + rangeTolerance

	if dist > maxDist {
		return false, fmt.Errorf("target out of range: %.2f > %.2f", dist, maxDist)
	}

	return true, nil
}
