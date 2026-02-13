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
	// Check Breacher
	breacher, ok := Classes["breacher"]
	if !ok {
		t.Fatalf("Expected 'breacher' class to be loaded")
	}
	if breacher.Stats.Health != 200 {
		t.Errorf("Expected Breacher HP 200, got %d", breacher.Stats.Health)
	}
	if breacher.Stats.Torque != 8 {
		t.Errorf("Expected Breacher Torque 8, got %d", breacher.Stats.Torque)
	}

	// Check Sapper (Syndicate)
	sapper, ok := Classes["sapper"]
	if !ok {
		t.Fatalf("Expected 'sapper' class to be loaded")
	}
	if sapper.Slots != 8 {
		t.Errorf("Expected Sapper Slots 8, got %d", sapper.Slots)
	}

	// Check Null-Walker
	walker, ok := Classes["null_walker"]
	if !ok {
		t.Fatalf("Expected 'null_walker' class to be loaded")
	}
	if walker.Stats.Flux != 10 {
		t.Errorf("Expected Null-Walker Flux 10, got %d", walker.Stats.Flux)
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
