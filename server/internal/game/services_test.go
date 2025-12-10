package game

import (
	"testing"
)

// Mock Repo for Service Tests
type MockPlayerRepo struct {
	P *Player
}

func (m *MockPlayerRepo) CreatePlayer(u, c string) (*Player, error) { return nil, nil }
func (m *MockPlayerRepo) LoadPlayer(u string) (*Player, error)      { return m.P, nil }
func (m *MockPlayerRepo) SavePlayerState(p *Player) error {
	m.P = p
	return nil
}

func TestMarketService_BuyItem(t *testing.T) {
	// Setup
	LoadGameData() // Ensure static data is loaded
	player := &Player{
		Username:  "Buyer",
		Solium:    10000,
		Inventory: []ItemStack{},
	}
	repo := &MockPlayerRepo{P: player}
	service := NewMarketService(repo)

	// 1. Buy Simple Item
	err := service.BuyItem(player, "auto_shotgun") // Price 500
	if err != nil {
		t.Errorf("Failed to buy simple item: %v", err)
	}
	if player.Solium != 9500 {
		t.Errorf("Expected 9500 Solium, got %d", player.Solium)
	}
	if len(player.Inventory) != 1 || player.Inventory[0].ItemID != "auto_shotgun" {
		t.Errorf("Inventory mismatch after buy")
	}

	// 2. Buy Organic Item (Check Traits)
	err = service.BuyItem(player, "synthetic_heart") // Price 5000, Organic
	if err != nil {
		t.Errorf("Failed to buy organic item: %v", err)
	}
	if player.Solium != 4500 {
		t.Errorf("Expected 4500 Solium, got %d", player.Solium)
	}

	lastItem := player.Inventory[len(player.Inventory)-1]
	if lastItem.ItemID != "synthetic_heart" {
		t.Errorf("Expected synthetic_heart, got %s", lastItem.ItemID)
	}
	if len(lastItem.Data) == 0 {
		t.Errorf("Expected organic traits (Data map), got empty")
	}
	if _, ok := lastItem.Data["quality"]; !ok {
		t.Errorf("Expected 'quality' trait")
	}

	// 3. Insufficient Funds
	err = service.BuyItem(player, "synthetic_heart") // Costs 5000, have 4500
	if err == nil {
		t.Errorf("Expected error for insufficient funds, got nil")
	}
}

func TestSurgeryService_GraftOrgan(t *testing.T) {
	// Setup
	LoadGameData()
	player := &Player{
		Username: "Surgeon",
		Ship: ShipLayout{
			HighSlots: []Slot{}, // Start empty
		},
		Inventory: []ItemStack{
			{ItemID: "mining_laser", Count: 1}, // Index 0, High Slot
			{ItemID: "laser_cannon", Count: 1}, // Index 1, High Slot
		},
	}
	repo := &MockPlayerRepo{P: player}
	service := NewSurgeryService(repo)

	// 1. Install into Empty Slot 0
	// Install "mining_laser" (Index 0) into "high_slots" (Slot 0)
	err := service.GraftOrgan(player, 0, "high_slots", 0)
	if err != nil {
		t.Errorf("Failed to graft organ: %v", err)
	}

	// Verify Inventory (Item removed)
	if len(player.Inventory) != 1 {
		t.Errorf("Expected inventory size 1, got %d", len(player.Inventory))
	}
	if player.Inventory[0].ItemID != "laser_cannon" {
		t.Errorf("Wrong item remaining in inventory")
	}

	// Verify Ship (Item installed)
	if len(player.Ship.HighSlots) != 1 {
		t.Errorf("Expected 1 high slot, got %d", len(player.Ship.HighSlots))
	}
	if player.Ship.HighSlots[0].Module.ItemID != "mining_laser" {
		t.Errorf("Expected mining_laser in slot 0")
	}

	// 2. Swap Logic
	// Install "laser_cannon" (Index 0 now) into "high_slots" (Slot 0) -> Should swap mining_laser out
	err = service.GraftOrgan(player, 0, "high_slots", 0)
	if err != nil {
		t.Errorf("Failed to swap organ: %v", err)
	}

	// Verify Ship
	if player.Ship.HighSlots[0].Module.ItemID != "laser_cannon" {
		t.Errorf("Expected laser_cannon in slot 0 after swap")
	}

	// Verify Inventory (mining_laser returned)
	if len(player.Inventory) != 1 {
		t.Errorf("Expected inventory size 1, got %d", len(player.Inventory))
	}
	if player.Inventory[0].ItemID != "mining_laser" {
		t.Errorf("Expected mining_laser returned to inventory, got %s", player.Inventory[0].ItemID)
	}
}
