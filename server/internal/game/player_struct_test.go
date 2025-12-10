package game

import (
	"encoding/json"
	"testing"
)

// TestInventoryStacking tests that simple items stack and complex items do not.
func TestInventoryStacking(t *testing.T) {
	player := &Player{
		Inventory: []ItemStack{},
	}

	// 1. Add simple item (Ammo)
	simpleItem := ItemStack{ItemID: "ammo_50cal", Count: 100}
	player.AddItem(simpleItem)

	if len(player.Inventory) != 1 {
		t.Errorf("Expected inventory size 1, got %d", len(player.Inventory))
	}
	if player.Inventory[0].Count != 100 {
		t.Errorf("Expected count 100, got %d", player.Inventory[0].Count)
	}

	// 2. Add same simple item (Stacking check)
	player.AddItem(simpleItem)
	if len(player.Inventory) != 1 {
		t.Errorf("Expected inventory size 1 after stacking, got %d", len(player.Inventory))
	}
	if player.Inventory[0].Count != 200 {
		t.Errorf("Expected count 200, got %d", player.Inventory[0].Count)
	}

	// 3. Add complex item (Bio-Organ with data)
	complexItem := ItemStack{
		ItemID: "heart_organ",
		Count:  1,
		Data: map[string]interface{}{
			"decay": 0.5,
			"tier":  2,
		},
	}
	player.AddItem(complexItem)

	if len(player.Inventory) != 2 {
		t.Errorf("Expected inventory size 2, got %d", len(player.Inventory))
	}

	// 4. Add SAME complex item again (Should NOT stack because it has data)
	// Note: Even if data is identical, our logic says "if len(Data) > 0 { append }"
	player.AddItem(complexItem)

	if len(player.Inventory) != 3 {
		t.Errorf("Expected inventory size 3 (non-stacking), got %d", len(player.Inventory))
	}
}

// TestJSONMarshaling verifies that the new structures marshal/unmarshal correctly.
func TestJSONMarshaling(t *testing.T) {
	// Construct a full player object
	p := &Player{
		Username: "TestUser",
		Inventory: []ItemStack{
			{ItemID: "ammo", Count: 50},
			{ItemID: "organ", Count: 1, Data: map[string]interface{}{"quality": "rotting"}},
		},
		Ship: ShipLayout{
			HighSlots: []Slot{
				{Module: &ItemStack{ItemID: "laser_cannon", Count: 1}},
			},
			MidSlots: []Slot{}, // Empty
		},
		GroundGear: GroundGear{
			PrimaryWeapon: &ItemStack{ItemID: "rifle", Count: 1},
			Helmet:        &ItemStack{ItemID: "visored_helm", Count: 1},
			Talents:       map[string]int{"xenobiology": 5},
		},
	}

	// Marshal
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Failed to marshal player: %v", err)
	}

	// Unmarshal
	var p2 Player
	err = json.Unmarshal(data, &p2)
	if err != nil {
		t.Fatalf("Failed to unmarshal player: %v", err)
	}

	// Verify Data
	if len(p2.Inventory) != 2 {
		t.Errorf("Inventory length mismatch. Expected 2, got %d", len(p2.Inventory))
	}

	// Check complex item data
	val, ok := p2.Inventory[1].Data["quality"]
	if !ok || val != "rotting" {
		t.Errorf("Failed to retrieve complex item data. Got %v", val)
	}

	// Check Ship Slots
	if len(p2.Ship.HighSlots) != 1 {
		t.Errorf("HighSlots length mismatch. Expected 1, got %d", len(p2.Ship.HighSlots))
	}
	if p2.Ship.HighSlots[0].Module.ItemID != "laser_cannon" {
		t.Errorf("Ship module ID mismatch. Expected laser_cannon, got %v", p2.Ship.HighSlots[0].Module.ItemID)
	}

	// Check Ground Gear
	if p2.GroundGear.PrimaryWeapon.ItemID != "rifle" {
		t.Errorf("PrimaryWeapon mismatch. Expected rifle, got %v", p2.GroundGear.PrimaryWeapon.ItemID)
	}
	if p2.GroundGear.Talents["xenobiology"] != 5 {
		t.Errorf("Talent mismatch. Expected 5, got %v", p2.GroundGear.Talents["xenobiology"])
	}
}
