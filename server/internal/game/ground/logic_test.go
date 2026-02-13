package ground

import (
	"testing"
	"time"
)

func TestValidateMovement(t *testing.T) {
	ctx := ValidationContext{
		LastPosition:  Vector2{X: 0, Y: 0},
		LastTimestamp: time.Now().Add(-1 * time.Second), // 1 second ago
		MaxSpeed:      10.0,
	}

	// Valid Move (Moved 10 units in 1 second)
	validState := ClientState{
		PositionX: 10,
		PositionY: 0,
	}

	ok, err := ValidateMovement(validState, ctx)
	if !ok {
		t.Errorf("Expected valid movement, got error: %v", err)
	}

	// Invalid Move (Moved 20 units in 1 second, Max is 10 + 2 buffer = 12)
	invalidState := ClientState{
		PositionX: 20,
		PositionY: 0,
	}

	ok, err = ValidateMovement(invalidState, ctx)
	if ok {
		t.Error("Expected invalid movement (speedhack), got ok")
	} else {
		t.Logf("Correctly detected invalid movement: %v", err)
	}
}
