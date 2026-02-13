package space

import (
	"testing"
)

func TestValidateVector(t *testing.T) {
	// Scenario: Ship drifting at constant velocity (10, 0, 0)
	lastPos := Vector3{X: 0, Y: 0, Z: 0}
	lastVel := Vector3{X: 10, Y: 0, Z: 0}
	dt := 1.0

	// Expected: (10, 0, 0)

	// Case 1: Perfect Match
	clientPos := Vector3{X: 10, Y: 0, Z: 0}
	ok, err := ValidateVector(clientPos, lastVel, lastPos, lastVel, dt)
	if !ok {
		t.Errorf("Expected valid vector, got error: %v", err)
	}

	// Case 2: Acceptable Deviation (Accelerated slightly)
	// Max Accel 50 -> 0.5 * 50 * 1^2 = 25 units allowed deviation + 2 buffer = 27
	clientPosAccel := Vector3{X: 30, Y: 0, Z: 0} // 20 units off
	ok, err = ValidateVector(clientPosAccel, lastVel, lastPos, lastVel, dt)
	if !ok {
		t.Errorf("Expected valid vector (within accel limits), got error: %v", err)
	}

	// Case 3: Teleportation (Too far)
	clientPosCheat := Vector3{X: 100, Y: 0, Z: 0} // 90 units off
	ok, err = ValidateVector(clientPosCheat, lastVel, lastPos, lastVel, dt)
	if ok {
		t.Error("Expected cheat detection, got ok")
	} else {
		t.Logf("Detected cheat: %v", err)
	}
}
