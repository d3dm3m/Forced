package protocol

import "encoding/json"

type Packet struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type MapDataPayload struct {
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	BiomeID string `json:"biome_id"`
	Data    []int  `json:"data"` // Flattened array of Tile IDs
}

type ShipStatsPayload struct {
	CurrentHealth float64 `json:"current_health"`
	MaxHealth     float64 `json:"max_health"`
	BioLoad       int     `json:"bio_load"`
	BioCapacity   int     `json:"bio_capacity"`
	Speed         float64 `json:"speed"`
}

type BuyItemPayload struct {
	ItemID string `json:"item_id"`
}

type LoginSuccessPayload struct {
	Message    string      `json:"message"`
	GroundGear interface{} `json:"ground_gear"`
}

type GraftOrganPayload struct {
	InventoryIndex int    `json:"inventory_index"`
	SlotType       string `json:"slot_type"`
	SlotIndex      int    `json:"slot_index"`
}

// Ground Payloads
type GroundMovementPayload struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"` // Z in 3D
	VelocityX float64 `json:"vx"`
	VelocityY float64 `json:"vy"` // Vz in 3D
	Timestamp int64   `json:"ts"`
}

type GroundEntityState struct {
	ID string  `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}

type GroundStatePayload struct {
	Entities []GroundEntityState `json:"entities"`
}

type GroundAttackPayload struct {
	TargetID string `json:"target_id"`
}

type SpaceAttackPayload struct {
	TargetID string `json:"target_id"`
	WeaponID string `json:"weapon_id"` // E.g., "missile_launcher"
}

type CombatEventPayload struct {
	AttackerID string  `json:"attacker_id"`
	TargetID   string  `json:"target_id"`
	Damage     float64 `json:"damage"`
	IsCrit     bool    `json:"is_crit"`
	IsBeam     bool    `json:"is_beam"` // New for beams
}

// Advanced Ballistics Payloads
type ProjectileSpawnPayload struct {
	ProjectileID string  `json:"projectile_id"`
	ShooterID    string  `json:"shooter_id"`
	TargetID     string  `json:"target_id"`
	Behavior     string  `json:"behavior"` // "missile", "linear", "instant"
	Outcome      string  `json:"outcome"`  // Hit, Miss
	StartX       float64 `json:"start_x"`
	StartY       float64 `json:"start_y"`
	StartZ       float64 `json:"start_z"` // New
	EndX         float64 `json:"end_x"`
	EndY         float64 `json:"end_y"`
	EndZ         float64 `json:"end_z"`   // New
	Speed        float64 `json:"speed"`
}

type CombatHitPayload struct {
	ProjectileID string  `json:"projectile_id"`
	TargetID     string  `json:"target_id"`
	Damage       float64 `json:"damage"`
}

const (
	PACKET_TYPE_SHIP_STATS         = "PACKET_TYPE_SHIP_STATS"
	PACKET_TYPE_BUY_ITEM           = "PACKET_TYPE_BUY_ITEM"
	PACKET_TYPE_GRAFT_ORGAN        = "PACKET_TYPE_GRAFT_ORGAN"
	PACKET_TYPE_GROUND_MOVEMENT    = "PACKET_TYPE_GROUND_MOVEMENT"
	PACKET_TYPE_GROUND_STATE       = "PACKET_TYPE_GROUND_STATE"
	PACKET_TYPE_GROUND_ATTACK      = "PACKET_TYPE_GROUND_ATTACK"
	PACKET_TYPE_FIRE_WEAPON        = "PACKET_TYPE_FIRE_WEAPON" // Shared Space/Ground logical fire? Or specific?
	PACKET_TYPE_COMBAT_EVENT       = "PACKET_TYPE_COMBAT_EVENT"
	PACKET_TYPE_PROJECTILE_SPAWN   = "PACKET_TYPE_PROJECTILE_SPAWN"
	PACKET_TYPE_COMBAT_HIT         = "PACKET_TYPE_COMBAT_HIT"

	// Behaviors
	BEHAVIOR_MISSILE = "missile"
	BEHAVIOR_LINEAR  = "linear"
	BEHAVIOR_INSTANT = "instant"
)
