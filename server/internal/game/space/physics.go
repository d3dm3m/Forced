package space

import (
	"fmt"
	"math"
)

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

// ValidateVector checks if the client's position is within a reasonable margin of the expected dead-reckoning position.
// clientPos: The position reported by the client.
// clientVel: The velocity reported by the client (or calculated from previous ticks).
// lastServerPos: The last valid position known by the server.
// lastServerVel: The last valid velocity known by the server (used for projection).
// dt: Delta time in seconds since the last update.
func ValidateVector(clientPos Vector3, clientVel Vector3, lastServerPos Vector3, lastServerVel Vector3, dt float64) (bool, error) {
	// 1. Calculate Expected Position based on Dead Reckoning (Inertia)
	// NewPos = OldPos + (OldVel * dt)
	expectedX := lastServerPos.X + (lastServerVel.X * dt)
	expectedY := lastServerPos.Y + (lastServerVel.Y * dt)
	expectedZ := lastServerPos.Z + (lastServerVel.Z * dt)

	// 2. Calculate Deviation (Distance between Expected and Reported)
	dx := clientPos.X - expectedX
	dy := clientPos.Y - expectedY
	dz := clientPos.Z - expectedZ
	distance := math.Sqrt(dx*dx + dy*dy + dz*dz)

	// 3. Define Tolerance
	// Network jitter, float imprecision, and client-side input acceleration (which happens *during* dt)
	// mean we can't expect a perfect match.
	// Allow for specific "max acceleration" deviation.
	// Deviation <= 0.5 * MaxAccel * dt^2 (Physics formula for distance traveled under accel)
	// Plus a buffer for latency jitter.

	const MaxAcceleration = 50.0 // Matches client-side acceleration_force
	const LatencyBuffer = 2.0 // Units of margin

	allowedDeviation := (0.5 * MaxAcceleration * dt * dt) + LatencyBuffer

	if distance > allowedDeviation {
		return false, fmt.Errorf("position deviation too high: %.2f > %.2f (dt: %.4f)", distance, allowedDeviation, dt)
	}

	return true, nil
}
