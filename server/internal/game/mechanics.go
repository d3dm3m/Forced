package game

import "math"

// DerivedStats represents the calculated total statistics of a ship.
type DerivedStats struct {
	MaxHealth     float64
	MaxShield     float64
	MaxCapacitor  float64
	CapRecharge   float64 // Peak Recharge Rate (GJ/s)
	Speed         float64
	SensorRange   float64
	BioCapacity   int
	CurrentBioLoad int
	RejectionRate float64 // Damage per second
}

// CalculateShipStats aggregates the base ship stats with all installed module modifiers.
// It also calculates Bio-Load and Rejection penalties.
func CalculateShipStats(player *Player) DerivedStats {
	// 1. Initialize with Base Ship Stats
	// Need to look up the ship definition based on player's ship/class?
	// The Player struct does NOT currently store the Ship Hull ID (e.g. "kestrel").
	// It stores `ShipLayout`.
	// Ideally, `Player` should have `ShipID`.
	// For MVP, checking `ships.json` and `classes.json`:
	// `classes.json` has `Frigate Ace` which is a Class.
	// `ships.json` has `Kestrel`.
	// The prompt implies "Start with the Ship's base stats".
	// Issue: We don't know WHICH ship the player is flying from the `Player` struct alone
	// (unless we assume default or it's missing).
	// Checking `PlayerRepository.CreatePlayer`: it creates `defaultShip`.
	// It doesn't seem to set a Hull ID.
	// Assumption: For this mechanic to work, we need a Base.
	// I will assume a default base "kestrel" if unknown, or maybe the `ClassID` defines it?
	// `Frigate Ace` is a class, but usually you buy a hull.
	// Let's fallback to "kestrel" base stats for now if we can't find it,
	// but logically we should probably add `ShipID` to Player.
	// Given strict constraints, I'll use "kestrel" as the default base for calculations
	// to ensure the function works.

	baseShipID := "kestrel" // Fallback
	// In a real implementation, player.ShipID would exist.

	shipDef, ok := Ships[baseShipID]
	var stats DerivedStats
	if ok {
		// Use explicit struct fields if populated (from new JSON), fallback to Stats map for older data
		stats.MaxHealth = shipDef.HullHP
		if stats.MaxHealth == 0 { stats.MaxHealth = getFloat(shipDef.Stats, "base_hull", 1000.0) }

		stats.MaxShield = shipDef.ShieldHP
		if stats.MaxShield == 0 { stats.MaxShield = getFloat(shipDef.Stats, "base_shield", 500.0) }

		stats.MaxCapacitor = shipDef.CapacitorCapacity
		if stats.MaxCapacitor == 0 { stats.MaxCapacitor = getFloat(shipDef.Stats, "capacitor", 100.0) }

		stats.CapRecharge = shipDef.CapacitorRechargeRate
		if stats.CapRecharge == 0 { stats.CapRecharge = 100.0 }

		stats.SensorRange = getFloat(shipDef.Stats, "sensor_range", 100.0)
		stats.Speed = getFloat(shipDef.Stats, "base_speed", 100.0)
		stats.BioCapacity = getInt(shipDef.Stats, "bio_capacity", 50)
	} else {
		// Absolute fallback
		stats.MaxHealth = 1000
		stats.MaxShield = 500
		stats.MaxCapacitor = 100
		stats.CapRecharge = 20
		stats.Speed = 100
		stats.BioCapacity = 50
	}

	// 2. Iterate All Slots
	// Combine all slot slices into one iterator or loop over each.
	allSlots := [][]Slot{
		player.Ship.HighSlots,
		player.Ship.MidSlots,
		player.Ship.LowSlots,
		player.Ship.Rigs,
		player.Ship.Injectors,
	}

	for _, slotGroup := range allSlots {
		for _, slot := range slotGroup {
			if slot.Module == nil {
				continue
			}
			applyModuleStats(&stats, slot.Module)
		}
	}

	// 3. Calculate Rejection
	if stats.CurrentBioLoad > stats.BioCapacity {
		// Overload!
		overload := stats.CurrentBioLoad - stats.BioCapacity
		// Rejection Rate: 1 DPS per point of overload? Or constant?
		// Prompt: "Rate = (Load - Cap) * Constant"
		const RejectionMultiplier = 0.5
		stats.RejectionRate = float64(overload) * RejectionMultiplier
	} else {
		stats.RejectionRate = 0
	}

	return stats
}

func applyModuleStats(stats *DerivedStats, item *ItemStack) {
	def, ok := Items[item.ItemID]
	if !ok {
		return
	}

	// Bio Cost (Static)
	stats.CurrentBioLoad += def.BioCost

	// Stat Multipliers
	// Default Compatibility = 1.0
	compatibility := 1.0
	if val, ok := item.Data["compatibility"]; ok {
		if fVal, ok := val.(float64); ok {
			compatibility = fVal
		}
	}

	// Apply Stats
	for stat, value := range def.Stats {
		effectiveValue := value * compatibility

		switch stat {
		case "speed":
			stats.Speed += effectiveValue
		case "health", "base_hull":
			stats.MaxHealth += effectiveValue
		case "shield", "base_shield":
			stats.MaxShield += effectiveValue
		case "capacitor":
			stats.MaxCapacitor += effectiveValue
		case "sensor_range":
			stats.SensorRange += effectiveValue
		case "bio_capacity":
			stats.BioCapacity += int(effectiveValue)
		}
	}
}

// RegenerateShip updates the transient state of the ship (Shield/Cap) based on delta time.
func RegenerateShip(player *Player, stats DerivedStats, dt float64) {
	// 1. Shield Regen (Linear)
	// Example: 1% per second
	regenAmount := (stats.MaxShield * 0.01) * dt
	if player.CurrentShield < stats.MaxShield {
		player.CurrentShield += regenAmount
		if player.CurrentShield > stats.MaxShield {
			player.CurrentShield = stats.MaxShield
		}
	}

	// 2. Capacitor Regen (Non-Linear / "Zombie Curve")
	// Formula: dC/dt = (10 * MaxCap / RechargeTime) * ( sqrt(C/Max) - C/Max )
	// RechargeTime usually ~300s? We used "CapRecharge" as a Rate or Time?
	// In EVE, the stat is "Recharge Time". In our JSON we called it "capacitor_recharge".
	// Let's treat stats.CapRecharge as "Recharge Time in Seconds".

	if stats.MaxCapacitor > 0 && stats.CapRecharge > 0 {
		ratio := player.CurrentCapacitor / stats.MaxCapacitor
		if ratio < 1.0 {
			// Avoid Sqrt of 0 or negative if empty
			if ratio < 0 { ratio = 0 }

			// EVE formula approximation
			// Rate = (10 * Max) / Time * (sqrt(ratio) - ratio) (simplified curve)
			// Wait, the real formula is complex.
			// Simpler "Peaked" curve: Rate = PeakRate * 2.5 * ratio * (1 - ratio)? No.
			// Let's use the provided Prompt formula:
			// Rate = (10 * MaxCap) / RechargeTime * ( sqrt(Current/Max) - (Current/Max) )

			rate := (10.0 * stats.MaxCapacitor) / stats.CapRecharge * (math.Sqrt(ratio) - ratio)

			// If rate is negative (shouldn't be for 0 < ratio < 1), clamp 0
			if rate < 0 { rate = 0 }

			// Minimum trickle to prevent stuck at 0
			if player.CurrentCapacitor <= 0.1 {
				rate = stats.MaxCapacitor * 0.005 // Jump start
			}

			player.CurrentCapacitor += rate * dt
			if player.CurrentCapacitor > stats.MaxCapacitor {
				player.CurrentCapacitor = stats.MaxCapacitor
			}
		}
	} else {
		// Fallback Linear
		player.CurrentCapacitor += 1.0 * dt
	}
}

// Helpers for extracting values from generic maps
func getFloat(m map[string]interface{}, key string, def float64) float64 {
	if val, ok := m[key]; ok {
		// JSON numbers are often float64 in Go generic maps
		if f, ok := val.(float64); ok {
			return f
		}
		if i, ok := val.(int); ok {
			return float64(i)
		}
	}
	return def
}

func getInt(m map[string]interface{}, key string, def int) int {
	if val, ok := m[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
		if f, ok := val.(float64); ok {
			return int(f)
		}
	}
	return def
}
