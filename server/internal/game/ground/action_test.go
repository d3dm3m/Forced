package ground

import (
	"biohorror/internal/game/data"
	"testing"
	"time"
)

func TestValidateAction(t *testing.T) {
	now := time.Now()

	def := data.AbilityDefinition{
		Range:       10.0,
		StaminaCost: 20,
		Cooldown:    1 * time.Second,
		Type:        data.AbilityTypeAttack,
	}

	tests := []struct {
		name      string
		req       ActionRequest
		ctx       ActionContext
		wantValid bool
	}{
		{
			name: "Valid Action",
			req: ActionRequest{
				TargetPos: Vector2{X: 5, Y: 0},
				Timestamp: now,
			},
			ctx: ActionContext{
				PlayerPos:      Vector2{X: 0, Y: 0},
				PlayerStamina:  100,
				GlobalCooldown: now.Add(-1 * time.Second),
				AbilityDef:     def,
			},
			wantValid: true,
		},
		{
			name: "Insufficient Stamina",
			req: ActionRequest{
				TargetPos: Vector2{X: 5, Y: 0},
				Timestamp: now,
			},
			ctx: ActionContext{
				PlayerPos:      Vector2{X: 0, Y: 0},
				PlayerStamina:  10, // Needs 20
				GlobalCooldown: now.Add(-1 * time.Second),
				AbilityDef:     def,
			},
			wantValid: false,
		},
		{
			name: "Out of Range",
			req: ActionRequest{
				TargetPos: Vector2{X: 12, Y: 0}, // Max 10 + 1 tolerance
				Timestamp: now,
			},
			ctx: ActionContext{
				PlayerPos:      Vector2{X: 0, Y: 0},
				PlayerStamina:  100,
				GlobalCooldown: now.Add(-1 * time.Second),
				AbilityDef:     def,
			},
			wantValid: false,
		},
		{
			name: "On Cooldown",
			req: ActionRequest{
				TargetPos: Vector2{X: 5, Y: 0},
				Timestamp: now,
			},
			ctx: ActionContext{
				PlayerPos:      Vector2{X: 0, Y: 0},
				PlayerStamina:  100,
				GlobalCooldown: now.Add(500 * time.Millisecond), // In future
				AbilityDef:     def,
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAction(tt.req, tt.ctx)
			if got != tt.wantValid {
				t.Errorf("ValidateAction() = %v, want %v (err: %v)", got, tt.wantValid, err)
			}
		})
	}
}
