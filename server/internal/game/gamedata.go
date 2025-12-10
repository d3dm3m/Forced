package game

import (
	"encoding/json"
	"log"
	"os"
)

// StatBlock represents the base statistics for a class or entity
type StatBlock struct {
	Health      int     `json:"health"`
	Speed       int     `json:"speed"` // Ground speed or Space agility
	Stamina     int     `json:"stamina"` // Used for ground actions
	Defense     int     `json:"defense"`
	SensorRange float64 `json:"sensor_range"` // For Space mainly
}

// ClassDefinition mirrors the Class Definitions in Game_Data.md
type ClassDefinition struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // "Ground" or "Space"
	Description string    `json:"description"`
	Archetype   string    `json:"archetype"` // e.g., "The Wall", "The Needle"
	Stats       StatBlock `json:"stats"`
	Slots       int       `json:"slots"` // Number of equipment slots/hardpoints
}

// ShipDefinition for ships.json
type ShipDefinition struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Faction     string                 `json:"faction"`
	Class       string                 `json:"class"`
	Stats       map[string]interface{} `json:"stats"`
	Description string                 `json:"description"`
}

// ItemDefinition mirrors the Item Definitions
type ItemDefinition struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Type      string             `json:"type"`      // e.g., "Weapon", "Module", "Resource"
	Mass      float64            `json:"mass"`
	Price     int                `json:"price"`     // Base price in Solium
	SlotType  string             `json:"slot_type"` // e.g., "high_slots", "cranial", "primary_weapon"
	IsOrganic bool               `json:"is_organic"`
	Stats     map[string]float64 `json:"stats"`    // Base stats modifiers (e.g., "speed": 10.0)
	BioCost   int                `json:"bio_cost"` // Biological strain on the hull
}

var (
	Classes = make(map[string]ClassDefinition)
	Ships   = make(map[string]ShipDefinition)
	Items   = make(map[string]ItemDefinition)
)

// LoadGameData populates the static data structures
func LoadGameData() {
	loadClasses()
	loadShips()
	loadItems()
}

func loadItems() {
	// Sample Items (Hardcoded for prototype)

	// Ground Weapons
	Items["auto_shotgun"] = ItemDefinition{
		ID: "auto_shotgun", Name: "Auto-Shotgun", Type: "Weapon", Mass: 5.0,
		Price: 500, SlotType: "primary_weapon", IsOrganic: false,
		Stats: map[string]float64{"damage": 20.0}, BioCost: 0,
	}
	Items["pistol"] = ItemDefinition{
		ID: "pistol", Name: "Service Pistol", Type: "Weapon", Mass: 1.5,
		Price: 200, SlotType: "sidearm", IsOrganic: false,
		Stats: map[string]float64{"damage": 10.0}, BioCost: 0,
	}

	// Space Modules
	Items["mining_laser"] = ItemDefinition{
		ID: "mining_laser", Name: "Mining Laser", Type: "Module", Mass: 2.0,
		Price: 1000, SlotType: "high_slots", IsOrganic: false,
		Stats: map[string]float64{"mining_yield": 5.0}, BioCost: 0,
	}
	Items["laser_cannon"] = ItemDefinition{
		ID: "laser_cannon", Name: "Laser Cannon", Type: "Module", Mass: 3.0,
		Price: 1200, SlotType: "high_slots", IsOrganic: false,
		Stats: map[string]float64{"damage": 50.0}, BioCost: 0,
	}

	// Bio-Grafts (Organic)
	Items["synthetic_heart"] = ItemDefinition{
		ID: "synthetic_heart", Name: "Synthetic Heart", Type: "Organ", Mass: 0.5,
		Price: 5000, SlotType: "thoracic", IsOrganic: true,
		Stats: map[string]float64{"speed": 20.0, "stamina_regen": 5.0}, BioCost: 10,
	}
	Items["ocular_implant"] = ItemDefinition{
		ID: "ocular_implant", Name: "Ocular Implant", Type: "Organ", Mass: 0.1,
		Price: 2500, SlotType: "cranial", IsOrganic: true,
		Stats: map[string]float64{"sensor_range": 50.0}, BioCost: 5,
	}

	// "Void Heart" for testing Overload/High Stats
	Items["void_heart"] = ItemDefinition{
		ID: "void_heart", Name: "Void Heart", Type: "Organ", Mass: 1.0,
		Price: 15000, SlotType: "thoracic", IsOrganic: true,
		Stats: map[string]float64{"speed": 100.0}, BioCost: 100, // High cost!
	}
}

func loadShips() {
	// Adjust path as needed. Assuming running from server/ root or close to it.
	// In production this might be an absolute path or relative to the executable.
	// Trying relative path "assets/data/ships.json"
	path := "assets/data/ships.json"

	// Check if we are running from cmd/ground or cmd/space
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Try going up levels if running from cmd subdirectories
		path = "../../assets/data/ships.json"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading ships.json: %v", err)
		return
	}

	var shipList []ShipDefinition
	if err := json.Unmarshal(data, &shipList); err != nil {
		log.Printf("Error unmarshalling ships.json: %v", err)
		return
	}

	for _, s := range shipList {
		Ships[s.ID] = s
		log.Printf("Loaded Ship: %s (%s)", s.Name, s.ID)
	}
}

func loadClasses() {
	// Adjust path as needed. Assuming running from server/ root or close to it.
	// In production this might be an absolute path or relative to the executable.
	// Trying relative path "assets/data/classes.json"
	path := "assets/data/classes.json"

	// Check if we are running from cmd/ground or cmd/space
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Try going up levels if running from cmd subdirectories
		path = "../../assets/data/classes.json"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading classes.json: %v", err)
		return
	}

	var classList []ClassDefinition
	if err := json.Unmarshal(data, &classList); err != nil {
		log.Printf("Error unmarshalling classes.json: %v", err)
		return
	}

	for _, c := range classList {
		Classes[c.ID] = c
		log.Printf("Loaded Class: %s (%s)", c.Name, c.ID)
	}
}
