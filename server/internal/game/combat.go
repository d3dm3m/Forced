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

// ResolveSpaceTurretFire calculates hit chance and damage for space combat.
func (pm *ProjectileManager) ResolveSpaceTurretFire(attacker *Player, target *Player, distance float64) (bool, float64) {
	// 1. Get Attacker Stats
	attackerStats := CalculateShipStats(attacker)
	// Proxy: Use SensorRange as Slew Rate proxy (high sensor = fast lock/track)
	// Normalized: 100 sensor range -> 1.0 rad/s tracking
	trackingSpeed := attackerStats.SensorRange / 100.0
	if trackingSpeed <= 0 {
		trackingSpeed = 0.5
	}

	// 2. Get Target Stats
	// targetStats := CalculateShipStats(target) // Unused until we map SigRadius
	// Proxy: We need Signature Radius. We'll use a derived value or default.
	// Since CalculateShipStats doesn't explicitly return SigRadius in DerivedStats (yet),
	// we will default it or pull from base ships if we had the ID.
	// For now, let's assume standard Frigate size (50.0).
	targetSig := 50.0

	// 3. Simulate Transversal
	// Random 0-100 m/s
	transversal := pm.rng.Float64() * 100.0

	// 4. Calculate Chance
	chance := CalculateTurretTracking(trackingSpeed, targetSig, distance, transversal)

	// 5. Roll
	hit := pm.rng.Float64() < chance
	damage := 0.0

	if hit {
		// Base damage (e.g. 50) modified by nothing for now (Linear hit is binary usually, or graze)
		// We'll return full damage for Hit.
		// Check Attacker gun damage? We assume standard Laser Cannon (50.0) from Items.
		damage = 50.0

		// Apply modifiers from attacker stats?
		// damage *= attackerStats.DamageMultiplier (if it existed)
	}

	return hit, damage
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

// -------------------------------------------------------------------------
// Hard Scifi Combat Math
// -------------------------------------------------------------------------

// CalculateTurretTracking determines the hit chance of a turret based on angular velocity.
// turretTracking: The weapon's tracking speed (rad/sec).
// targetRadius: The target's signature radius (m).
// distance: Range to target (m).
// transversalVelocity: Speed of target perpendicular to the shooter (m/s).
//
// Formula: Chance = 0.5 ^ ( ( (Transversal / Range) / Tracking ) ^ 2 )
// Note: This matches standard "EVE-like" tracking mechanics.
func CalculateTurretTracking(turretTracking float64, targetRadius float64, distance float64, transversalVelocity float64) float64 {
	if distance <= 0 {
		return 1.0 // Point blank
	}

	angularVelocity := transversalVelocity / distance // rad/sec

	// Standard Tracking Formula
	// The exponent part compares the angular tracking demand vs the gun's capability.
	// We use targetRadius usually to mitigate tracking (Signature Resolution / Target Radius).
	// Assuming `turretTracking` includes the resolution factor or is a raw score.
	// Prompt says: "Hit chance = Weapon Slew Rate vs Target Angular Velocity".
	// Let's implement the canonical formula:
	// Chance = 0.5 ^ ( (Angular / Tracking) * (SignatureResolution / TargetRadius) ) ^ 2
	// For this MVP function, we lack SignatureResolution input, so we assume 1:1 or ignore it,
	// OR we assume turretTracking is the *effective* tracking against a standard target.
	// We will follow the prompt's implied simple comparison but use the curve.

	// Simple: Chance = 0.5 ^ ( (Angular / Tracking) ^ 2 )
	exponent := math.Pow(angularVelocity/turretTracking, 2)
	chance := math.Pow(0.5, exponent)

	return chance
}

// CalculateMissileDamage calculates applied damage based on explosion physics.
// baseDamage: Warhead damage.
// explosionRadius: Size of the explosion (m).
// explosionVel: Expansion velocity of the explosion (m/s).
// targetRadius: Target's signature radius (m).
// targetVel: Target's speed (m/s).
func CalculateMissileDamage(baseDamage float64, explosionRadius float64, explosionVel float64, targetRadius float64, targetVel float64) float64 {
	// Factor 1: Signature (Small targets take less from big booms)
	sigFactor := targetRadius / explosionRadius

	// Factor 2: Velocity (Fast targets outrun the boom)
	// Formula: (ExplosionVel / TargetVel) * (TargetRadius / ExplosionRadius)
	// We take the minimum of 1.0, SigFactor, and VelFactor.

	// Avoid div by zero
	if targetVel <= 0 {
		targetVel = 0.001
	}
	if explosionRadius <= 0 {
		explosionRadius = 1.0
	}

	velFactor := (explosionVel / targetVel) * sigFactor

	// Applied factor is the smallest of: 1.0 (Full hit), SigFactor, or VelFactor
	appliedFactor := math.Min(1.0, math.Min(sigFactor, velFactor))

	return baseDamage * appliedFactor
}

// CalculateLockTime determines how long it takes to lock a target.
// sourceScanRes: The attacker's sensor strength/resolution (mm).
// targetSCS: The target's sensor cross-section (m^2).
//
// Prompt Formula: Lock Time = Scanner Res / Target SCS
func CalculateLockTime(sourceScanRes float64, targetSCS float64) float64 {
	if targetSCS <= 0 {
		return 999.0 // Cannot lock stealth target
	}
	// Per prompt instruction.
	// Note: This implies that higher ScanRes makes locking SLOWER, or ScanRes is a "Scan Delay" value.
	return sourceScanRes / targetSCS
}

// ResolveGroundDamage calculates damage for ground combat based on stats.
func ResolveGroundDamage(attacker *Player, target *Player) float64 {
	// 1. Get Weapon Damage
	weaponDamage := 5.0 // Default Unarmed
	var itemID string

	if attacker.GroundGear.PrimaryWeapon != nil {
		itemID = attacker.GroundGear.PrimaryWeapon.ItemID
	} else if attacker.GroundGear.Sidearm != nil {
		itemID = attacker.GroundGear.Sidearm.ItemID
	}

	if itemDef, ok := Items[itemID]; ok {
		if val, ok := itemDef.Stats["damage"]; ok {
			weaponDamage = val
		}
	}

	// 2. Get Target Defense
	defense := 0.0
	if classDef, ok := Classes[target.ClassID]; ok {
		defense = float64(classDef.Stats.Defense)
	}

	// 3. Calculate Mitigation (MOBA Formula)
	// Damage Multiplier = 100 / (100 + Defense)
	multiplier := 100.0 / (100.0 + defense)

	finalDamage := weaponDamage * multiplier
	return finalDamage
}
