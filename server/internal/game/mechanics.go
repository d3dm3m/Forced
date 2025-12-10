package game

// DerivedStats represents the calculated total statistics of a ship.
type DerivedStats struct {
	MaxHealth    float64
	Speed        float64
	SensorRange  float64
	BioCapacity  int
	CurrentBioLoad int
	RejectionRate  float64 // Damage per second
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
		stats.MaxHealth = getFloat(shipDef.Stats, "base_hull", 1000.0)
		stats.Speed = getFloat(shipDef.Stats, "base_speed", 100.0) // "speed" might not be in JSON, check defaults
		// Checking ships.json from context: "base_hull": 1000, "base_shield": 500, "capacitor": 100
		// No speed in ships.json?
		// Checking classes.json: "Frigate Ace" has stats: { "speed": 50 }
		// Maybe base stats come from CLASS?
		// The prompt says "Start with the Ship's base stats (from ships.json)".
		// If ships.json lacks speed, I'll default it.
		stats.BioCapacity = getInt(shipDef.Stats, "bio_capacity", 50) // Default capacity
	} else {
		// Absolute fallback
		stats.MaxHealth = 1000
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
		case "sensor_range":
			stats.SensorRange += effectiveValue
		case "bio_capacity":
			stats.BioCapacity += int(effectiveValue)
		}
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
