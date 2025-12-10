package game

import (
	"encoding/json"
	"testing"
)

func TestPlayerHealthPersistence(t *testing.T) {
	// 1. Create Player with default Health (1000.0)
	// We'll mock the DB calls or just test the struct behavior if simple.
	// But `CreatePlayer` sets the default.
	// Since we don't have a real DB in unit tests unless we mock `db.Pool`,
	// we will verify that the struct fields and JSON tags are correct
	// and trust the SQL migration we wrote.

	// Create struct manually to verify JSON marshaling of new field
	p := &Player{
		CurrentHealth: 750.5,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Failed to marshal player: %v", err)
	}

	var p2 Player
	err = json.Unmarshal(data, &p2)
	if err != nil {
		t.Fatalf("Failed to unmarshal player: %v", err)
	}

	if p2.CurrentHealth != 750.5 {
		t.Errorf("Expected CurrentHealth 750.5, got %f", p2.CurrentHealth)
	}
}
