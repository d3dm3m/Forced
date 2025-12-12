package space

import (
	"testing"
	"time"
)

func TestValidateSpaceAction(t *testing.T) {
	now := time.Now()

	stats := ModuleStats{
		Range:         100.0,
		CapacitorCost: 50.0,
		Cooldown:      2 * time.Second,
		RequiresLock:  true,
	}

	tests := []struct {
		name      string
		req       SpaceActionRequest
		ctx       SpaceActionContext
		wantValid bool
	}{
		{
			name: "Valid Shot",
			req: SpaceActionRequest{
				TargetPos: Vector3{X: 50, Y: 0, Z: 0},
				Timestamp: now,
			},
			ctx: SpaceActionContext{
				ShipPos:        Vector3{X: 0, Y: 0, Z: 0},
				ShipCapacitor:  100.0,
				GlobalCooldown: now.Add(-1 * time.Second),
				ModuleStats:    stats,
				TargetLocked:   true,
			},
			wantValid: true,
		},
		{
			name: "Insufficient Capacitor",
			req: SpaceActionRequest{
				TargetPos: Vector3{X: 50, Y: 0, Z: 0},
				Timestamp: now,
			},
			ctx: SpaceActionContext{
				ShipPos:        Vector3{X: 0, Y: 0, Z: 0},
				ShipCapacitor:  10.0, // Needs 50
				GlobalCooldown: now.Add(-1 * time.Second),
				ModuleStats:    stats,
				TargetLocked:   true,
			},
			wantValid: false,
		},
		{
			name: "No Lock",
			req: SpaceActionRequest{
				TargetPos: Vector3{X: 50, Y: 0, Z: 0},
				Timestamp: now,
			},
			ctx: SpaceActionContext{
				ShipPos:        Vector3{X: 0, Y: 0, Z: 0},
				ShipCapacitor:  100.0,
				GlobalCooldown: now.Add(-1 * time.Second),
				ModuleStats:    stats,
				TargetLocked:   false, // Logic failure
			},
			wantValid: false,
		},
		{
			name: "Out of Range",
			req: SpaceActionRequest{
				TargetPos: Vector3{X: 110, Y: 0, Z: 0}, // Max 100 + 5
				Timestamp: now,
			},
			ctx: SpaceActionContext{
				ShipPos:        Vector3{X: 0, Y: 0, Z: 0},
				ShipCapacitor:  100.0,
				GlobalCooldown: now.Add(-1 * time.Second),
				ModuleStats:    stats,
				TargetLocked:   true,
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateSpaceAction(tt.req, tt.ctx)
			if got != tt.wantValid {
				t.Errorf("ValidateSpaceAction() = %v, want %v (err: %v)", got, tt.wantValid, err)
			}
		})
	}
}
