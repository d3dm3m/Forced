package game

import (
	"testing"
)

func TestLoadGameData(t *testing.T) {
	// 1. Ensure map is empty before load (though it's global, this is first test)
	if len(Classes) != 0 {
		t.Log("Classes map was not empty, assuming already loaded or state carried over.")
	}

	// 2. Load Data
	LoadGameData()

	// 3. Assertions
	// Check Marine
	marine, ok := Classes["marine"]
	if !ok {
		t.Fatalf("Expected 'marine' class to be loaded")
	}
	if marine.Stats.Health != 150 {
		t.Errorf("Expected Marine HP 150, got %d", marine.Stats.Health)
	}

	// Check Sapper
	sapper, ok := Classes["sapper"]
	if !ok {
		t.Fatalf("Expected 'sapper' class to be loaded")
	}
	if sapper.Slots != 6 {
		t.Errorf("Expected Sapper Slots 6, got %d", sapper.Slots)
	}

	// Check Frigate
	frigate, ok := Classes["frigate_ace"]
	if !ok {
		t.Fatalf("Expected 'frigate_ace' class to be loaded")
	}
	if frigate.Type != "Space" {
		t.Errorf("Expected Frigate Type 'Space', got %s", frigate.Type)
	}

	// Check Item
	shotgun, ok := Items["auto_shotgun"]
	if !ok {
		t.Fatalf("Expected 'auto_shotgun' item to be loaded")
	}
	if shotgun.Mass != 5.0 {
		t.Errorf("Expected Shotgun Mass 5.0, got %f", shotgun.Mass)
	}
}
