package game

import (
	"fmt"
	"math/rand"
	"time"
)

// MarketService handles buying and selling of items.
type MarketService struct {
	Repo PlayerRepository
}

func NewMarketService(repo PlayerRepository) *MarketService {
	return &MarketService{Repo: repo}
}

// BuyItem handles the purchase of an item by a player.
// It checks currency, deducts cost, generates item (with traits if organic), and adds to inventory.
func (s *MarketService) BuyItem(player *Player, itemID string) error {
	// 1. Validate Item Exists
	itemDef, ok := Items[itemID]
	if !ok {
		return fmt.Errorf("item not found: %s", itemID)
	}

	// 2. Validate Funds
	if player.Solium < itemDef.Price {
		return fmt.Errorf("insufficient funds: have %d, need %d", player.Solium, itemDef.Price)
	}

	// 3. Deduct Cost
	player.Solium -= itemDef.Price

	// 4. Generate Item
	newItem := ItemStack{
		ItemID: itemID,
		Count:  1,
	}

	if itemDef.IsOrganic {
		newItem.Data = generateOrganicTraits(itemID)
	}

	// 5. Add to Inventory
	player.AddItem(newItem)

	// 6. Persist State
	if err := s.Repo.SavePlayerState(player); err != nil {
		// Rollback (In-memory only, DB is consistent via transaction ideally, but here manual)
		player.Solium += itemDef.Price
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}

// generateOrganicTraits creates random properties for bio-grafts.
func generateOrganicTraits(itemID string) map[string]interface{} {
	traits := make(map[string]interface{})
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Base Stats
	traits["decay"] = 0.0

	// Quality Roll
	roll := r.Intn(100)
	if roll > 90 {
		traits["quality"] = "pristine"
		traits["compatibility"] = 1.0 // 100% match
	} else if roll > 50 {
		traits["quality"] = "standard"
		traits["compatibility"] = 0.8
	} else {
		traits["quality"] = "questionable"
		traits["compatibility"] = 0.5
		traits["decay"] = float64(r.Intn(20)) // Starts pre-decayed
	}

	return traits
}

// SurgeryService handles the installation and removal of bio-grafts and modules.
type SurgeryService struct {
	Repo PlayerRepository
}

func NewSurgeryService(repo PlayerRepository) *SurgeryService {
	return &SurgeryService{Repo: repo}
}

// GraftOrgan installs an item from the inventory into a specific ship slot.
// inventoryIndex: The index of the item in player.Inventory.
// slotType: "high_slots", "mid_slots", "low_slots", "rigs", "injectors".
// slotIndex: The index of the slot in the respective array.
func (s *SurgeryService) GraftOrgan(player *Player, inventoryIndex int, slotType string, slotIndex int) error {
	// 1. Validate Inventory Index
	if inventoryIndex < 0 || inventoryIndex >= len(player.Inventory) {
		return fmt.Errorf("invalid inventory index: %d", inventoryIndex)
	}

	if slotIndex < 0 {
		return fmt.Errorf("invalid slot index: %d", slotIndex)
	}

	if slotType == "" {
		return fmt.Errorf("invalid slot type")
	}

	// 2. Retrieve Item
	itemToInstall := player.Inventory[inventoryIndex]

	// Check Static Data for Slot Type
	itemDef, ok := Items[itemToInstall.ItemID]
	if !ok {
		return fmt.Errorf("unknown item: %s", itemToInstall.ItemID)
	}

	// Map generic slot types (high_slots) to specific definition types if needed.
	if itemDef.SlotType != slotType {
		return fmt.Errorf("item %s fits in %s, not %s", itemDef.Name, itemDef.SlotType, slotType)
	}

	// 3. Access the Target Slot Array
	var targetSlots []Slot
	switch slotType {
	case "high_slots":
		targetSlots = player.Ship.HighSlots
	case "mid_slots":
		targetSlots = player.Ship.MidSlots
	case "low_slots":
		targetSlots = player.Ship.LowSlots
	case "rigs":
		targetSlots = player.Ship.Rigs
	case "injectors":
		targetSlots = player.Ship.Injectors
	default:
		return fmt.Errorf("unknown ship slot type: %s", slotType)
	}

	// 4. Handle Swap (If slot is occupied)
	var removedItem *ItemStack
	if slotIndex < len(targetSlots) {
		// Slot exists
		existingSlot := targetSlots[slotIndex]
		if existingSlot.Module != nil {
			removedItem = existingSlot.Module
		}
	} else if slotIndex == len(targetSlots) {
		// Appending to new slot
	} else {
		return fmt.Errorf("slot index out of bounds (gap in slots)")
	}

	// 5. Perform the Operation

	// Remove from Inventory by Index
	player.Inventory = append(player.Inventory[:inventoryIndex], player.Inventory[inventoryIndex+1:]...)

	// Add removed item back to inventory (Infinite Inventory)
	if removedItem != nil {
		player.AddItem(*removedItem)
	}

	// Install into Ship
	newSlot := Slot{Module: &itemToInstall}

	switch slotType {
	case "high_slots":
		player.Ship.HighSlots = updateSlotSlice(player.Ship.HighSlots, slotIndex, newSlot)
	case "mid_slots":
		player.Ship.MidSlots = updateSlotSlice(player.Ship.MidSlots, slotIndex, newSlot)
	case "low_slots":
		player.Ship.LowSlots = updateSlotSlice(player.Ship.LowSlots, slotIndex, newSlot)
	case "rigs":
		player.Ship.Rigs = updateSlotSlice(player.Ship.Rigs, slotIndex, newSlot)
	case "injectors":
		player.Ship.Injectors = updateSlotSlice(player.Ship.Injectors, slotIndex, newSlot)
	}

	// 6. Persist
	if err := s.Repo.SavePlayerState(player); err != nil {
		return fmt.Errorf("failed to save surgery result: %w", err)
	}

	return nil
}

// Helper to handle the "Append or Replace" logic for slots
func updateSlotSlice(slots []Slot, index int, newSlot Slot) []Slot {
	if index < len(slots) {
		slots[index] = newSlot
		return slots
	}
	// If index is exactly len, append.
	if index == len(slots) {
		return append(slots, newSlot)
	}
	// Should be caught by validation, but return original if out of bounds
	return slots
}
