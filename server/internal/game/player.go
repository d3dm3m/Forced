package game

import (
	"context"
	"encoding/json"
	"fmt"
	"biohorror/internal/db"
	"time"

	"github.com/jackc/pgx/v5"
)

// ItemStack represents an item in the inventory or a slot.
type ItemStack struct {
	ItemID string                 `json:"item_id"`
	Count  int                    `json:"count"`
	Data   map[string]interface{} `json:"data,omitempty"` // If set, item is non-stackable
}

// Slot represents a tiered slot in the ship layout.
type Slot struct {
	Module *ItemStack `json:"module,omitempty"`
}

// ShipLayout represents the organic ship layout, stored as JSONB
type ShipLayout struct {
	HighSlots []Slot `json:"high_slots"` // Cranial Mounts
	MidSlots  []Slot `json:"mid_slots"`  // Vascular Systems
	LowSlots  []Slot `json:"low_slots"`  // Viscera/Structure
	Rigs      []Slot `json:"rigs"`       // Genetic Splices
	Injectors []Slot `json:"injectors"`  // Injection Ports
}

// GroundGear represents the player's equipment and talents.
type GroundGear struct {
	// Layer 1: The Loadout (External / Handheld)
	PrimaryWeapon *ItemStack `json:"primary_weapon,omitempty"`
	Sidearm       *ItemStack `json:"sidearm,omitempty"`

	// Layer 2: The Shell (Wearable Protection)
	Helmet  *ItemStack `json:"helmet,omitempty"`
	Exosuit *ItemStack `json:"exosuit,omitempty"`

	// Layer 3: The Meat/Chrome (Internal Augments)
	Cranial      *ItemStack `json:"cranial,omitempty"`
	Thoracic     *ItemStack `json:"thoracic,omitempty"`
	Manipulators *ItemStack `json:"manipulators,omitempty"`
	Locomotors   *ItemStack `json:"locomotors,omitempty"`

	// Skill Tree
	Talents map[string]int `json:"talents"`
}

// IndustrialStats represents the core bio-mechanical attributes.
type IndustrialStats struct {
	Torque  int `json:"torque"`  // Physical strength/load
	Compute int `json:"compute"` // Processing power
	Synapse int `json:"synapse"` // Reaction speed
	Flux    int `json:"flux"`    // Exotic energy manipulation
}

// Motherboard represents the grid for installing Chipsets.
type Motherboard struct {
	Slots    int         `json:"slots"` // Total available grid size
	Chips    []ItemStack `json:"chips"` // Installed chips
	Stats    IndustrialStats `json:"stats"` // Aggregated stats
}

// Player represents the player state
type Player struct {
	ID            string         `json:"id"`
	Username      string         `json:"username"`
	ClassID       string         `json:"class_id"`
	PositionX     float64        `json:"position_x"`
	PositionY     float64        `json:"position_y"`
	Inventory     []ItemStack    `json:"inventory"`      // Changed to slice of structs
	Ship          ShipLayout     `json:"ship_layout"`    // JSONB
	GroundGear    GroundGear     `json:"ground_gear"`    // JSONB
	Motherboard   Motherboard    `json:"motherboard"`    // JSONB (New in Sprint 21)
	Solium           int            `json:"solium"`            // Currency
	Skills           map[string]int `json:"skills"`            // Ship/Space Skills
	CurrentHealth    float64        `json:"current_health"`    // Ship Health (Space) or Player Health (Ground)
	CurrentShield    float64        `json:"current_shield"`    // Transient
	CurrentCapacitor float64        `json:"current_capacitor"` // Transient
	SystemID         string         `json:"system_id"`         // Current Star System (e.g. "Sol-0")
	CreatedAt        time.Time      `json:"created_at"`
}

type PlayerRepository interface {
	CreatePlayer(username string, classID string) (*Player, error)
	LoadPlayer(username string) (*Player, error)
	SavePlayerState(player *Player) error
}

// AddItem adds an item to the player's inventory.
// If the item has Data, it is treated as unique and appended.
// If the item has no Data, it attempts to stack with existing items.
func (p *Player) AddItem(newItem ItemStack) {
	if len(newItem.Data) > 0 {
		// Non-stackable
		p.Inventory = append(p.Inventory, newItem)
		return
	}

	// Try to stack
	for i, item := range p.Inventory {
		if item.ItemID == newItem.ItemID && len(item.Data) == 0 {
			p.Inventory[i].Count += newItem.Count
			return
		}
	}
	p.Inventory = append(p.Inventory, newItem)
}

// RemoveItem removes a count of an item from the inventory.
// For simple items (no Data), it reduces the count or removes the stack.
// For complex items (Data), it removes the first matching item found (by ID).
func (p *Player) RemoveItem(itemID string, count int) bool {
	for i, item := range p.Inventory {
		if item.ItemID == itemID {
			// Check if we can satisfy the count from this stack/item
			if item.Count >= count {
				p.Inventory[i].Count -= count
				if p.Inventory[i].Count == 0 {
					// Remove the item from the slice
					p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
				}
				return true
			}
			return false // Not enough items in this specific stack
		}
	}
	return false // Item not found
}

type PostgresPlayerRepository struct{}

func NewPostgresPlayerRepository() *PostgresPlayerRepository {
	return &PostgresPlayerRepository{}
}

func (r *PostgresPlayerRepository) CreatePlayer(username string, classID string) (*Player, error) {
	// Validate classID exists in static data
	if _, ok := Classes[classID]; !ok {
		return nil, fmt.Errorf("invalid class ID: %s", classID)
	}

	// Default spawn coordinates (The Cradle)
	defaultX := 0.0
	defaultY := 0.0

	// Default empty inventory
	defaultInventory := []ItemStack{}
	inventoryJson, err := json.Marshal(defaultInventory)
	if err != nil {
		return nil, err
	}

	// Default ship layout (tiered slots)
	defaultShip := ShipLayout{
		HighSlots: []Slot{},
		MidSlots:  []Slot{},
		LowSlots:  []Slot{},
		Rigs:      []Slot{},
		Injectors: []Slot{},
	}
	shipJson, err := json.Marshal(defaultShip)
	if err != nil {
		return nil, err
	}

	// Default Ground Gear
	defaultGroundGear := GroundGear{
		Talents: make(map[string]int),
	}
	groundGearJson, err := json.Marshal(defaultGroundGear)
	if err != nil {
		return nil, err
	}

	// Default Skills
	defaultSkills := make(map[string]int)
	skillsJson, err := json.Marshal(defaultSkills)
	if err != nil {
		return nil, err
	}

	// Default Solium
	defaultSolium := 1000 // Starter cash

	// Default Health (Safe Value)
	defaultHealth := 1000.0

	// Default System
	defaultSystem := "Sol-0"
	// Default Motherboard
	defaultMotherboard := Motherboard{
		Slots: 0,
		Chips: []ItemStack{},
		Stats: IndustrialStats{Torque: 0, Compute: 0, Synapse: 0, Flux: 0},
	}
	motherboardJson, err := json.Marshal(defaultMotherboard)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO players (username, class_id, position_x, position_y, inventory, ship_layout, ground_gear, motherboard, solium, skills, current_health, system_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at
	`

	p := &Player{
		Username:      username,
		ClassID:       classID,
		PositionX:     defaultX,
		PositionY:     defaultY,
		Inventory:     defaultInventory,
		Ship:          defaultShip,
		GroundGear:    defaultGroundGear,
		Motherboard:   defaultMotherboard,
		Solium:        defaultSolium,
		Skills:        defaultSkills,
		CurrentHealth: defaultHealth,
		SystemID:      defaultSystem,
	}

	err = db.Pool.QueryRow(context.Background(), query,
		username, classID, defaultX, defaultY,
		string(inventoryJson), string(shipJson), string(groundGearJson), string(motherboardJson),
		defaultSolium, string(skillsJson), defaultHealth, defaultSystem,
	).Scan(&p.ID, &p.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	return p, nil
}

func (r *PostgresPlayerRepository) LoadPlayer(username string) (*Player, error) {
	query := `
		SELECT id, username, class_id, position_x, position_y, inventory, ship_layout, ground_gear, motherboard, solium, skills, current_health, system_id, created_at
		FROM players
		WHERE username = $1
	`

	p := &Player{}
	var inventoryBytes []byte
	var shipBytes []byte
	var groundGearBytes []byte
	var motherboardBytes []byte
	var skillsBytes []byte

	err := db.Pool.QueryRow(context.Background(), query, username).Scan(
		&p.ID,
		&p.Username,
		&p.ClassID,
		&p.PositionX,
		&p.PositionY,
		&inventoryBytes,
		&shipBytes,
		&groundGearBytes,
		&motherboardBytes,
		&p.Solium,
		&skillsBytes,
		&p.CurrentHealth,
		&p.SystemID,
		&p.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to load player: %w", err)
	}

	if len(inventoryBytes) > 0 {
		if err := json.Unmarshal(inventoryBytes, &p.Inventory); err != nil {
			return nil, fmt.Errorf("failed to unmarshal inventory: %w", err)
		}
	} else {
		p.Inventory = []ItemStack{}
	}

	if len(shipBytes) > 0 {
		if err := json.Unmarshal(shipBytes, &p.Ship); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ship layout: %w", err)
		}
	}

	if len(groundGearBytes) > 0 {
		if err := json.Unmarshal(groundGearBytes, &p.GroundGear); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ground gear: %w", err)
		}
	} else {
		p.GroundGear = GroundGear{Talents: make(map[string]int)}
	}

	if len(motherboardBytes) > 0 {
		if err := json.Unmarshal(motherboardBytes, &p.Motherboard); err != nil {
			return nil, fmt.Errorf("failed to unmarshal motherboard: %w", err)
		}
	} else {
		p.Motherboard = Motherboard{
			Slots: 0,
			Chips: []ItemStack{},
			Stats: IndustrialStats{Torque: 0, Compute: 0, Synapse: 0, Flux: 0},
		}
	}

	if len(skillsBytes) > 0 {
		if err := json.Unmarshal(skillsBytes, &p.Skills); err != nil {
			return nil, fmt.Errorf("failed to unmarshal skills: %w", err)
		}
	} else {
		p.Skills = make(map[string]int)
	}

	return p, nil
}

func (r *PostgresPlayerRepository) SavePlayerState(player *Player) error {
	shipJson, err := json.Marshal(player.Ship)
	if err != nil {
		return err
	}

	inventoryJson, err := json.Marshal(player.Inventory)
	if err != nil {
		return err
	}

	groundGearJson, err := json.Marshal(player.GroundGear)
	if err != nil {
		return err
	}

	skillsJson, err := json.Marshal(player.Skills)
	if err != nil {
		return err
	}

	motherboardJson, err := json.Marshal(player.Motherboard)
	if err != nil {
		return err
	}

	query := `
		UPDATE players
		SET position_x = $1, position_y = $2, inventory = $3, ship_layout = $4, ground_gear = $5, motherboard = $6, solium = $7, skills = $8, current_health = $9, system_id = $10
		WHERE id = $11
	`

	_, err = db.Pool.Exec(context.Background(), query,
		player.PositionX,
		player.PositionY,
		string(inventoryJson),
		string(shipJson),
		string(groundGearJson),
		string(motherboardJson),
		player.Solium,
		string(skillsJson),
		player.CurrentHealth,
		player.SystemID,
		player.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to save player state: %w", err)
	}

	return nil
}
