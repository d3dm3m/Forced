package game

import (
	"testing"
)

func TestCalculateShipStats_BioLoad(t *testing.T) {
	LoadGameData() // Ensure static data

	// Setup Player with a ship
	player := &Player{
		Ship: ShipLayout{
			HighSlots: []Slot{},
			MidSlots:  []Slot{},
			LowSlots:  []Slot{},
			Rigs:      []Slot{},
			Injectors: []Slot{},
		},
	}

	// 1. Baseline
	baseStats := CalculateShipStats(player)
	if baseStats.CurrentBioLoad != 0 {
		t.Errorf("Expected 0 BioLoad, got %d", baseStats.CurrentBioLoad)
	}

	// 2. Install Standard Heart (BioCost 10, Speed 20)
	stdHeart := ItemStack{
		ItemID: "synthetic_heart",
		Count: 1,
		Data: map[string]interface{}{"compatibility": 1.0},
	}
	// Manually inject into slot for test (skipping SurgeryService to test logic directly)
	player.Ship.HighSlots = append(player.Ship.HighSlots, Slot{Module: &stdHeart})

	stats1 := CalculateShipStats(player)
	if stats1.CurrentBioLoad != 10 {
		t.Errorf("Expected 10 BioLoad, got %d", stats1.CurrentBioLoad)
	}
	if stats1.Speed <= baseStats.Speed {
		t.Errorf("Expected Speed increase (Base %f -> New %f)", baseStats.Speed, stats1.Speed)
	}
	expectedSpeed := baseStats.Speed + 20.0
	if stats1.Speed != expectedSpeed {
		t.Errorf("Expected Speed %f, got %f", expectedSpeed, stats1.Speed)
	}

	// 3. Install Low Compatibility Heart (BioCost 10, Speed 20 * 0.5 = 10)
	// Add to another slot
	badHeart := ItemStack{
		ItemID: "synthetic_heart",
		Count: 1,
		Data: map[string]interface{}{"compatibility": 0.5},
	}
	player.Ship.MidSlots = append(player.Ship.MidSlots, Slot{Module: &badHeart})

	stats2 := CalculateShipStats(player)
	// BioLoad should sum fully (10 + 10 = 20)
	if stats2.CurrentBioLoad != 20 {
		t.Errorf("Expected 20 BioLoad, got %d", stats2.CurrentBioLoad)
	}
	// Speed should increase by 10 (Total +30 from base)
	expectedSpeed2 := baseStats.Speed + 20.0 + 10.0
	if stats2.Speed != expectedSpeed2 {
		t.Errorf("Expected Speed %f, got %f", expectedSpeed2, stats2.Speed)
	}

	// 4. Overload (BioCapacity is 50 default in mechanics.go fallback, or 0 if defined elsewhere)
	// Let's add a massive item: Void Heart (Cost 100)
	voidHeart := ItemStack{
		ItemID: "void_heart",
		Count: 1,
		Data: map[string]interface{}{"compatibility": 1.0},
	}
	player.Ship.LowSlots = append(player.Ship.LowSlots, Slot{Module: &voidHeart})

	stats3 := CalculateShipStats(player)
	totalLoad := 10 + 10 + 100 // 120
	if stats3.CurrentBioLoad != totalLoad {
		t.Errorf("Expected 120 BioLoad, got %d", stats3.CurrentBioLoad)
	}

	// Check Rejection
	// Capacity 50. Overload = 120 - 50 = 70.
	// Rate = 70 * 0.5 = 35.0
	if stats3.RejectionRate <= 0 {
		t.Errorf("Expected Rejection Rate > 0, got %f", stats3.RejectionRate)
	}
	if stats3.RejectionRate != 35.0 {
		t.Errorf("Expected Rejection Rate 35.0, got %f (Cap %d)", stats3.RejectionRate, stats3.BioCapacity)
	}
}
