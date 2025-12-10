package game

import (
	"math"
	"math/rand"
	"time"
	"biohorror/pkg/protocol"
	"github.com/google/uuid"
)

// Projectile represents a simulated shot in flight.
type Projectile struct {
	ID         string
	ShooterID  string
	TargetID   string
	Behavior   string // "missile", "linear"
	Outcome    string // "Hit", "Miss"

	// Position 3D
	X, Y, Z    float64

	// Velocity 3D (For Linear/Missile physics)
	VX, VY, VZ float64

	// Target Position 3D (For Homing/Interpolation check)
	TargetX, TargetY, TargetZ float64

	Speed      float64
	TurnRate   float64 // Radians/sec for missiles

	Damage     float64
}

// ProjectileManager handles the lifecycle of projectiles.
type ProjectileManager struct {
	ActiveProjectiles map[string]*Projectile
	rng               *rand.Rand
}

func NewProjectileManager() *ProjectileManager {
	return &ProjectileManager{
		ActiveProjectiles: make(map[string]*Projectile),
		rng:               rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CalculateHitChance determines if a shot lands based on stats.
func (pm *ProjectileManager) CalculateHitChance(attacker *Player, target *Player) bool {
	// Base Chance
	chance := 0.8

	// Stats Modifiers (Sprint 18 requirement)
	// Ground Gear check
	if attacker.GroundGear.Talents != nil {
		if marks, ok := attacker.GroundGear.Talents["marksmanship"]; ok {
			chance += float64(marks) * 0.05
		}
	}

	// TODO: Add Ship Stats check for Space Combat

	return pm.rng.Float64() < chance
}

// SpawnProjectile creates a new projectile and adds it to the simulation.
func (pm *ProjectileManager) SpawnProjectile(
	shooterID, targetID string,
	sx, sy, sz, ex, ey, ez float64,
	speed, damage, turnRate float64,
	behavior, outcome string,
) *Projectile {

	id := uuid.New().String()

	// Initial Velocity Direction
	dx := ex - sx
	dy := ey - sy
	dz := ez - sz
	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

	vx, vy, vz := 0.0, 0.0, 0.0
	if dist > 0 {
		vx = (dx / dist) * speed
		vy = (dy / dist) * speed
		vz = (dz / dist) * speed
	}

	proj := &Projectile{
		ID:         id,
		ShooterID:  shooterID,
		TargetID:   targetID,
		Behavior:   behavior,
		Outcome:    outcome,
		X: sx, Y: sy, Z: sz,
		TargetX: ex, TargetY: ey, TargetZ: ez,
		VX: vx, VY: vy, VZ: vz,
		Speed:      speed,
		TurnRate:   turnRate,
		Damage:     damage,
	}

	pm.ActiveProjectiles[id] = proj
	return proj
}

// UpdateSimulation ticks the projectiles. Returns list of impacts.
func (pm *ProjectileManager) UpdateSimulation(dt float64) []*Projectile {
	var impacts []*Projectile

	for id, p := range pm.ActiveProjectiles {
		// Distance to Target (Current)
		dx := p.TargetX - p.X
		dy := p.TargetY - p.Y
		dz := p.TargetZ - p.Z
		dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

		moveDist := p.Speed * dt

		if moveDist >= dist {
			// Impact!
			p.X = p.TargetX
			p.Y = p.TargetY
			p.Z = p.TargetZ
			impacts = append(impacts, p)
			delete(pm.ActiveProjectiles, id)
			continue
		}

		if p.Behavior == protocol.BEHAVIOR_MISSILE && dist > 0 {
			// Homing Logic: Rotate Velocity towards Target
			// Simplified: Just re-calculate normalized vector to target and blend?
			// Full TurnRate physics is complex. For MVP, we'll just steer directly
			// if within TurnRate, or assume perfect homing for now (Behavior = Missile just means it hits).
			// If Outcome == "Miss", we might target a false point?
			// For now, standard homing:

			// Update velocity vector to point at target
			p.VX = (dx / dist) * p.Speed
			p.VY = (dy / dist) * p.Speed
			p.VZ = (dz / dist) * p.Speed
		}

		// Move
		p.X += p.VX * dt
		p.Y += p.VY * dt
		p.Z += p.VZ * dt
	}

	return impacts
}
