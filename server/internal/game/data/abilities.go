package data

import (
	"time"
)

// AbilityType defines the classification of an ability
type AbilityType string

const (
	AbilityTypeBuff   AbilityType = "Buff"
	AbilityTypeAttack AbilityType = "Attack"
	AbilityTypeUtility AbilityType = "Utility"
)

// AbilityDefinition defines the static data for a class ability
type AbilityDefinition struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	ClassID     string        `json:"class_id"` // Restricted to specific class
	Type        AbilityType   `json:"type"`
	Description string        `json:"description"`
	Cooldown    time.Duration `json:"cooldown"`
	Duration    time.Duration `json:"duration"` // 0 if instant
	Range       float64       `json:"range"`    // Max cast range
	StaminaCost int           `json:"stamina_cost"`
	
	// Modifiers (Simplistic for prototype)
	StatsModifier map[string]float64 `json:"stats_modifier,omitempty"` 
	// e.g. {"defense": 50.0, "speed": -0.5} (Multipliers or Additive)
}

var Abilities = make(map[string]AbilityDefinition)

func LoadAbilities() {
	// Breacher: Siege Mode
	// Immobilizes self to vastly increase defense and firing stability
	Abilities["siege_mode"] = AbilityDefinition{
		ID:          "siege_mode",
		Name:        "Siege Mode",
		ClassID:     "breacher",
		Type:        AbilityTypeBuff,
		Description: "Anchor to the ground, increasing Defense by 200% but reducing Speed to 0.",
		Cooldown:    30 * time.Second,
		Duration:    10 * time.Second,
		StaminaCost: 50,
		StatsModifier: map[string]float64{
			"defense_mult": 2.0, // +200%
			"speed_mult":   0.0, // Rooted
		},
	}

	// Scavenger: Stealth Cloak
	// Becomes invisible to AI and Radar, speed boost
	Abilities["stealth_cloak"] = AbilityDefinition{
		ID:          "stealth_cloak",
		Name:        "Stealth Cloak",
		ClassID:     "scavenger",
		Type:        AbilityTypeUtility,
		Description: "Bend light to become invisible and move faster.",
		Cooldown:    45 * time.Second,
		Duration:    5 * time.Second,
		StaminaCost: 40,
		StatsModifier: map[string]float64{
			"speed_mult": 1.5, // +50% Speed
			// "stealth": 1.0, // Flag
		},
	}
}
