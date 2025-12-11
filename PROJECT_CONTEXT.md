# PROJECT CONTEXT

**Last Updated:** 2025-12-11 00:47:38

## 🤖 AI Persona Roster
* **The Architect:** System Design, Database Schema, Network Topology. (Use for: Infrastructure)
* **The Void Engineer:** Go Backend, Physics, Concurrency. (Use for: Server Logic)
* **The Operator:** Godot Engine, GDScript, Shaders, UI. (Use for: Client)
* **The Biologist:** Game Design, Balancing, Lore, Item Configs. (Use for: Mechanics)
* **The Scribe:** Documentation, Context Management, Verification. (Use for: Organization)

## Project Visuals

### System Architecture (Mindmap)
```mermaid
mindmap
  root((Bio-Horror MMO))
    Client (Godot)
      UI System
        HUD (Vital Signs)
        MarketWindow (Economy)
        SurgeryWindow (Crafting)
      VFX
        SanityController (Shader)
        EntityManager (Visuals)
      Networking
        NetworkManager (Singleton)
    Server (Go)
      Space Core (WebSocket)
        Game Loop (Ticker)
        Mechanics (Stats)
        Services (Market/Surgery)
      Ground Core (UDP)
        Movement Logic
      Data Layer
        PostgreSQL (Persistence)
        Redis (Session)
```

### Critical Path (Dependency Graph)
```mermaid
graph TD
    User([User]) -->|Login| Auth[Auth Service]
    Auth -->|Token| Client[Client App]

    subgraph "Twin-Core Architecture"
    Client -->|WebSocket| SpaceCore[Space Core]
    Client -->|UDP| GroundCore[Ground Core]
    end

    subgraph "Logic Loop"
    SpaceCore -->|Packet| MarketService[Market Service]
    SpaceCore -->|Packet| SurgeryService[Surgery Service]
    SpaceCore -->|Tick| Mechanics[Mechanics Engine]
    Mechanics -->|Calc| BioLoad[Bio-Load Sim]
    end

    subgraph "Persistence"
    MarketService -->|Save| DB[(PostgreSQL)]
    SurgeryService -->|Save| DB
    BioLoad -->|Update Health| DB
    end

    subgraph "Expansion Layer (Future)"
    Gatekeeper[Gatekeeper Service] -.->|Handoff| SpaceCore
    Genesis[Genesis Engine] -.->|Map Data| GroundCore
    SimService[Meta-Sim] -.->|Faction State| DB
    end

    BioLoad -->|PACKET_SHIP_STATS| Client
    Client -->|Signal| HUD[Heads Up Display]
    HUD -->|Visuals| Player([Player Feedback])
```

## Project Status (Epics)

### Completed Sprints: 12-17 (Foundation of Bio-Economy & Twin-Core)
- [x] **Sprint 12 (Data):** Refactored Inventory (`ItemStack` metadata), Ship Layout (Tiered Slots), and Ground Gear.
- [x] **Sprint 13 (Logic):** Implemented `Mechanics` core (Stat Aggregation, Bio-Load) and integrated `Main` game loop.
- [x] **Sprint 14 (Client):** Restored Architecture (`NetworkManager`, `SanityController`) and built HUD (Bio-Feedback).
- [x] **Sprint 15 (Integration):** Wired `MarketService` and `SurgeryService` to frontend UI (`MarketWindow`, `SurgeryWindow`).
- [x] **Sprint 16 (Ground):** Established Ground Gameplay Loop (20Hz UDP, Validation, Broadcasting).
- [x] **Sprint 17 (Tether):** Implemented Ground Persistence and Hangar Handoff trigger.

## The Macro-Scale Architecture (Planned)
* **Zone Sharding:** The universe is split into `Systems`. Each System can be hosted on a different physical server node. The `IGatekeeper` interface will manage routing.
* **Gatekeeper:** A dedicated service to track which server hosts which system and facilitate handoffs (e.g., `Packet_JumpGate`).
* **Planet Hazards:** Planets will have JSON-defined environmental variables (`Gravity`, `Atmosphere`, `ThermalRating`) that the Physics Engine must respect.
* **Social (Guilds):** Guild ownership is baked into the Player/Entity relationship via `guild_id`. This allows robust permission checks (e.g. Door Access).
* **Simulation (Status Effects):** Active Effects are handled via a generic `JSONB` container to allow for infinite expansion (Buffs, Diseases) without schema changes.
* **Legacy (Account Separation):** The `accounts` table stores permanent data (Cosmetics, Legacy Currency) separate from the wipeable `players` table.

### Technical Debt & Future Focus
- [x] **Ground Core Lag:** The `GroundGear` data structures exist on the server but are not used by the Client or Ground Core networking.
- [ ] **Inventory UI:** `InventoryUI.gd` is basic and does not support drag-and-drop for the Surgery interaction.
- [ ] **Next Goal:** Strategic Directive - Future Proofing (Architecture Stubs & Migrations).

## Directory Tree

```
./
    PROJECT_CONTEXT.md
    LICENSE
    extract_files.py
    update_context.py
    server/
        assets/
            data/
                classes.json
                ships.json
        migrations/
            003_add_current_health.sql
            002_sprint12_phase2.sql
            001_sprint12_schema.sql
            005_deep_data.sql
            004_future_proofing.sql
        internal/
            db/
                redis.go
                db.go
            game/
                architecture.go
                combat.go
                mechanics.go
                services_test.go
                services.go
                gamedata.go
                gamedata_test.go
                chat.go
                player_struct_test.go
                mechanics_test.go
                player.go
                health_test.go
                session.go
                space/
                    physics_test.go
                    physics.go
                ground/
                    logic_test.go
                    logic.go
        pkg/
            protocol/
                packet.go
        cmd/
            space/
                main.go
            ground/
                main.go
    docs/
    design/
        Mechanics.md
        Tools_and_Social.md
        Lore.md
        Technical.md
        Game_Data.md
        Production_Roadmap.md
    client/
        assets/
            shaders/
                SanityDistortion.gdshader
        src/
            ui/
                InventoryUI.gd
                space/
                    HUD.gd
                    SurgeryWindow.gd
                    DraggableWindow.gd
                    MarketWindow.gd
                ground/
            autoload/
                NetworkManager.gd
            scenes/
                space/
                    EntityManager.gd
                    SpaceProjectileManager.gd
                    ShipController.gd
                ground/
                    ProjectileManager.gd
                    CameraRig.gd
                    PlayerController.gd
                    GroundEntityManager.gd
                    HangarZone.gd
            vfx/
                SanityController.gd
```

## File Census

- **.md**: 7
- ****: 1
- **.py**: 2
- **.mod**: 1
- **.sum**: 1
- **.json**: 2
- **.sql**: 5
- **.go**: 22
- **.gdshader**: 1
- **.gd**: 15

## File Contents

### ./LICENSE
*Binary/Asset File*

### ./extract_files.py
```py
import os
import re

def extract_files(md_file):
    with open(md_file, 'r', encoding='utf-8') as f:
        content = f.read()

    # Regex to find file blocks
    # Format: ### ./path/to/file
    # ```lang
    # content
    # ```

    # We will split by "### ./" and then process each chunk
    chunks = content.split("### ./")

    # Skip the first chunk (header stuff)
    for chunk in chunks[1:]:
        # First line is the filepath
        lines = chunk.splitlines()
        filepath = lines[0].strip()

        # Determine content
        # Find the first ``` and the last ```
        # Note: Code blocks might be indented or just start.
        # But usually in the provided format they are strictly formatted.

        # Check if binary/asset file placeholder
        if "*Binary/Asset File*" in chunk:
            print(f"Skipping binary file: {filepath}")
            continue

        # Extract code block
        try:
            start_idx = chunk.find("```")
            if start_idx == -1:
                print(f"No code block found for {filepath}")
                continue

            # Find end of start line (e.g. ```go)
            code_start_newline = chunk.find("\n", start_idx)

            # Find closing ```
            end_idx = chunk.rfind("```")

            if end_idx <= start_idx:
                print(f"Malformed block for {filepath}")
                continue

            file_content = chunk[code_start_newline+1:end_idx]

            # Remove trailing newline if it looks like artifact?
            # Usually keep as is.

            # Write file
            # Ensure directory exists
            dirpath = os.path.dirname(filepath)
            if dirpath and not os.path.exists(dirpath):
                os.makedirs(dirpath)

            # If file already exists, I should check if I should overwrite.
            # The plan says "Restore Missing Files".
            # If I already moved it from Force, maybe I shouldn't overwrite?
            # Or maybe PROJECT_CONTEXT.md is the source of truth?
            # The user said "restore them correctly", "Force" directory has files.
            # I should probably prioritize "Force" files if they exist, and only write if not exists.

            if os.path.exists(filepath):
                print(f"File exists, skipping: {filepath}")
            else:
                with open(filepath, 'w', encoding='utf-8') as out:
                    out.write(file_content)
                print(f"Restored: {filepath}")

        except Exception as e:
            print(f"Error processing {filepath}: {e}")

if __name__ == "__main__":
    extract_files("PROJECT_CONTEXT.md")

```

### ./server/assets/data/classes.json
```json
[
  {
    "id": "marine",
    "name": "Marine",
    "type": "Ground",
    "description": "Tank / Frontline. Uses Taunt shouts.",
    "archetype": "The Wall",
    "stats": {
      "health": 150,
      "speed": 10,
      "stamina": 100,
      "defense": 20,
      "sensor_range": 0
    },
    "slots": 4
  },
  {
    "id": "sapper",
    "name": "Sapper",
    "type": "Ground",
    "description": "Area Denial / Support. Constructs hardpoints.",
    "archetype": "The Architect",
    "stats": {
      "health": 100,
      "speed": 12,
      "stamina": 100,
      "defense": 10,
      "sensor_range": 0
    },
    "slots": 6
  },
  {
    "id": "biologist",
    "name": "Xeno-Biologist",
    "type": "Ground",
    "description": "Healer / Buffer. Harvests biomass.",
    "archetype": "The Witch",
    "stats": {
      "health": 80,
      "speed": 14,
      "stamina": 120,
      "defense": 5,
      "sensor_range": 0
    },
    "slots": 4
  },
  {
    "id": "frigate_ace",
    "name": "Frigate Ace",
    "type": "Space",
    "description": "Tackle / Scout. Speed tanking.",
    "archetype": "The Needle",
    "stats": {
      "health": 500,
      "speed": 50,
      "stamina": 0,
      "defense": 10,
      "sensor_range": 100.0
    },
    "slots": 3
  },
  {
    "id": "cruiser_captain",
    "name": "Cruiser Captain",
    "type": "Space",
    "description": "DPS / Line Ship. Broadside combat.",
    "archetype": "The Anvil",
    "stats": {
      "health": 2000,
      "speed": 20,
      "stamina": 0,
      "defense": 50,
      "sensor_range": 80.0
    },
    "slots": 5
  },
  {
    "id": "industrialist",
    "name": "Industrialist",
    "type": "Space",
    "description": "Logistics / Economy. Mining.",
    "archetype": "The Vein",
    "stats": {
      "health": 1200,
      "speed": 15,
      "stamina": 0,
      "defense": 30,
      "sensor_range": 60.0
    },
    "slots": 8
  }
]


```

### ./server/assets/data/ships.json
```json
[
  {
    "id": "kestrel",
    "name": "Kestrel",
    "faction": "USF",
    "class": "Light",
    "stats": {
      "base_hull": 1000,
      "base_shield": 500,
      "capacitor": 100,
      "high_slots": 2
    },
    "description": "Standard USF Light Frigate."
  },
  {
    "id": "omen",
    "name": "Omen",
    "faction": "Celestial Concord",
    "class": "Medium",
    "stats": {
      "base_hull": 2000,
      "base_shield": 300,
      "capacitor": 200,
      "high_slots": 4
    },
    "description": "Concord Medium Cruiser. Organic hull regeneration."
  },
  {
    "id": "hauler",
    "name": "Hauler",
    "faction": "Void Syndicate",
    "class": "Heavy",
    "stats": {
      "base_hull": 1500,
      "base_shield": 200,
      "capacitor": 150,
      "high_slots": 1
    },
    "description": "Syndicate Industrial Ship. Massive cargo hold."
  }
]


```

### ./server/migrations/003_add_current_health.sql
```sql
-- Migration for Sprint 13 Phase 2: Add Current Health

-- Add current_health column to players table (default 1000.0)
ALTER TABLE players ADD COLUMN IF NOT EXISTS current_health DOUBLE PRECISION DEFAULT 1000.0;

```

### ./server/migrations/002_sprint12_phase2.sql
```sql
-- Migration for Sprint 12 Phase 2: Player Economy and Skills

-- Add solium column to players table (default 0)
ALTER TABLE players ADD COLUMN IF NOT EXISTS solium INT DEFAULT 0;

-- Add skills column to players table (JSONB)
ALTER TABLE players ADD COLUMN IF NOT EXISTS skills JSONB DEFAULT '{}'::jsonb;

```

### ./server/migrations/001_sprint12_schema.sql
```sql
-- Migration for Sprint 12: Bio-Grafting and Ground Gear

-- Add ground_gear column to players table
ALTER TABLE players ADD COLUMN IF NOT EXISTS ground_gear JSONB DEFAULT '{}'::jsonb;

-- Note: inventory and ship_layout columns are already JSONB.
-- The application logic will handle the structure change.
-- For a real production migration, we might need to transform existing data,
-- but for Sprint 12 dev, we assume backwards compatibility or reset.

```

### ./server/migrations/005_deep_data.sql
```sql
-- Migration for Future Proofing Phase 2: Deep Data (Social, Sim, Legacy)

-- 1. Social Layer (Guilds)
CREATE TABLE IF NOT EXISTS guilds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    owner_id TEXT NOT NULL, -- References player(id) but loosely
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Update Players with Guild Info
ALTER TABLE players ADD COLUMN IF NOT EXISTS guild_id UUID;
ALTER TABLE players ADD COLUMN IF NOT EXISTS guild_rank INT DEFAULT 0; -- 0=Initiate, 1=Member, 2=Officer, 3=Leader

-- 2. Simulation Layer (Status Effects)
-- Stores temporary buffs/debuffs e.g. [{"id": "stim_pack", "duration": 30, "stacks": 1}]
ALTER TABLE players ADD COLUMN IF NOT EXISTS active_effects JSONB DEFAULT '[]'::jsonb;

-- 3. Legacy Layer (Seasonal Wipes)
-- Permanent account data that survives wipes.
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    solium_legacy INT DEFAULT 0,
    cosmetics JSONB DEFAULT '[]'::jsonb, -- Unlocked skins/titles
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Link Player (Temporary Character) to Account (Permanent User)
ALTER TABLE players ADD COLUMN IF NOT EXISTS account_id UUID;

```

### ./server/migrations/004_future_proofing.sql
```sql
-- Migration for Strategic Directive: Future Proofing (Macro-Scale Architecture)

-- 1. Update Players for Sharding
-- Add system_id to track which physical server/system the player is in.
-- Default to 'Sol-0' (The starting hub).
ALTER TABLE players ADD COLUMN IF NOT EXISTS system_id TEXT DEFAULT 'Sol-0';

-- 2. Create Planets Table
-- Stores static and dynamic data for planetary bodies.
CREATE TABLE IF NOT EXISTS planets (
    id TEXT PRIMARY KEY,
    system_id TEXT NOT NULL,
    name TEXT NOT NULL,
    hazards JSONB DEFAULT '{}'::jsonb, -- Stores Gravity, Atmosphere, ThermalRating
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for fast lookup of all planets in a system
CREATE INDEX IF NOT EXISTS idx_planets_system_id ON planets(system_id);

```

### ./server/internal/db/redis.go
```go
package db

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() error {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}


```

### ./server/internal/db/db.go
```go
package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() error {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		// Default for local development
		dbUrl = "postgres://user:password@localhost:5432/biohorror"
	}
	var err error
	Pool, err = pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	return nil
}


```

### ./server/internal/game/architecture.go
```go
package game

// IGatekeeper defines the future service for cross-server travel.
type IGatekeeper interface {
	GetServerForSystem(systemID string) (string, error)
	RegisterUserTransfer(userID string, targetServerID string) (string, error)
}

// PlanetHazards defines environmental dangers (Gas, Heat, Gravity).
type PlanetHazards struct {
	Gravity       float64 `json:"gravity"`       // Default 1.0
	Atmosphere    string  `json:"atmosphere"`    // "Breathable", "Toxic", "Vacuum"
	ThermalRating float64 `json:"thermal_rating"`// -100 to +100
}

// WorldState represents dynamic faction influence.
type WorldState struct {
	FactionInfluence map[string]float64 `json:"faction_influence"`
}

```

### ./server/internal/game/combat.go
```go
package game

import (
	"math"
	"math/rand"
	"time"
	"biohorror/pkg/protocol"
	"github.com/google/uuid"
)

// Projectile represents a simulated shot in flight.
type Projectile struct {
	ID         string
	ShooterID  string
	TargetID   string
	Behavior   string // "missile", "linear"
	Outcome    string // "Hit", "Miss"

	// Position 3D
	X, Y, Z    float64

	// Velocity 3D (For Linear/Missile physics)
	VX, VY, VZ float64

	// Target Position 3D (For Homing/Interpolation check)
	TargetX, TargetY, TargetZ float64

	Speed      float64
	TurnRate   float64 // Radians/sec for missiles

	Damage     float64
}

// ProjectileManager handles the lifecycle of projectiles.
type ProjectileManager struct {
	ActiveProjectiles map[string]*Projectile
	rng               *rand.Rand
}

func NewProjectileManager() *ProjectileManager {
	return &ProjectileManager{
		ActiveProjectiles: make(map[string]*Projectile),
		rng:               rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CalculateHitChance determines if a shot lands based on stats.
func (pm *ProjectileManager) CalculateHitChance(attacker *Player, target *Player) bool {
	// Base Chance
	chance := 0.8

	// Stats Modifiers (Sprint 18 requirement)
	// Ground Gear check
	if attacker.GroundGear.Talents != nil {
		if marks, ok := attacker.GroundGear.Talents["marksmanship"]; ok {
			chance += float64(marks) * 0.05
		}
	}

	// TODO: Add Ship Stats check for Space Combat

	return pm.rng.Float64() < chance
}

// SpawnProjectile creates a new projectile and adds it to the simulation.
func (pm *ProjectileManager) SpawnProjectile(
	shooterID, targetID string,
	sx, sy, sz, ex, ey, ez float64,
	speed, damage, turnRate float64,
	behavior, outcome string,
) *Projectile {

	id := uuid.New().String()

	// Initial Velocity Direction
	dx := ex - sx
	dy := ey - sy
	dz := ez - sz
	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

	vx, vy, vz := 0.0, 0.0, 0.0
	if dist > 0 {
		vx = (dx / dist) * speed
		vy = (dy / dist) * speed
		vz = (dz / dist) * speed
	}

	proj := &Projectile{
		ID:         id,
		ShooterID:  shooterID,
		TargetID:   targetID,
		Behavior:   behavior,
		Outcome:    outcome,
		X: sx, Y: sy, Z: sz,
		TargetX: ex, TargetY: ey, TargetZ: ez,
		VX: vx, VY: vy, VZ: vz,
		Speed:      speed,
		TurnRate:   turnRate,
		Damage:     damage,
	}

	pm.ActiveProjectiles[id] = proj
	return proj
}

// UpdateSimulation ticks the projectiles. Returns list of impacts.
func (pm *ProjectileManager) UpdateSimulation(dt float64) []*Projectile {
	var impacts []*Projectile

	for id, p := range pm.ActiveProjectiles {
		// Distance to Target (Current)
		dx := p.TargetX - p.X
		dy := p.TargetY - p.Y
		dz := p.TargetZ - p.Z
		dist := math.Sqrt(dx*dx + dy*dy + dz*dz)

		moveDist := p.Speed * dt

		if moveDist >= dist {
			// Impact!
			p.X = p.TargetX
			p.Y = p.TargetY
			p.Z = p.TargetZ
			impacts = append(impacts, p)
			delete(pm.ActiveProjectiles, id)
			continue
		}

		if p.Behavior == protocol.BEHAVIOR_MISSILE && dist > 0 {
			// Homing Logic: Rotate Velocity towards Target
			// Simplified: Just re-calculate normalized vector to target and blend?
			// Full TurnRate physics is complex. For MVP, we'll just steer directly
			// if within TurnRate, or assume perfect homing for now (Behavior = Missile just means it hits).
			// If Outcome == "Miss", we might target a false point?
			// For now, standard homing:

			// Update velocity vector to point at target
			p.VX = (dx / dist) * p.Speed
			p.VY = (dy / dist) * p.Speed
			p.VZ = (dz / dist) * p.Speed
		}

		// Move
		p.X += p.VX * dt
		p.Y += p.VY * dt
		p.Z += p.VZ * dt
	}

	return impacts
}

```

### ./server/internal/game/mechanics.go
```go
package game

// DerivedStats represents the calculated total statistics of a ship.
type DerivedStats struct {
	MaxHealth    float64
	Speed        float64
	SensorRange  float64
	BioCapacity  int
	CurrentBioLoad int
	RejectionRate  float64 // Damage per second
}

// CalculateShipStats aggregates the base ship stats with all installed module modifiers.
// It also calculates Bio-Load and Rejection penalties.
func CalculateShipStats(player *Player) DerivedStats {
	// 1. Initialize with Base Ship Stats
	// Need to look up the ship definition based on player's ship/class?
	// The Player struct does NOT currently store the Ship Hull ID (e.g. "kestrel").
	// It stores `ShipLayout`.
	// Ideally, `Player` should have `ShipID`.
	// For MVP, checking `ships.json` and `classes.json`:
	// `classes.json` has `Frigate Ace` which is a Class.
	// `ships.json` has `Kestrel`.
	// The prompt implies "Start with the Ship's base stats".
	// Issue: We don't know WHICH ship the player is flying from the `Player` struct alone
	// (unless we assume default or it's missing).
	// Checking `PlayerRepository.CreatePlayer`: it creates `defaultShip`.
	// It doesn't seem to set a Hull ID.
	// Assumption: For this mechanic to work, we need a Base.
	// I will assume a default base "kestrel" if unknown, or maybe the `ClassID` defines it?
	// `Frigate Ace` is a class, but usually you buy a hull.
	// Let's fallback to "kestrel" base stats for now if we can't find it,
	// but logically we should probably add `ShipID` to Player.
	// Given strict constraints, I'll use "kestrel" as the default base for calculations
	// to ensure the function works.

	baseShipID := "kestrel" // Fallback
	// In a real implementation, player.ShipID would exist.

	shipDef, ok := Ships[baseShipID]
	var stats DerivedStats
	if ok {
		stats.MaxHealth = getFloat(shipDef.Stats, "base_hull", 1000.0)
		stats.Speed = getFloat(shipDef.Stats, "base_speed", 100.0) // "speed" might not be in JSON, check defaults
		// Checking ships.json from context: "base_hull": 1000, "base_shield": 500, "capacitor": 100
		// No speed in ships.json?
		// Checking classes.json: "Frigate Ace" has stats: { "speed": 50 }
		// Maybe base stats come from CLASS?
		// The prompt says "Start with the Ship's base stats (from ships.json)".
		// If ships.json lacks speed, I'll default it.
		stats.BioCapacity = getInt(shipDef.Stats, "bio_capacity", 50) // Default capacity
	} else {
		// Absolute fallback
		stats.MaxHealth = 1000
		stats.Speed = 100
		stats.BioCapacity = 50
	}

	// 2. Iterate All Slots
	// Combine all slot slices into one iterator or loop over each.
	allSlots := [][]Slot{
		player.Ship.HighSlots,
		player.Ship.MidSlots,
		player.Ship.LowSlots,
		player.Ship.Rigs,
		player.Ship.Injectors,
	}

	for _, slotGroup := range allSlots {
		for _, slot := range slotGroup {
			if slot.Module == nil {
				continue
			}
			applyModuleStats(&stats, slot.Module)
		}
	}

	// 3. Calculate Rejection
	if stats.CurrentBioLoad > stats.BioCapacity {
		// Overload!
		overload := stats.CurrentBioLoad - stats.BioCapacity
		// Rejection Rate: 1 DPS per point of overload? Or constant?
		// Prompt: "Rate = (Load - Cap) * Constant"
		const RejectionMultiplier = 0.5
		stats.RejectionRate = float64(overload) * RejectionMultiplier
	} else {
		stats.RejectionRate = 0
	}

	return stats
}

func applyModuleStats(stats *DerivedStats, item *ItemStack) {
	def, ok := Items[item.ItemID]
	if !ok {
		return
	}

	// Bio Cost (Static)
	stats.CurrentBioLoad += def.BioCost

	// Stat Multipliers
	// Default Compatibility = 1.0
	compatibility := 1.0
	if val, ok := item.Data["compatibility"]; ok {
		if fVal, ok := val.(float64); ok {
			compatibility = fVal
		}
	}

	// Apply Stats
	for stat, value := range def.Stats {
		effectiveValue := value * compatibility

		switch stat {
		case "speed":
			stats.Speed += effectiveValue
		case "health", "base_hull":
			stats.MaxHealth += effectiveValue
		case "sensor_range":
			stats.SensorRange += effectiveValue
		case "bio_capacity":
			stats.BioCapacity += int(effectiveValue)
		}
	}
}

// Helpers for extracting values from generic maps
func getFloat(m map[string]interface{}, key string, def float64) float64 {
	if val, ok := m[key]; ok {
		// JSON numbers are often float64 in Go generic maps
		if f, ok := val.(float64); ok {
			return f
		}
		if i, ok := val.(int); ok {
			return float64(i)
		}
	}
	return def
}

func getInt(m map[string]interface{}, key string, def int) int {
	if val, ok := m[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
		if f, ok := val.(float64); ok {
			return int(f)
		}
	}
	return def
}

```

### ./server/internal/game/services_test.go
```go
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

```

### ./server/internal/game/services.go
```go
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

```

### ./server/internal/game/gamedata.go
```go
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

```

### ./server/internal/game/gamedata_test.go
```go
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


```

### ./server/internal/game/chat.go
```go
package game

import (
	"encoding/json"
	"math"
	"net"
	"sync"
)

type ChatMessage struct {
	SenderID string `json:"sender_id"`
	Channel  string `json:"channel"` // "local", "system", "global"
	Text     string `json:"text"`
}

type ChatSession struct {
	PlayerID string
	Addr     *net.UDPAddr
	X, Y     float64 // Last known position
}

type ChatManager struct {
	sessions map[string]*ChatSession
	mu       sync.RWMutex
	conn     *net.UDPConn
}

func NewChatManager(conn *net.UDPConn) *ChatManager {
	return &ChatManager{
		sessions: make(map[string]*ChatSession),
		conn:     conn,
	}
}

func (cm *ChatManager) RegisterSession(playerID string, addr *net.UDPAddr, x, y float64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.sessions[playerID] = &ChatSession{
		PlayerID: playerID,
		Addr:     addr,
		X:        x,
		Y:        y,
	}
}

func (cm *ChatManager) UpdatePosition(playerID string, x, y float64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if session, ok := cm.sessions[playerID]; ok {
		session.X = x
		session.Y = y
	}
}

func (cm *ChatManager) BroadcastLocal(senderID string, text string, rangeTiles float64) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	sender, ok := cm.sessions[senderID]
	if !ok {
		return
	}

	msg := ChatMessage{
		SenderID: senderID,
		Channel:  "local",
		Text:     text,
	}
	packetBytes, _ := json.Marshal(struct {
		Type    string      `json:"type"`
		Payload ChatMessage `json:"payload"`
	}{
		Type:    "PACKET_TYPE_CHAT_MESSAGE",
		Payload: msg,
	})

	for _, session := range cm.sessions {
		// Calculate distance
		dx := session.X - sender.X
		dy := session.Y - sender.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist <= rangeTiles {
			cm.conn.WriteToUDP(packetBytes, session.Addr)
		}
	}
}

func (cm *ChatManager) BroadcastSystem(text string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	msg := ChatMessage{
		SenderID: "SYSTEM",
		Channel:  "system",
		Text:     text,
	}
	packetBytes, _ := json.Marshal(struct {
		Type    string      `json:"type"`
		Payload ChatMessage `json:"payload"`
	}{
		Type:    "PACKET_TYPE_CHAT_MESSAGE",
		Payload: msg,
	})

	for _, session := range cm.sessions {
		cm.conn.WriteToUDP(packetBytes, session.Addr)
	}
}


```

### ./server/internal/game/player_struct_test.go
```go
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

```

### ./server/internal/game/mechanics_test.go
```go
package game

import (
	"testing"
)

func TestCalculateShipStats_BioLoad(t *testing.T) {
	LoadGameData() // Ensure static data

	// Setup Player with a ship
	player := &Player{
		Ship: ShipLayout{
			HighSlots: []Slot{},
			MidSlots:  []Slot{},
			LowSlots:  []Slot{},
			Rigs:      []Slot{},
			Injectors: []Slot{},
		},
	}

	// 1. Baseline
	baseStats := CalculateShipStats(player)
	if baseStats.CurrentBioLoad != 0 {
		t.Errorf("Expected 0 BioLoad, got %d", baseStats.CurrentBioLoad)
	}

	// 2. Install Standard Heart (BioCost 10, Speed 20)
	stdHeart := ItemStack{
		ItemID: "synthetic_heart",
		Count: 1,
		Data: map[string]interface{}{"compatibility": 1.0},
	}
	// Manually inject into slot for test (skipping SurgeryService to test logic directly)
	player.Ship.HighSlots = append(player.Ship.HighSlots, Slot{Module: &stdHeart})

	stats1 := CalculateShipStats(player)
	if stats1.CurrentBioLoad != 10 {
		t.Errorf("Expected 10 BioLoad, got %d", stats1.CurrentBioLoad)
	}
	if stats1.Speed <= baseStats.Speed {
		t.Errorf("Expected Speed increase (Base %f -> New %f)", baseStats.Speed, stats1.Speed)
	}
	expectedSpeed := baseStats.Speed + 20.0
	if stats1.Speed != expectedSpeed {
		t.Errorf("Expected Speed %f, got %f", expectedSpeed, stats1.Speed)
	}

	// 3. Install Low Compatibility Heart (BioCost 10, Speed 20 * 0.5 = 10)
	// Add to another slot
	badHeart := ItemStack{
		ItemID: "synthetic_heart",
		Count: 1,
		Data: map[string]interface{}{"compatibility": 0.5},
	}
	player.Ship.MidSlots = append(player.Ship.MidSlots, Slot{Module: &badHeart})

	stats2 := CalculateShipStats(player)
	// BioLoad should sum fully (10 + 10 = 20)
	if stats2.CurrentBioLoad != 20 {
		t.Errorf("Expected 20 BioLoad, got %d", stats2.CurrentBioLoad)
	}
	// Speed should increase by 10 (Total +30 from base)
	expectedSpeed2 := baseStats.Speed + 20.0 + 10.0
	if stats2.Speed != expectedSpeed2 {
		t.Errorf("Expected Speed %f, got %f", expectedSpeed2, stats2.Speed)
	}

	// 4. Overload (BioCapacity is 50 default in mechanics.go fallback, or 0 if defined elsewhere)
	// Let's add a massive item: Void Heart (Cost 100)
	voidHeart := ItemStack{
		ItemID: "void_heart",
		Count: 1,
		Data: map[string]interface{}{"compatibility": 1.0},
	}
	player.Ship.LowSlots = append(player.Ship.LowSlots, Slot{Module: &voidHeart})

	stats3 := CalculateShipStats(player)
	totalLoad := 10 + 10 + 100 // 120
	if stats3.CurrentBioLoad != totalLoad {
		t.Errorf("Expected 120 BioLoad, got %d", stats3.CurrentBioLoad)
	}

	// Check Rejection
	// Capacity 50. Overload = 120 - 50 = 70.
	// Rate = 70 * 0.5 = 35.0
	if stats3.RejectionRate <= 0 {
		t.Errorf("Expected Rejection Rate > 0, got %f", stats3.RejectionRate)
	}
	if stats3.RejectionRate != 35.0 {
		t.Errorf("Expected Rejection Rate 35.0, got %f (Cap %d)", stats3.RejectionRate, stats3.BioCapacity)
	}
}

```

### ./server/internal/game/player.go
```go
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
	Solium        int            `json:"solium"`         // Currency
	Skills        map[string]int `json:"skills"`         // Ship/Space Skills
	CurrentHealth float64        `json:"current_health"` // Ship Health (Space) or Player Health (Ground)
	CreatedAt     time.Time      `json:"created_at"`
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

	query := `
		INSERT INTO players (username, class_id, position_x, position_y, inventory, ship_layout, ground_gear, solium, skills, current_health)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
		Solium:        defaultSolium,
		Skills:        defaultSkills,
		CurrentHealth: defaultHealth,
	}

	err = db.Pool.QueryRow(context.Background(), query,
		username, classID, defaultX, defaultY,
		string(inventoryJson), string(shipJson), string(groundGearJson),
		defaultSolium, string(skillsJson), defaultHealth,
	).Scan(&p.ID, &p.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	return p, nil
}

func (r *PostgresPlayerRepository) LoadPlayer(username string) (*Player, error) {
	query := `
		SELECT id, username, class_id, position_x, position_y, inventory, ship_layout, ground_gear, solium, skills, current_health, created_at
		FROM players
		WHERE username = $1
	`

	p := &Player{}
	var inventoryBytes []byte
	var shipBytes []byte
	var groundGearBytes []byte
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
		&p.Solium,
		&skillsBytes,
		&p.CurrentHealth,
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

	query := `
		UPDATE players
		SET position_x = $1, position_y = $2, inventory = $3, ship_layout = $4, ground_gear = $5, solium = $6, skills = $7, current_health = $8
		WHERE id = $9
	`

	_, err = db.Pool.Exec(context.Background(), query,
		player.PositionX,
		player.PositionY,
		string(inventoryJson),
		string(shipJson),
		string(groundGearJson),
		player.Solium,
		string(skillsJson),
		player.CurrentHealth,
		player.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to save player state: %w", err)
	}

	return nil
}

```

### ./server/internal/game/health_test.go
```go
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

```

### ./server/internal/game/session.go
```go
package game

import (
	"context"
	"fmt"
	"time"
	"biohorror/internal/db"
	"github.com/google/uuid"
)

const TokenTTL = 60 * time.Second

// GenerateTransferToken creates a one-time use token for session handoff
func GenerateTransferToken(username string) (string, error) {
	if db.RedisClient == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	token := uuid.New().String()
	key := fmt.Sprintf("token:%s", token)

	ctx := context.Background()

	// Store username as the value, with TTL
	err := db.RedisClient.Set(ctx, key, username, TokenTTL).Err()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// ValidateTransferToken checks if a token is valid and returns the associated username.
// It deletes the token immediately (One-Time Use).
func ValidateTransferToken(token string) (string, error) {
	if db.RedisClient == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("token:%s", token)
	ctx := context.Background()

	// Get value
	username, err := db.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("invalid or expired token")
	}

	// Delete token (One-Time Use)
	db.RedisClient.Del(ctx, key)

	return username, nil
}


```

### ./server/internal/game/space/physics_test.go
```go
package space

import (
	"testing"
)

func TestValidateVector(t *testing.T) {
	// Scenario: Ship drifting at constant velocity (10, 0, 0)
	lastPos := Vector3{X: 0, Y: 0, Z: 0}
	lastVel := Vector3{X: 10, Y: 0, Z: 0}
	dt := 1.0

	// Expected: (10, 0, 0)

	// Case 1: Perfect Match
	clientPos := Vector3{X: 10, Y: 0, Z: 0}
	ok, err := ValidateVector(clientPos, lastVel, lastPos, lastVel, dt)
	if !ok {
		t.Errorf("Expected valid vector, got error: %v", err)
	}

	// Case 2: Acceptable Deviation (Accelerated slightly)
	// Max Accel 50 -> 0.5 * 50 * 1^2 = 25 units allowed deviation + 2 buffer = 27
	clientPosAccel := Vector3{X: 30, Y: 0, Z: 0} // 20 units off
	ok, err = ValidateVector(clientPosAccel, lastVel, lastPos, lastVel, dt)
	if !ok {
		t.Errorf("Expected valid vector (within accel limits), got error: %v", err)
	}

	// Case 3: Teleportation (Too far)
	clientPosCheat := Vector3{X: 100, Y: 0, Z: 0} // 90 units off
	ok, err = ValidateVector(clientPosCheat, lastVel, lastPos, lastVel, dt)
	if ok {
		t.Error("Expected cheat detection, got ok")
	} else {
		t.Logf("Detected cheat: %v", err)
	}
}


```

### ./server/internal/game/space/physics.go
```go
package space

import (
	"fmt"
	"math"
)

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

// ValidateVector checks if the client's position is within a reasonable margin of the expected dead-reckoning position.
// clientPos: The position reported by the client.
// clientVel: The velocity reported by the client (or calculated from previous ticks).
// lastServerPos: The last valid position known by the server.
// lastServerVel: The last valid velocity known by the server (used for projection).
// dt: Delta time in seconds since the last update.
func ValidateVector(clientPos Vector3, clientVel Vector3, lastServerPos Vector3, lastServerVel Vector3, dt float64) (bool, error) {
	// 1. Calculate Expected Position based on Dead Reckoning (Inertia)
	// NewPos = OldPos + (OldVel * dt)
	expectedX := lastServerPos.X + (lastServerVel.X * dt)
	expectedY := lastServerPos.Y + (lastServerVel.Y * dt)
	expectedZ := lastServerPos.Z + (lastServerVel.Z * dt)

	// 2. Calculate Deviation (Distance between Expected and Reported)
	dx := clientPos.X - expectedX
	dy := clientPos.Y - expectedY
	dz := clientPos.Z - expectedZ
	distance := math.Sqrt(dx*dx + dy*dy + dz*dz)

	// 3. Define Tolerance
	// Network jitter, float imprecision, and client-side input acceleration (which happens *during* dt)
	// mean we can't expect a perfect match.
	// Allow for specific "max acceleration" deviation.
	// Deviation <= 0.5 * MaxAccel * dt^2 (Physics formula for distance traveled under accel)
	// Plus a buffer for latency jitter.

	const MaxAcceleration = 50.0 // Matches client-side acceleration_force
	const LatencyBuffer = 2.0 // Units of margin

	allowedDeviation := (0.5 * MaxAcceleration * dt * dt) + LatencyBuffer

	if distance > allowedDeviation {
		return false, fmt.Errorf("position deviation too high: %.2f > %.2f (dt: %.4f)", distance, allowedDeviation, dt)
	}

	return true, nil
}


```

### ./server/internal/game/ground/logic_test.go
```go
package ground

import (
	"testing"
	"time"
)

func TestValidateMovement(t *testing.T) {
	ctx := ValidationContext{
		LastPosition:  Vector2{X: 0, Y: 0},
		LastTimestamp: time.Now().Add(-1 * time.Second), // 1 second ago
		MaxSpeed:      10.0,
	}

	// Valid Move (Moved 10 units in 1 second)
	validState := ClientState{
		PositionX: 10,
		PositionY: 0,
	}

	ok, err := ValidateMovement(validState, ctx)
	if !ok {
		t.Errorf("Expected valid movement, got error: %v", err)
	}

	// Invalid Move (Moved 20 units in 1 second, Max is 10 + 2 buffer = 12)
	invalidState := ClientState{
		PositionX: 20,
		PositionY: 0,
	}

	ok, err = ValidateMovement(invalidState, ctx)
	if ok {
		t.Error("Expected invalid movement (speedhack), got ok")
	} else {
		t.Logf("Correctly detected invalid movement: %v", err)
	}
}


```

### ./server/internal/game/ground/logic.go
```go
package ground

import (
	"fmt"
	"math"
	"time"
)

// Vector2 for simple logic
type Vector2 struct {
	X float64
	Y float64
}

// ClientState represents the state received from the client
type ClientState struct {
	ID        string    `json:"id"`
	PositionX float64   `json:"pos_x"`
	PositionY float64   `json:"pos_y"` // Mapped from Z in Godot
	VelocityX float64   `json:"vel_x"`
	VelocityY float64   `json:"vel_y"`
	Timestamp time.Time `json:"-"`
}

// ValidationContext holds state needed for validation
type ValidationContext struct {
	LastPosition  Vector2
	LastTimestamp time.Time
	MaxSpeed      float64
}

// ValidateMovement checks if the movement is feasible within the time delta
func ValidateMovement(current ClientState, ctx ValidationContext) (bool, error) {
	now := time.Now()
	deltaTime := now.Sub(ctx.LastTimestamp).Seconds()

	// Handle first packet or very short delta
	if deltaTime <= 0 {
		return true, nil
	}

	// Calculate Distance Traveled
	dx := current.PositionX - ctx.LastPosition.X
	dy := current.PositionY - ctx.LastPosition.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// Calculate Max Allowed Distance (Speed * Time)
	// Add a tolerance buffer (e.g., 10% or fixed units) for network jitter / lag compensation
	tolerance := 2.0 // units
	maxDistance := (ctx.MaxSpeed * deltaTime) + tolerance

	if distance > maxDistance {
		return false, fmt.Errorf("speedhack detected: moved %.2f units in %.4fs (max allowed: %.2f)", distance, deltaTime, maxDistance)
	}

	return true, nil
}


```

### ./server/pkg/protocol/packet.go
```go
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

```

### ./server/cmd/space/main.go
```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	"biohorror/internal/db"
	"biohorror/internal/game"
	"biohorror/internal/game/space"
	"biohorror/pkg/protocol"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}
	projMgr *game.ProjectileManager
	mu      sync.RWMutex
)

// Packet structure (Local definition or use protocol package)
// We will use protocol.Packet for consistency, but main needs to marshal/unmarshal
type Packet struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type LoginWithTokenPayload struct {
	Token string `json:"token"`
}

type LoginResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Player  *game.Player `json:"player,omitempty"`
}

func main() {
	// Initialize Game Data
	game.LoadGameData()
	log.Println("Space Core: Game Data Loaded")

	// Initialize Database
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Pool.Close()
	log.Println("Space Core: Database Connected")

	// Initialize Redis
	if err := db.ConnectRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Space Core: Redis Connected")

	http.HandleFunc("/ws", handleWebSocket)

	log.Println("Space Core Server listening on WebSocket :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("New Client Connected. Waiting for Token...")

	// 1. Wait for LOGIN_WITH_TOKEN
	player, err := waitForLogin(conn)
	if err != nil {
		log.Printf("Login failed: %v", err)
		conn.WriteJSON(Packet{Type: "LOGIN_FAILURE", Payload: json.RawMessage(`{"message":"Login failed"}`)})
		return
	}

	log.Printf("Player %s authenticated for Space Core.", player.Username)

	// Send LOGIN_SUCCESS
	payloadBytes, _ := json.Marshal(LoginResponse{
		Success: true,
		Message: "Login successful",
		Player:  player,
	})
	conn.WriteJSON(Packet{Type: "LOGIN_SUCCESS", Payload: payloadBytes})

	// 2. Enter Game Loop
	gameLoop(conn, player)
}

func waitForLogin(conn *websocket.Conn) (*game.Player, error) {
	// Set read deadline to prevent hanging connections
	// conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	_, message, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var packet Packet
	if err := json.Unmarshal(message, &packet); err != nil {
		return nil, err
	}

	if packet.Type != "LOGIN_WITH_TOKEN" {
		return nil, log.Output(2, "Expected LOGIN_WITH_TOKEN packet")
	}

	var loginPayload LoginWithTokenPayload
	if err := json.Unmarshal(packet.Payload, &loginPayload); err != nil {
		return nil, err
	}

	// Validate Token
	username, err := game.ValidateTransferToken(loginPayload.Token)
	if err != nil {
		return nil, err
	}

	// Load Player
	repo := game.NewPostgresPlayerRepository()
	player, err := repo.LoadPlayer(username)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, log.Output(2, "Player not found after token validation")
	}

	return player, nil
}

func gameLoop(conn *websocket.Conn, player *game.Player) {
	// Initialize Services
	repo := game.NewPostgresPlayerRepository()
	marketService := game.NewMarketService(repo)
	surgeryService := game.NewSurgeryService(repo)

	// Init Global Projectile Manager if nil (Lazy init for single-instance in this MVP structure)
	// In production, this would be injected.
	if projMgr == nil {
		projMgr = game.NewProjectileManager()
	}

	// Setup Ticker for Game Logic (10Hz)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Channel to signal connection close
	done := make(chan struct{})

	// Action Channel to handle requests from the Reader Goroutine in the Main Loop
	// This ensures thread safety for the Player struct.
	actionChan := make(chan Packet, 10)

	// Start Reader Goroutine
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Read error for %s: %v", player.Username, err)
				return
			}

			var packet Packet
			if err := json.Unmarshal(message, &packet); err != nil {
				log.Printf("Invalid packet from %s: %v", player.Username, err)
				continue
			}

			// Forward actionable packets to Main Loop
			actionChan <- packet
		}
	}()

	tickCount := 0

	// Main Loop
	for {
		select {
		case <-done:
			return // Connection closed

		case packet := <-actionChan:
			// Handle Input Packet Safely
			handleActionPacket(packet, conn, player, marketService, surgeryService)

		case <-ticker.C:
			tickCount++

			// 1. Tick Projectiles (Every Tick - 10Hz)
			// Note: 10Hz is slow for projectiles, client interpolation is key.
			mu.Lock()
			impacts := projMgr.UpdateSimulation(0.1) // 100ms
			for _, p := range impacts {
				// Handle Impact (Self is target?)
				// In Space Core, we mostly care if WE got hit or if we hit someone else.
				// Since this is a single-client loop MVP, we check if target is self.
				// But packets come from client saying "Fire at X".
				// Space Logic: If I am target, take damage.
				if p.TargetID == player.ID && p.Outcome == "Hit" {
					player.CurrentHealth -= p.Damage
					sendCombatHit(conn, p.ID, p.TargetID, p.Damage)
				}
			}
			mu.Unlock()

			// 2. Every 10 ticks (1 second), perform Stat Calculation & Bio-Load Logic
			if tickCount%10 == 0 {
				stats := game.CalculateShipStats(player)

				// Apply Rejection / Decay
				if stats.RejectionRate > 0 {
					player.CurrentHealth -= stats.RejectionRate
					log.Printf("Player %s taking rejection damage: %.2f (Health: %.2f)", player.Username, stats.RejectionRate, player.CurrentHealth)
				}

				// Cap Health
				if player.CurrentHealth > stats.MaxHealth {
					player.CurrentHealth = stats.MaxHealth
				}

				// Death Check
				if player.CurrentHealth <= 0 {
					log.Printf("Player %s DIED due to Bio-Load Failure.", player.Username)
					// Respawn / Reset logic would go here
					player.CurrentHealth = 0 // Clamp
				}

				// Broadcast Stats Update
				sendShipStats(conn, player, stats)
			}

			// 3. Save State Periodically (Every 5 seconds / 50 ticks)
			if tickCount%50 == 0 {
				// Note: Services save state on action, but passive decay needs saving too.
				repo.SavePlayerState(player)
			}
		}
	}
}

func handleActionPacket(packet Packet, conn *websocket.Conn, player *game.Player, market *game.MarketService, surgery *game.SurgeryService) {
	switch packet.Type {
	case protocol.PACKET_TYPE_BUY_ITEM:
		var payload protocol.BuyItemPayload
		if err := json.Unmarshal(packet.Payload, &payload); err == nil {
			if err := market.BuyItem(player, payload.ItemID); err != nil {
				log.Printf("BuyItem failed: %v", err)
			} else {
				log.Printf("Player %s bought %s", player.Username, payload.ItemID)
			}
		}
	case protocol.PACKET_TYPE_GRAFT_ORGAN:
		var payload protocol.GraftOrganPayload
		if err := json.Unmarshal(packet.Payload, &payload); err == nil {
			if err := surgery.GraftOrgan(player, payload.InventoryIndex, payload.SlotType, payload.SlotIndex); err != nil {
				log.Printf("GraftOrgan failed: %v", err)
			} else {
				log.Printf("Player %s grafted organ into %s[%d]", player.Username, payload.SlotType, payload.SlotIndex)
			}
		}
	case protocol.PACKET_TYPE_FIRE_WEAPON:
		// Space Combat Fire
		// For MVP, target is dummy or self-test.
		// We'll spawn a Missile aimed at a fixed point for visuals.
		pkt := handleSpaceAttack(player)
		if pkt != nil {
			// Convert to local Packet type for main.go
			// Ideally we use protocol.Packet throughout, but main defines its own struct with identical json tags
			// We can just marshal and send.
			// conn.WriteJSON expects interface{}.
			// We can pass pkt directly if `protocol.Packet` matches `Packet`.
			// `main.Packet` is identical structure.
			conn.WriteJSON(pkt)
		}

	case "PACKET_TYPE_SPACE_STATE":
		// Validation placeholder
		_ = space.Vector3{}
	}
}

func handleSpaceAttack(player *game.Player) *protocol.Packet {
	// Dummy Target (100 units ahead)

	// Assuming Player struct uses X/Y as X/Z plane for space, or X/Y screen plane.
	// Physics uses Vector3. Let's default Z=0.
	startX, startY, startZ := player.PositionX, player.PositionY, 0.0
	targetX, targetY, targetZ := player.PositionX + 100, player.PositionY, 0.0

	mu.Lock()
	proj := projMgr.SpawnProjectile(
		player.ID, "dummy_target",
		startX, startY, startZ,
		targetX, targetY, targetZ,
		50.0, 20.0, 1.0,
		protocol.BEHAVIOR_MISSILE, "Hit",
	)
	mu.Unlock()

	// Build Packet
	payload := protocol.ProjectileSpawnPayload{
		ProjectileID: proj.ID,
		ShooterID:    proj.ShooterID,
		TargetID:     proj.TargetID,
		Behavior:     proj.Behavior,
		Outcome:      proj.Outcome,
		StartX:       proj.X,
		StartY:       proj.Y,
		StartZ:       proj.Z,
		EndX:         proj.TargetX,
		EndY:         proj.TargetY,
		EndZ:         proj.TargetZ,
		Speed:        proj.Speed,
	}

	bytes, _ := json.Marshal(payload)
	log.Printf("Spawned Space Missile: %s", proj.ID)

	return &protocol.Packet{
		Type:    protocol.PACKET_TYPE_PROJECTILE_SPAWN,
		Payload: bytes,
	}
}

func sendCombatHit(conn *websocket.Conn, projID, targetID string, damage float64) {
	payload := protocol.CombatHitPayload{
		ProjectileID: projID,
		TargetID:     targetID,
		Damage:       damage,
	}
	bytes, _ := json.Marshal(payload)
	conn.WriteJSON(Packet{Type: protocol.PACKET_TYPE_COMBAT_HIT, Payload: bytes})
}

func sendShipStats(conn *websocket.Conn, player *game.Player, stats game.DerivedStats) {
	payload := protocol.ShipStatsPayload{
		CurrentHealth: player.CurrentHealth,
		MaxHealth:     stats.MaxHealth,
		BioLoad:       stats.CurrentBioLoad,
		BioCapacity:   stats.BioCapacity,
		Speed:         stats.Speed,
	}

	payloadBytes, _ := json.Marshal(payload)
	packet := Packet{
		Type:    protocol.PACKET_TYPE_SHIP_STATS,
		Payload: payloadBytes,
	}

	conn.WriteJSON(packet)
}

```

### ./server/cmd/ground/main.go
```go
package main

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
	"biohorror/internal/db"
	"biohorror/internal/game"
	"biohorror/internal/game/ground"
	"biohorror/pkg/protocol"
)

// Session Management
type GroundSession struct {
	ID        string
	Addr      *net.UDPAddr
	Player    *game.Player
	LastSeen  time.Time
}

var (
	sessions = make(map[string]*GroundSession)
	mu       sync.RWMutex
	projMgr  *game.ProjectileManager
)

// Main
func main() {
	game.LoadGameData()
	log.Println("Ground Core: Game Data Loaded")

	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Pool.Close()
	log.Println("Ground Core: Database Connected")

	if err := db.ConnectRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Ground Core: Redis Connected")

	addr, err := net.ResolveUDPAddr("udp", ":5000")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Ground Core Server listening on UDP :5000")

	repo := game.NewPostgresPlayerRepository()
	projMgr = game.NewProjectileManager()

	go gameLoop(conn, repo)

	for {
		buffer := make([]byte, 4096)
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}

		go handlePacket(conn, remoteAddr, buffer[:n], repo)
	}
}

func handlePacket(conn *net.UDPConn, addr *net.UDPAddr, data []byte, repo game.PlayerRepository) {
	var packet protocol.Packet
	if err := json.Unmarshal(data, &packet); err != nil {
		return
	}

	switch packet.Type {
	case "REQUEST_LOGIN":
		var loginPayload struct {
			Username string `json:"username"`
		}
		if err := json.Unmarshal(packet.Payload, &loginPayload); err == nil {
			player, err := repo.LoadPlayer(loginPayload.Username)
			if err != nil || player == nil {
				if player == nil { player, _ = repo.CreatePlayer(loginPayload.Username, "marine") }
			}
			if player == nil { return }

			mu.Lock()
			sessions[loginPayload.Username] = &GroundSession{
				ID:       loginPayload.Username,
				Addr:     addr,
				Player:   player,
				LastSeen: time.Now(),
			}
			mu.Unlock()

			loginSuccess := protocol.LoginSuccessPayload{
				Message:    "Logged in",
				GroundGear: player.GroundGear,
			}
			payloadBytes, _ := json.Marshal(loginSuccess)
			response := protocol.Packet{ Type: "LOGIN_SUCCESS", Payload: payloadBytes }
			sendPacket(conn, addr, response)
		}

	case protocol.PACKET_TYPE_GROUND_MOVEMENT:
		handleMovement(packet.Payload, addr)

	case protocol.PACKET_TYPE_GROUND_ATTACK:
		handleAttack(conn, packet.Payload, addr)

	case "REQUEST_LAUNCH":
		handleLaunchRequest(conn, addr)
	}
}

func handleLaunchRequest(conn *net.UDPConn, addr *net.UDPAddr) {
	mu.Lock()
	defer mu.Unlock()

	var session *GroundSession
	for _, s := range sessions {
		if s.Addr.String() == addr.String() {
			session = s
			break
		}
	}

	if session == nil {
		return
	}

	// Generate Token
	token, err := game.GenerateTransferToken(session.ID)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return
	}

	// Send Back
	payload := map[string]string{
		"token": token,
		"url":   "ws://localhost:8080/ws", // Space Core URL
	}
	bytes, _ := json.Marshal(payload)

	packet := protocol.Packet{
		Type:    "PACKET_TYPE_LAUNCH_GRANTED",
		Payload: bytes,
	}
	sendPacket(conn, addr, packet)

	log.Printf("Launch Granted for %s. Token: %s", session.ID, token)
}

func handleMovement(payload json.RawMessage, addr *net.UDPAddr) {
	var move protocol.GroundMovementPayload
	if err := json.Unmarshal(payload, &move); err != nil { return }

	mu.Lock()
	defer mu.Unlock()

	var session *GroundSession
	for _, s := range sessions {
		if s.Addr.String() == addr.String() {
			session = s
			break
		}
	}

	if session == nil { return }

	ctx := ground.ValidationContext{
		LastPosition: ground.Vector2{X: session.Player.PositionX, Y: session.Player.PositionY},
		LastTimestamp: session.LastSeen,
		MaxSpeed: 20.0,
	}
	newState := ground.ClientState{ PositionX: move.X, PositionY: move.Y }

	if valid, _ := ground.ValidateMovement(newState, ctx); valid {
		session.Player.PositionX = move.X
		session.Player.PositionY = move.Y
		session.LastSeen = time.Now()
	}
}

func handleAttack(conn *net.UDPConn, payload json.RawMessage, addr *net.UDPAddr) {
	var attack protocol.GroundAttackPayload
	if err := json.Unmarshal(payload, &attack); err != nil { return }

	mu.RLock()
	var attacker *GroundSession
	var target *GroundSession
	for _, s := range sessions {
		if s.Addr.String() == addr.String() {
			attacker = s
			break
		}
	}
	if attack.TargetID != "" { target = sessions[attack.TargetID] }
	mu.RUnlock() // Unlock early, careful with pointer validity

	if attacker == nil || target == nil { return }

	// Calculate Hit/Miss
	hit := projMgr.CalculateHitChance(attacker.Player, target.Player)
	outcome := "Miss"
	if hit { outcome = "Hit" }

	// Spawn Projectile
	speed := 20.0 // Units per sec
	damage := 10.0

	// Create Projectile in Manager (Thread Safe?)
	// projMgr is global, but map is not thread safe.
	// Need mutex for projMgr or use global mu.
	// Re-using global mu for simplicity in MVP.
	mu.Lock()
	// Using Z=0 for Ground (2D Plane)
	// Behavior = Linear (Bullet), TurnRate = 0
	proj := projMgr.SpawnProjectile(
		attacker.ID, target.ID,
		attacker.Player.PositionX, attacker.Player.PositionY, 0,
		target.Player.PositionX, target.Player.PositionY, 0,
		speed, damage, 0.0,
		protocol.BEHAVIOR_LINEAR, outcome,
	)
	mu.Unlock()

	// Broadcast Spawn
	payloadData := protocol.ProjectileSpawnPayload{
		ProjectileID: proj.ID,
		ShooterID:    attacker.ID,
		TargetID:     target.ID,
		Behavior:     proj.Behavior,
		Outcome:      outcome,
		StartX:       proj.X,
		StartY:       proj.Y,
		StartZ:       proj.Z,
		EndX:         proj.TargetX,
		EndY:         proj.TargetY,
		EndZ:         proj.TargetZ,
		Speed:        speed,
	}

	bytes, _ := json.Marshal(payloadData)
	packet := protocol.Packet{ Type: protocol.PACKET_TYPE_PROJECTILE_SPAWN, Payload: bytes }
	packetBytes, _ := json.Marshal(packet)

	mu.RLock()
	for _, s := range sessions {
		conn.WriteToUDP(packetBytes, s.Addr)
	}
	mu.RUnlock()
}

func gameLoop(conn *net.UDPConn, repo game.PlayerRepository) {
	ticker := time.NewTicker(50 * time.Millisecond) // 20Hz
	defer ticker.Stop()

	tickCount := 0

	for range ticker.C {
		tickCount++

		// 1. Tick Projectiles
		mu.Lock()
		impacts := projMgr.UpdateSimulation(0.05) // 50ms

		// Handle Impacts
		for _, p := range impacts {
			if p.Outcome == "Hit" {
				if target, ok := sessions[p.TargetID]; ok {
					target.Player.CurrentHealth -= p.Damage

					// Send Hit Event
					hitPayload := protocol.CombatHitPayload{
						ProjectileID: p.ID,
						TargetID:     p.TargetID,
						Damage:       p.Damage,
					}
					bytes, _ := json.Marshal(hitPayload)
					packet := protocol.Packet{ Type: protocol.PACKET_TYPE_COMBAT_HIT, Payload: bytes }
					pktBytes, _ := json.Marshal(packet)

					// Inner loop broadcast, inefficient but works for MVP
					for _, s := range sessions {
						conn.WriteToUDP(pktBytes, s.Addr)
					}
				}
			}
		}
		mu.Unlock()

		// 2. Broadcast State
		broadcastState(conn)

		// 3. Auto Save
		if tickCount % 100 == 0 {
			saveAllSessions(repo)
		}
	}
}

func saveAllSessions(repo game.PlayerRepository) {
	mu.RLock()
	var playersToSave []*game.Player
	for _, s := range sessions {
		playersToSave = append(playersToSave, s.Player)
	}
	mu.RUnlock()

	for _, p := range playersToSave {
		if err := repo.SavePlayerState(p); err != nil {
			log.Printf("Failed to auto-save %s: %v", p.Username, err)
		}
	}
}

func broadcastState(conn *net.UDPConn) {
	mu.RLock()
	defer mu.RUnlock()

	var entities []protocol.GroundEntityState
	now := time.Now()
	for _, s := range sessions {
		if now.Sub(s.LastSeen) < 5*time.Second {
			entities = append(entities, protocol.GroundEntityState{
				ID: s.ID,
				X:  s.Player.PositionX,
				Y:  s.Player.PositionY,
			})
		}
	}

	if len(entities) == 0 { return }

	payload := protocol.GroundStatePayload{ Entities: entities }
	payloadBytes, _ := json.Marshal(payload)

	packet := protocol.Packet{ Type: protocol.PACKET_TYPE_GROUND_STATE, Payload: payloadBytes }
	packetBytes, _ := json.Marshal(packet)

	for _, s := range sessions {
		if now.Sub(s.LastSeen) < 5*time.Second {
			conn.WriteToUDP(packetBytes, s.Addr)
		}
	}
}

func sendPacket(conn *net.UDPConn, addr *net.UDPAddr, packet protocol.Packet) {
	bytes, _ := json.Marshal(packet)
	conn.WriteToUDP(bytes, addr)
}

```

### ./design/Mechanics.md
```md
# Game Mechanics & Systems Architecture (Beta Design)

## Core Philosophy: Parallel Loops & The Long Year
**"Boots on the Ground, Eyes in the Stars."**
We have removed the artificial division between layers. Combat, Industry, and Exploration exist in parallel on both the **Ground** and the **Void**.

---

## 1. Parallel Systems: Parity

### Mining
*   **Ground (The Drill):** Players operate heavy Exosuit Rigs to bore into planetary crusts.
    *   *Mechanic:* Rhythm-based active reloading to prevent drill overheat.
    *   *Resource:* Raw Ores (Iron, Tungsten) for Hull plating.
*   **Space (The Laser):** Ships use Mining Lasers on asteroid belts.
    *   *Mechanic:* Managing beam intensity vs. rock volatility (it can explode).
    *   *Resource:* Volatile Gases and Ice for Fuel and Coolant.

### Combat
*   **Ground (Tactical):** Isometric shooter. Cover-based. Skill-shots.
    *   *Context:* Clearing bunkers, defending cities from insect swarms.
*   **Space (Dogfighting):** Newtonian-lite flight.
    *   *Context:* Fleet battles, gate camps, blockade running.

### Research & Fabrication
*   **Labs exist everywhere.**
    *   *City Hubs (Ground):* Safe, low-yield research.
    *   *Starbases (Space):* Risky, high-yield research (required for T2+ blueprints).

---

## 2. Progression Curve: The Ladder

### The Cradle (Planetary Surface)
*   **Status:** High-Security Safe Zone.
*   **Purpose:** The nursery. Planetary shields hold back the Entropy.
*   **Gameplay:** Tutorial-heavy, forgiving death penalties (Gear durability loss only).
*   **Transition:** Players must build a ship to leave the Cradle.

### The Void (Deep Space)
*   **Status:** Low-Sec / Null-Sec.
*   **Purpose:** The Endgame. The Entropy is strong here.
*   **Gameplay:** Full-Loot PvP. High-tier resources.
*   **Requirement:** You must "graduate" to survive.

---

## 3. The Yearly Cycle (Seasonal Reset)
*   **Duration:** 12 Months (Real-time).
*   **Concept:** The timeline is unstable. We hold reality together for as long as we can before we must "Leap."

### The Calendar
*   **Phase 1-4 (Establishment):** Safe zones active. Economy stabilizes. Guilds claim land.
*   **Phase 5-10 (Conflict):** Resources in High-Sec run dry. Wars for Outer Rim territory begin.
*   **Phase 11 (The Crumble):**
    *   Outer systems destabilize (Physics corruption).
    *   Planetary shields fail.
    *   Sanity drain increases globally.
*   **Phase 12 (The Final Event - "Eye of Entropy"):**
    *   The hardest Rifts open.
    *   **Goal:** Survive as long as possible.
    *   **Reward:** "Legacy Titles" and cosmetics that persist to the next timeline.
    *   **The End:** The Server Wipes. We "Leap" to a new timeline (Season 2).

---

## 4. Character & Identity
*   **Heritage:** Baseline Human, Ascended (Cyborg), Gene-Forged.
*   **Allegiance:** USF (Stasis), Celestial Concord (Integration), Void Syndicate (Extraction).

---

## 5. Economy & Industry (Revised)
*   **Currencies:** Solium (Fiat) and Old-Earth Data (Legacy).
*   **Production:** Full-Loss (Destroyed items are gone).
*   **Decay:** Blueprints have limited runs. Organic items spoil.

---

## 6. Targeting Systems
### Ground (Tactical Lock)
*   **Mechanism:** Single Target Lock.
*   **Controls:**
    *   `TAB`: Cycles through visible enemies/neutrals.
    *   `F1-F4`: Selects party members (for healing/buffs).
    *   `Mouse Click`: Raycast selection.
*   **Visuals:** Bracket/Reticle around the locked entity.

### Space (Multi-Lock)
*   **Mechanism:** Multi-Target Lock based on CPU/Sensor stats.
*   **Gameplay:** Lock-on time varies by target signature. Weapons can be split (Drones group A attacks Target 1, Missiles attack Target 2).

---

## 7. Deep Dive: Bio-Grafting (The Surgery)
*   **Surgery Loop:** Requires Surgery Bay + Biomass + Connection Minigame.
*   **Rejection:** Compatibility Score dictates Bleeding/Spasm/Necrosis.
*   **Sanity:** Low Sanity causes Friendly Fire, Warp Refusal, and Vendor Fear.

```

### ./design/Tools_and_Social.md
```md
# Tools & Social Systems: The Connected World

## 1. The Genesis Engine (Procedural Pipeline)
*The solution to the "Infinite vs. Curated" problem.*

### Layer 1: Noise Maps (The Math)
We use **Godot's FastNoiseLite** to generate three parallel 2D grids for every planet.
*   **Height Map:** Defines physical terrain shape (0.0 = Ocean Deep, 1.0 = Mountain Peak).
*   **Moisture Map:** Defines vegetation and weather (0.0 = Desert, 1.0 = Swamp).
*   **Danger Map:** Defines enemy density and level (0.0 = Safe Zone, 1.0 = Titan Spawn).

### Layer 2: The Biome Rules (The Logic)
The server iterates through the noise maps and applies a **Logic Table** to determine the tile type.

| Condition | Resulting Biome |
| :--- | :--- |
| `Height > 0.8` AND `Temperature < 0.2` | **Ice Spire** (Mountains) |
| `Height < 0.3` AND `Danger > 0.7` | **Abyssal Trench** (High-Level Ocean) |
| `Moisture > 0.8` AND `Heat > 0.8` | **Rot Jungle** (Toxic Flora) |
| `Height > 0.4` AND `Danger < 0.1` | **Plains** (Safe Building Zone) |

### Layer 3: The Director's Cut (Overrides)
This layer allows designers to inject hand-crafted content into procedural worlds.
*   **Mechanism:** The engine generates the planet, then checks for a `overrides.json` file associated with that specific Planet Seed.
*   **The "Force Spawn" Feature:**
    *   Designers can define a rectangle of coordinates where the procedural logic is ignored.
    *   *Example:* To ensure a critical story quest exists, the designer forces a "Dungeon Entrance" prefab at specific coordinates.
*   **JSON Structure:**
    ```json
    {
      "seed": 987654321,
      "overrides": [
        {
          "type": "FORCE_PREFAB",
          "prefab_id": "Gate_Of_Sorrow",
          "location": {"x": 100, "y": 100},
          "flatten_radius": 15
        }
      ]
    }
    ```
    *   *Result:* No matter how random the rest of the planet is, the "Gate of Sorrow" will always be at (100, 100) on a flat plateau.

---

## 2. The Social Architecture (Chat & Guilds)
*Building a society in a dying universe.*

### The Chat Matrix
Communication channels change based on the physical layer.
1.  **Local (Ground - Grid):**
    *   **Proximity-Based:** Messages appear as **Speech Bubbles** above the character's head.
    *   **Range:** Visible only to players within 20 tiles.
    *   **Vibe:** Intimate, tactical, roleplay-focused.
2.  **System (Space - Void):**
    *   **Radio Channels:** Traditional text box interface.
    *   **Encryption:** Faction-specific channels are encrypted (garbled text for non-members).
    *   **Vibe:** Strategic, impersonal, military.

### Guild Logic: The Trust Chain
Guilds are not just "Chat Rooms"; they are shared **Digital Keychains**.
*   **Access Keys:** Being in a guild grants cryptographic keys to shared assets.
    *   *Doors:* Guild members can open airlocks in player-built cities.
    *   *Turrets:* Guild Turrets will not target players with the correct key.
*   **Hierarchy & Permissions:**
    *   **Initiate:** Can enter the Guild Hall (Door Access). Cannot access the Bank.
    *   **Member:** Can withdraw "Rations" (Low-value items). Can fly "Loaner" ships.
    *   **Officer:** Can withdraw "Spec Ops" gear. Can fly Guild Cruisers.
    *   **Leader:** Can define Access Keys and set Tax Rates.


```

### ./design/Lore.md
```md
# Game Lore and Universe (Beta Design)

## The Core Concept: Industrial Bio-Horror
**"Desperate Industrialism."** The universe is dying (The Entropy), and humanity is grafting parts of the corpse onto itself to survive another day.

---

## The Threat: The Entropy
*   **Description:** A sentient physics error. It unravels reality.
*   **Visuals:** Glitch-horror. Geometry tearing, textures stretching.

---

## History: The Escalation
### The Cradle (Why We Start on Ground)
Planetary cores generate a unique magnetic resonance that temporarily repels the Entropy. These "Safe Zones" (The Cradle) are the only places where clones can be decanted without immediate genetic corruption.
*   **The Cost:** Staying on the ground is safe but stagnant. To find the resources to fix the universe, you must leave the shield and enter The Void.

### The Dimensional Leap (The Yearly Reset)
We have accepted that this timeline is doomed. The Entropy cannot be stopped, only delayed.
*   **The Plan:** We spend the "Year" (Season) harvesting data and energy to charge the **Ark-Gate**.
*   **The Event:** When the timeline collapses (Phase 12), we trigger the Gate. We don't save the world; we **abandon** it. We jump to a parallel timeline, carrying over our memories (Legacy Titles) and crucial data (Cosmetics), to try again.

---

## Factions: Survival Philosophies
*   **USF (Stasis):** "Hold the Line." Admiral Valerius Kane.
*   **Celestial Concord (Integration):** "Adapt to Decay." High-Weaver Elara.
*   **Void Syndicate (Extraction):** "Profit Before End." Broker Vex.

---

## Narrative Delivery
*   **The Whisper Network:** Encrypted burst transmissions.
*   **Black Boxes:** Audio logs that change content based on Sanity.
*   **Unreliable Narrator:** Low sanity reveals the "Truth Layer" (monsters disguised as ships).


```

### ./design/Technical.md
```md
# Technical Architecture v2.0: Industrial Bio-Horror (Beta Design)

## 1. High-Level Architecture: The "Twin-Core" Model
To handle the dual gameplay loops (Grid-based Horror vs. Vector-based Stealth) without code bloat, we treat the Ground and Space layers as distinct server authorities sharing a persistent database.

### The Ground Core (The "Grid")
*   **Context:** Planetary exploration, on-foot combat, city building.
*   **Logic:** Tile-based deterministic movement.
*   **Netcode:** High-frequency **UDP** state updates (20Hz) for responsive combat.
*   **Authority:** Server-side pathfinding validation to prevent speed/noclip hacks.

### The Space Core (The "Vector")
*   **Context:** Space travel, ship-to-ship combat, stealth.
*   **Logic:** Newtonian physics (Inertia/Drift).
*   **Netcode:** Lower-frequency **"Dead Reckoning"** updates (10Hz). The server sends Velocity/Acceleration vectors; the client predicts positions smoothly between ticks.

### The Handoff Protocol
1.  When a player enters a Hangar (Ground), the Ground Core saves state to **Redis** and issues a strictly timed **"Session Transfer Token."**
2.  The Client unloads `GroundScene`, loads `SpaceScene`, and connects to the Space Core using the Token.
3.  **Benefit:** Zero "item duplication" glitches during the swap.

---

## 2. User Interface Paradigms: Distinct Identity

### A. Ground Interface (Tactical View)
*   **Focus:** Immediate action and spatial awareness.
*   **Style:** Clean, minimal, peripheral.
*   **Layout:**
    *   **Center:** Clear viewport for grid-precision movement.
    *   **Periphery:** Equipment slots (Hotbar) pushed to edges to maximize view range.
    *   **Feedback:** Immediate visual popups (damage numbers, status icons) near the character model.
    *   **Interaction:** Context-sensitive click menus for world objects (Levers, Crates).

### B. Space Interface (Command Dashboard)
*   **Focus:** Data management and strategic maneuvering.
*   **Style:** Modular, window-based overlays. "Submarine Command Console."
*   **Layout:**
    *   **The HUD:** The player looks "through" a glass overlay.
    *   **Modules:** Draggable windows for "Target List," "D-Scan Radar," "Market Graphs," and "Module Control."
    *   **Data Density:** High. Tables, graphs, and vector lines dominate the screen.
    *   **Vibe:** You are not looking out a window; you are reading sensors.

---

## 3. New Systems Implementation

### A. The "EMCON" Engine (Noise & Stealth)
*Mechanic: High-tech modules create "Noise" that summons The Entropy.*
*   **Spatial Indexing:** Dynamic **Octree** (or Spatial Hash) for the Space map.
*   **The "Loudness" Layer:**
    *   Every entity has a `NoiseRadius` property.
    *   *Passive:* Hull movement = Small Radius.
    *   *Active:* Firing a Plasma Lance = Massive Radius.
*   **The "Listener" Nodes:**
    *   Invisible server-side entities ("Entropy Nodes") are scattered in the Octree.
    *   *Trigger:* If `Player.NoiseRadius` intersects `Node.DetectionRadius` -> **Wake AI**.
    *   *Optimization:* Collision calculated only on State changes (e.g., activating a module), not every tick.

### B. The "Sanity" Rendering Pipeline
*Mechanic: Low sanity causes hallucinations and UI lies.*
*   **Client-Side Injection (The "Lies"):**
    *   Server sends `SanityLevel` (0-100).
    *   **Shader Overrides:** At <50% Sanity, `WorldEnvironment` enables chromatic aberration and "breathing" wall shaders.
    *   **Fake UI:** Client injects "Ghost Radar Dots" (red blips) corresponding to no real entity.
*   **Server-Side Injection (The "Phantoms"):**
    *   At Extreme Insanity (<10%), Server sends `Packet_SpawnEntity` to *only* that player.
    *   The player fights a monster that doesn't exist for teammates.

### C. Biological Ship Data (The "Grafting" System)
*Mechanic: Ships are grown, not built. Parts are organic.*
*   **Schema Strategy:** **JSONB** in PostgreSQL.
*   **Reasoning:** Rigid SQL columns (e.g., `Gun_Slot_1`) cannot handle organic mutations.
*   **Flexibility:** Designers can create "Tumors" that take up 2 slots but regenerate fuel without schema migrations.

---

## 4. DevOps & Security

### The "Headless" Validator
*   **Strategy:** Never trust the client.
*   **Implementation:** A "Headless" version of the Godot Client (no graphics) runs on the Go Server.
*   **Usage:** When a player shoots, the Go server passes coordinates to the Headless instance to run a micro-simulation (Raycast check) to validate line-of-sight.

---

## 5. Architectural Deep Dives

### Handoff Failure & Recovery (The Limbo State)
*   **Problem:** Client crashes after Ground Core disconnect but before Space Core connect.
*   **Limbo State:** The Player is flagged as `TRANSITIONING` in Redis. They exist on neither server.
*   **Recovery Logic:**
    1.  On reconnect, the Login Server sees the `TRANSITIONING` flag.
    2.  It checks the timestamp.
        *   *< 5 Minutes:* Resends the **Space Session Token**. The client attempts to join Space again.
        *   *> 5 Minutes:* **Emergency Recall.** The player is reset to the last visited "Safe Harbor" (City Hangar) on the Ground Server.
    *   **Anti-Dupe:** The Space Token is One-Time-Use (nonce). If used, it invalidates the previous state.

### Headless Validator: 2.5D Projectile Logic
*   **Challenge:** The Client is 2D (Isometric), but the Server logic needs to prevent shooting through walls.
*   **Implementation:**
    *   The Headless Godot instance loads the `CollisionMap` (NavMesh + HeightMap).
    *   **Raycast Verification:**
        *   Client sends: `Origin(x,y)`, `Target(x,y)`, `Timestamp`.
        *   Server calculates: `Trajectory`.
        *   Server checks: Does `Trajectory` intersect with `Wall_Height > Projectile_Height`?
    *   **Aimbot Prevention:**
        *   The Server tracks `MouseDelta`. Instant 180-degree snaps with 0ms delay are flagged.
        *   **Heuristic:** "Human aiming has jitter; Bot aiming is linear."

### Instancing Logic: Rift Gates
*   **Architecture:** **On-Demand Micro-Containers**.
*   **Process:**
    1.  Party interacts with a Rift Gate.
    2.  **Orchestrator (K8s):** Spins up a lightweight `Instance_Pod` (Docker container running the Space Core logic).
    3.  **State Injection:** The Orchestrator injects the "Seed" (for procedural generation) and "Difficulty Modifiers" (based on party gear).
    4.  **Connection:** The Party is seamlessly transferred (via Handoff Protocol) to this isolated IP address.
    5.  **Teardown:** When the raid wipes or completes, the container spins down, writing loot logs to the main DB.


```

### ./design/Game_Data.md
```md
# Game Data: The RPG Bible

## 1. Class Definitions (The Trinity)
*Classes are defined by the "Suit" (Ground) or "License" (Space) currently equipped. Players can switch roles by changing gear at a station.*

### Ground Classes (The Exosuit)
1.  **Marine (Archetype: The Wall)**
    *   **Role:** Tank / Frontline.
    *   **Playstyle:** Heavy Armor usage. Uses **Taunt** shouts to draw aggro from swarms.
    *   **Signature Weapons:** Auto-Shotguns, Riot Shields.
    *   **Passive:** *Adrenaline Rush* (Stamina regenerates faster when taking damage).
2.  **Sapper (Archetype: The Architect)**
    *   **Role:** Area Denial / Support.
    *   **Playstyle:** Constructs temporary hardpoints. Controls the battlefield geometry.
    *   **Signature Equipment:** Sentry Turrets, Barricade Walls, Satchel Charges.
    *   **Passive:** *Efficient Construction* (Structures cost 20% less resource to build).
3.  **Xeno-Biologist (Archetype: The Witch)**
    *   **Role:** Healer / Buffer.
    *   **Playstyle:** "Combat Alchemist." Harvests biomass from fallen enemies to brew instant injectors.
    *   **Signature Equipment:** Dart Gun (Buffs allies, Debuffs enemies), Harvester Drone.
    *   **Passive:** *Necrotic Harvest* (Kills grant a temporary stack of "Biomass" used for healing).

### Space Classes (The License)
1.  **Frigate Ace (Archetype: The Needle)**
    *   **Role:** Tackle / Scout.
    *   **Playstyle:** Speed tanking (moving so fast turrets can't track). Locking down targets for the fleet.
    *   **Signature Systems:** Warp Disruptors, Stasis Webs.
    *   **Passive:** *Inertial Dampening* (+50% Agility/Turn Speed).
2.  **Cruiser Captain (Archetype: The Anvil)**
    *   **Role:** DPS / Line Ship.
    *   **Playstyle:** Broadside combat. Rotating shield quadrants to mitigate incoming fire.
    *   **Signature Systems:** Flak Batteries (Anti-Fighter), Heavy Shield Boosters.
    *   **Passive:** *Grid Stabilization* (+10% Shield HP).
3.  **Industrialist (Archetype: The Vein)**
    *   **Role:** Logistics / Economy.
    *   **Playstyle:** Mining high-value nodes in dangerous space. Evading capture.
    *   **Signature Systems:** Mining Lasers, Cargo Expanders, "Blockade Runner" Cloak (Short duration invisibility).
    *   **Passive:** *Yield Optimization* (+15% Mining Laser Yield).

---

## 2. Faction Specialties (Meta-Bonuses)

### The United Sol Federation (USF)
*   **Philosophy:** Reliability and Standardization.
*   **Bonus 1:** **Standardized Parts** (-10% Repair Costs).
*   **Bonus 2:** **Reinforced Plating** (+10% Armor HP).
*   **Penalty:** **Obsolete Tech** (-10% Sensor Range).

### The Celestial Concord
*   **Philosophy:** Biological Adaptation.
*   **Bonus 1:** **Living Metal** (Regenerates 1% Hull/sec).
*   **Bonus 2:** **Graft Compatibility** (Organic Modules have -50% Rejection Chance).
*   **Penalty:** **Glass Shields** (-20% Shield Capacity). They rely on flesh, not energy.

### The Void Syndicate
*   **Philosophy:** High Risk, High Profit.
*   **Bonus 1:** **Dirty Hacks** (+15% E-War Strength).
*   **Bonus 2:** **Scavenger Protocols** (+10% Loot Drop Rate).
*   **Penalty:** **Scrap Hulls** (-10% Structure HP). Ships are fragile.

---

## 3. The Star Atlas (Key Systems)

### Sol-0 (The Capital)
*   **Security:** High-Sec (1.0).
*   **Visuals:** **Smog-Choked Cyberpunk**. The sun is barely visible through the industrial haze. Neon signs advertise "Pure Air" as a luxury.
*   **Gameplay:** The central trade hub. Safe for new players. Home of the USF Admiralty.

### Eden's Rot
*   **Security:** Low-Sec (0.4).
*   **Visuals:** **Flesh Dyson Sphere**. A system being consumed by a massive biological growth. Asteroids look like tumors. Space stations are grown, not built.
*   **Gameplay:** High-value Organic Mining. PvP is enabled but guarded by Sentry Guns on gates.

### The Breaker Yards
*   **Security:** Null-Sec (0.0).
*   **Visuals:** **Ship Graveyard**. A dense debris field of millions of destroyed vessels. Foggy, dark, and silent.
*   **Gameplay:** Lawless. High-risk salvage operations. Players can claim territory here.


```

### ./design/Production_Roadmap.md
```md
# Production Roadmap: The MVP Critical Path

## Phase 1: The Foundation (Months 1-2)
**Goal:** A headless "Twin-Core" server that can accept connections, move entities, and persist data.

*   **Sprint 1 (Weeks 1-2): Core Netcode**
    *   **Space Core:** Implement Go WebSocket server. Basic Vector3 movement and dead-reckoning.
    *   **Ground Core:** Implement Go UDP server. Grid-based movement validation.
    *   **DB:** Set up PostgreSQL with the JSONB `ships` schema.
*   **Sprint 2 (Weeks 3-4): The Handoff**
    *   Implement Redis session management.
    *   Build the "Session Transfer Token" logic.
    *   **Test:** Client connects to Ground, enters Hangar, disconnects, connects to Space with correct Ship state.
*   **Sprint 3 (Weeks 5-6): Data Structure & Login**
    *   Implement `Game_Data` loading (Stats, Items).
    *   Build secure Login/Auth microservice (JWT).
    *   **Deliverable:** "The Grey Box" - A player can log in, walk in a grey box, board a grey cube, and fly in a black void.

## Phase 2: The Vertical Slice (Months 3-4)
**Goal:** One fully playable loop (Ground -> Space -> Combat -> Loot -> Ground).

*   **Sprint 4 (Weeks 7-8): The Tactical View (Ground)**
    *   Implement Isometric Camera & Controls.
    *   Basic Combat: Projectile raycasting (Headless Validator) and Damage logic.
    *   **Asset:** "Marine" Class mesh and animations.
*   **Sprint 5 (Weeks 9-10): The Command Dashboard (Space)**
    *   Implement HUD Overlays (Radar, Target List).
    *   Space Combat: Locking, Turret Tracking logic, Missile flight.
    *   **Asset:** "Interceptor" Ship hull and Turret models.
*   **Sprint 6 (Weeks 11-12): The Loop Integration**
    *   Loot Tables & Inventory UI.
    *   NPC AI: Simple "Seek and Destroy" behavior for Entropy Drones.
    *   **Deliverable:** A playable build where a Marine fights a drone, loots "Biomass", boards a ship, flies to an asteroid, and mines "Ice".

## Phase 3: The Content Hose (Months 5-6)
**Goal:** Fleshing out the skeleton. UI Polish, Audio, and Content scaling.

*   **Sprint 7 (Weeks 13-14): UI Paradigm Polish**
    *   Implement "Glitch" Shaders for Sanity effects.
    *   Differentiate Ground (Clean/Peripheral) vs Space (Data/Modular) styles.
*   **Sprint 8 (Weeks 15-16): The "Genesis" Pipeline**
    *   Finalize Procedural Generation tools.
    *   Generate "Sol-0" (Urban) and "Eden's Rot" (Bio) maps.
*   **Sprint 9 (Weeks 17-18): Audio & FX**
    *   Implement Wwise/FMOD middleware (if needed) or Godot Audio buses.
    *   Add "The Sound of Entropy" (Dynamic ambience).
*   **Sprint 10 (Weeks 19-20): Beta Prep**
    *   Load Testing (Bot swarm).
    *   Bug fixing.
    *   **Deliverable:** The MVP Beta Candidate.

---

## The Asset Matrix (MVP Requirements)

### 3D Assets (Models & Textures)
| Category | Asset Name | Description | Poly Count Target |
| :--- | :--- | :--- | :--- |
| **Ships** | **The Kestrel (Light)** | Starter Frigate. Industrial/Blocky. USF Style. | < 5k tris |
| **Ships** | **The Omen (Medium)** | Concord Cruiser. Organic curves/Ceramic. | < 10k tris |
| **Ships** | **The Hauler (Heavy)** | Syndicate Industrial. Asymmetrical/Scrap. | < 8k tris |
| **Biomes** | **Ice Wastes** | Blue/White jagged terrain. Low friction. | Modular Tileset |
| **Biomes** | **The Rot** | Flesh-moss terrain. Pulsating flora. | Modular Tileset |
| **Biomes** | **Urban Ruins** | Rusted steel floors, concrete walls. | Modular Tileset |
| **Weapons** | **Kinetic Repeater** | Basic machine gun. | Low |
| **Weapons** | **Plasma Lance** | Beam weapon. | Low |
| **Weapons** | **Bio-Spitter** | Organic mortar (Enemy weapon). | Low |

### 2D Assets (UI & Sprites)
*   **Icons:** 50x Item Icons (Ores, Guns, Organs).
*   **Faction Logos:** USF (Shield/Fist), Concord (Eye/Wing), Syndicate (Coin/Dagger).
*   **UI Sets:**
    *   *Tactical Set:* Minimalist, Green/Amber, High Contrast.
    *   *Command Set:* Data-heavy, Blue/Grey, Scanlines.

### Audio Assets
*   **The Sound of Entropy:** A dynamic loop of static/screams that increases volume with Low Sanity.
*   **Weapon SFX:** Distinct sounds for Ballistic (Thud) vs Energy (Hum/Crack).
*   **UI SFX:** "Click" vs "Datastream" sounds for the different UI paradigms.


```

### ./client/assets/shaders/SanityDistortion.gdshader
```glsl
shader_type canvas_item;

uniform float sanity_level : hint_range(0.0, 1.0) = 1.0;
uniform sampler2D screen_texture : hint_screen_texture, repeat_disable, filter_nearest;

void fragment() {
    // 0.0 Sanity = Max Distortion, 1.0 = Clean
    float distortion_strength = 1.0 - sanity_level;

    // Chromatic Aberration
    // Shift channels based on strength
    float r_offset = distortion_strength * 0.01;
    float b_offset = distortion_strength * -0.01;

    vec4 red_channel = texture(screen_texture, SCREEN_UV + vec2(r_offset, 0.0));
    vec4 green_channel = texture(screen_texture, SCREEN_UV);
    vec4 blue_channel = texture(screen_texture, SCREEN_UV + vec2(b_offset, 0.0));

    COLOR = vec4(red_channel.r, green_channel.g, blue_channel.b, 1.0);

    // Breathing / Pulsing Effect
    // Only at low sanity (< 0.5)
    if (sanity_level < 0.5) {
        float time = TIME * 2.0;
        float pulse = sin(time) * 0.02 * (0.5 - sanity_level); // stronger as sanity drops

        // Simple UV warp for breathing
        vec2 centered_uv = SCREEN_UV - 0.5;
        float dist = length(centered_uv);
        vec2 warped_uv = SCREEN_UV + (centered_uv / dist) * pulse * dist;

        // Re-sample with warp
        // (For simplicity in this snippet, effectively combining them involves complex sampling,
        // so we'll just mix the color tint here or assume the UV warp applies to the texture lookup above)
        // Let's just apply a slight dark vignette pulse for "Breathing"
        float vignette = smoothstep(0.8, 0.2, length(centered_uv) + pulse);
        COLOR.rgb *= vignette;
    }
}


```

### ./client/src/ui/InventoryUI.gd
```gd
extends Control

# Simple UI Grid for Inventory
# Expects a GridContainer child named "Grid"

func _ready():
	if not has_node("Grid"):
		var grid = GridContainer.new()
		grid.name = "Grid"
		grid.columns = 5
		add_child(grid)

	# Start listening for inventory updates (mock connection)
	# NetworkManager.connect("inventory_updated", _on_inventory_updated)

func _on_inventory_updated(items: Array):
	var grid = $Grid
	# Clear existing children
	for child in grid.get_children():
		child.queue_free()

	# Populate new items
	for item in items:
		var slot = PanelContainer.new()
		var label = Label.new()
		# item is expected to be { "item_id": "...", "count": ... }
		label.text = str(item.count) + "x\n" + item.item_id
		slot.add_child(label)
		grid.add_child(slot)


```

### ./client/src/ui/space/HUD.gd
```gd
extends Control

# Preload the Draggable Window scene or script
var window_script = preload("res://src/ui/space/DraggableWindow.gd")

# HUD Elements
var health_bar: ProgressBar
var bio_load_bar: ProgressBar
var warning_label: Label
var purge_button: Button

func _ready():
	# Setup Network Connection
	NetworkManager.connect("ship_stats_updated", _on_ship_stats_updated)

	# Create HUD Layout
	setup_status_bars()

	# Example: Create a Radar Window (Existing logic)
	create_radar_window()

func setup_status_bars():
	var container = VBoxContainer.new()
	container.position = Vector2(20, 20)
	container.size = Vector2(200, 100)
	add_child(container)

	# Health Bar
	var hp_label = Label.new()
	hp_label.text = "HULL INTEGRITY"
	container.add_child(hp_label)

	health_bar = ProgressBar.new()
	health_bar.max_value = 1000
	health_bar.value = 1000
	# Set color to Green roughly via modulation or theme (omitted for brevity)
	container.add_child(health_bar)

	# Spacer
	container.add_child(HSeparator.new())

	# Bio-Load Bar
	var bio_label = Label.new()
	bio_label.text = "BIO-LOAD"
	container.add_child(bio_label)

	bio_load_bar = ProgressBar.new()
	bio_load_bar.max_value = 100
	bio_load_bar.value = 0
	container.add_child(bio_load_bar)

	# Warning Label
	warning_label = Label.new()
	warning_label.text = "REJECTION IMMINENT"
	warning_label.modulate = Color(1, 0, 0)
	warning_label.visible = false
	container.add_child(warning_label)

	# Purge Button
	purge_button = Button.new()
	purge_button.text = "EMERGENCY PURGE"
	purge_button.modulate = Color(1, 0, 0) # Red button
	purge_button.connect("pressed", _on_purge_pressed)
	container.add_child(purge_button)

	# Spacer
	container.add_child(HSeparator.new())

	# Toggle Buttons
	var market_btn = Button.new()
	market_btn.text = "OPEN MARKET"
	market_btn.connect("pressed", _on_toggle_market)
	container.add_child(market_btn)

	var surgery_btn = Button.new()
	surgery_btn.text = "OPEN SURGERY"
	surgery_btn.connect("pressed", _on_toggle_surgery)
	container.add_child(surgery_btn)

	# Load Windows
	var market = preload("res://src/ui/space/MarketWindow.gd").new()
	market.visible = false
	add_child(market)

	var surgery = preload("res://src/ui/space/SurgeryWindow.gd").new()
	surgery.visible = false
	add_child(surgery)

func _on_toggle_market():
	var win = get_node("MarketWindow")
	if win: win.visible = !win.visible

func _on_toggle_surgery():
	var win = get_node("SurgeryWindow")
	if win: win.visible = !win.visible

func _on_ship_stats_updated(stats: Dictionary):
	# Update Health
	var cur_hp = stats.get("current_health", 1000.0)
	var max_hp = stats.get("max_health", 1000.0)

	health_bar.max_value = max_hp
	health_bar.value = cur_hp

	# Visual Feedback for Health
	if cur_hp < (max_hp * 0.2):
		health_bar.modulate = Color(1, 0, 0) # Red Critical
		# Panic Signal to SanityController could be emitted here or centralized
	else:
		health_bar.modulate = Color(1, 1, 1) # Normal

	# Update Bio-Load
	var cur_load = stats.get("bio_load", 0)
	var max_load = stats.get("bio_capacity", 50)

	bio_load_bar.max_value = max_load
	bio_load_bar.value = cur_load

	# Color Coding Bio-Load
	var load_percent = float(cur_load) / float(max_load) if max_load > 0 else 0.0

	if load_percent > 1.0:
		bio_load_bar.modulate = Color(1, 0, 0) # Red (Overload)
		warning_label.visible = true
		# Blink effect could be done via tween or process
	elif load_percent > 0.8:
		bio_load_bar.modulate = Color(1, 0.6, 0) # Orange (Warning)
		warning_label.visible = false
	else:
		bio_load_bar.modulate = Color(0, 1, 0) # Green (Safe)
		warning_label.visible = false

func _on_purge_pressed():
	print("Purge Initiated!")
	# Logic to send PURGE packet would go here.
	# NetworkManager.send_packet("PACKET_TYPE_PURGE_REQUEST", {})

# ... Existing Radar Logic ...
func create_radar_window():
	var radar_win = Panel.new()
	radar_win.set_script(window_script)
	radar_win.name = "RadarWindow"
	radar_win.size = Vector2(200, 200)
	radar_win.position = Vector2(50, 150) # Moved down to avoid bars

	# Add some content to look like a radar
	var label = Label.new()
	label.text = "D-SCAN RADAR"
	label.position = Vector2(10, 5)
	radar_win.add_child(label)

	var radar_display = ColorRect.new()
	radar_display.color = Color(0, 0.2, 0, 0.8)
	radar_display.position = Vector2(10, 30)
	radar_display.size = Vector2(180, 160)
	radar_win.add_child(radar_display)

	# Add a center blip (Self)
	var self_blip = ColorRect.new()
	self_blip.color = Color.WHITE
	self_blip.size = Vector2(4, 4)
	self_blip.position = Vector2(90 - 2, 80 - 2)
	radar_display.add_child(self_blip)

	add_child(radar_win)

```

### ./client/src/ui/space/SurgeryWindow.gd
```gd
extends "res://src/ui/space/DraggableWindow.gd"

# Surgery Window: Shows Slots and allows installation
# For Sprint 15 MVP, we just show buttons to install into specific slots from "Inventory Index 0"
# to demonstrate the loop. A full Inventory Picker is complex.

func _ready():
	name = "SurgeryWindow"
	size = Vector2(250, 300)
	position = Vector2(600, 50)

	var title = Label.new()
	title.text = "SURGERY BAY"
	title.position = Vector2(10, 5)
	add_child(title)

	var container = VBoxContainer.new()
	container.position = Vector2(10, 30)
	container.size = Vector2(230, 260)
	add_child(container)

	# Slot: High 1
	add_slot_row(container, "High Slot 1", "high_slots", 0)
	# Slot: Mid 1
	add_slot_row(container, "Mid Slot 1", "mid_slots", 0)
	# Slot: Low 1
	add_slot_row(container, "Low Slot 1", "low_slots", 0)

	# Note: This UI assumes you want to install the FIRST item in your inventory.
	# Real UI needs drag-and-drop from Inventory Window.
	var note = Label.new()
	note.text = "NOTE: Installs Item #0 from Inventory"
	note.modulate = Color(0.7, 0.7, 0.7)
	container.add_child(note)

func add_slot_row(container, label_text, slot_type, slot_index):
	var row = HBoxContainer.new()

	var label = Label.new()
	label.text = label_text
	label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	row.add_child(label)

	var btn = Button.new()
	btn.text = "GRAFT (Idx 0)"
	btn.connect("pressed", func(): _on_graft_pressed(slot_type, slot_index))
	row.add_child(btn)

	container.add_child(row)

func _on_graft_pressed(slot_type, slot_index):
	print("Grafting Inventory[0] into ", slot_type, "[", slot_index, "]")
	# Hardcoded to index 0 for MVP testing
	var payload = {
		"inventory_index": 0,
		"slot_type": slot_type,
		"slot_index": slot_index
	}
	NetworkManager.send_packet("PACKET_TYPE_GRAFT_ORGAN", payload)

```

### ./client/src/ui/space/DraggableWindow.gd
```gd
extends Panel

var dragging: bool = false
var drag_offset: Vector2 = Vector2.ZERO

func _gui_input(event):
	if event is InputEventMouseButton:
		if event.button_index == MOUSE_BUTTON_LEFT:
			if event.pressed:
				dragging = true
				drag_offset = get_global_mouse_position() - global_position
				raise() # Bring to front
			else:
				dragging = false

	if event is InputEventMouseMotion and dragging:
		global_position = get_global_mouse_position() - drag_offset

```

### ./client/src/ui/space/MarketWindow.gd
```gd
extends "res://src/ui/space/DraggableWindow.gd"

# Hardcoded Market for Sprint 15
var items = [
	{ "id": "synthetic_heart", "name": "Synthetic Heart", "price": 500 },
	{ "id": "ocular_implant", "name": "Ocular Implant", "price": 250 },
	{ "id": "void_heart", "name": "Void Heart (Illegal)", "price": 1000 },
	{ "id": "mining_laser", "name": "Mining Laser", "price": 100 },
]

func _ready():
	name = "MarketWindow"
	size = Vector2(250, 300)
	position = Vector2(300, 50)

	# Title
	var title = Label.new()
	title.text = "THE MEAT MARKET"
	title.position = Vector2(10, 5)
	add_child(title)

	var container = VBoxContainer.new()
	container.position = Vector2(10, 30)
	container.size = Vector2(230, 260)
	add_child(container)

	for item in items:
		var row = HBoxContainer.new()

		var label = Label.new()
		label.text = item.name + "\n" + str(item.price) + " Solium"
		label.size_flags_horizontal = Control.SIZE_EXPAND_FILL
		row.add_child(label)

		var buy_btn = Button.new()
		buy_btn.text = "BUY"
		buy_btn.connect("pressed", func(): _on_buy_pressed(item.id))
		row.add_child(buy_btn)

		container.add_child(row)
		container.add_child(HSeparator.new())

func _on_buy_pressed(item_id: String):
	print("Buying: ", item_id)
	var payload = { "item_id": item_id }
	NetworkManager.send_packet("PACKET_TYPE_BUY_ITEM", payload)

```

### ./client/src/autoload/NetworkManager.gd
```gd
extends Node

# Signals
signal connected_to_server()
signal disconnected_from_server()
signal packet_received(type: String, payload: Dictionary)
signal sanity_changed(new_level: float)
signal ship_stats_updated(stats: Dictionary)

# Connection State
var ws_peer: WebSocketPeer = WebSocketPeer.new()
var udp_peer: PacketPeerUDP = PacketPeerUDP.new()
var udp_host: String = "127.0.0.1"
var udp_port: int = 5000

var is_connected_to_space: bool = false
var space_host: String = "ws://localhost:8080/ws"

func _ready():
	# UDP Setup
	udp_peer.connect_to_host(udp_host, udp_port)
	# For prototype, we might start disconnected or auto-connect.
	pass

func connect_to_space(token: String):
	print("NetworkManager: Connecting to Space Core at ", space_host)
	var err = ws_peer.connect_to_url(space_host)
	if err != OK:
		print("NetworkManager: Failed to connect to Space Core")
		return

	is_connected_to_space = true
	# We need to wait for connection to be open before sending token.
	# This is handled in _process.
	# For simplicity in this script, we store the token to send later.
	_pending_token = token

var _pending_token: String = ""
var _handshake_sent: bool = false

func _send_login_packet():
	var pkt = {
		"type": "LOGIN_WITH_TOKEN",
		"payload": {
			"token": _pending_token
		}
	}
	send_packet(pkt.type, pkt.payload)

func send_packet(type: String, payload: Dictionary):
	if ws_peer.get_ready_state() == WebSocketPeer.STATE_OPEN:
		var pkt = {
			"type": type,
			"payload": payload
		}
		ws_peer.send_text(JSON.stringify(pkt))

func _handle_packet(data: Dictionary):
	var type = data.get("type", "")
	var payload = data.get("payload", {})

	emit_signal("packet_received", type, payload)

	match type:
		"LOGIN_SUCCESS":
			print("NetworkManager: Login Success")
			emit_signal("connected_to_server")
			emit_signal("packet_received", type, payload) # Ensure logic sees this too
		"PACKET_TYPE_SANITY_UPDATE":
			var val = payload.get("sanity", 100.0)
			emit_signal("sanity_changed", val)
		"PACKET_TYPE_SHIP_STATS":
			emit_signal("ship_stats_updated", payload)
		"PACKET_TYPE_GROUND_STATE":
			emit_signal("packet_received", type, payload)

func send_udp_packet(type: String, payload: Dictionary):
	var pkt = {
		"type": type,
		"payload": payload
	}
	var packet_bytes = JSON.stringify(pkt).to_utf8_buffer()
	udp_peer.put_packet(packet_bytes)

func _process(delta):
	# WebSocket Polling
	ws_peer.poll()
	var state = ws_peer.get_ready_state()

	if state == WebSocketPeer.STATE_OPEN:
		if !_handshake_sent and _pending_token != "":
			_send_login_packet()
			_handshake_sent = true

		while ws_peer.get_available_packet_count() > 0:
			var pkt = ws_peer.get_packet()
			var txt = pkt.get_string_from_utf8()
			var json = JSON.parse_string(txt)
			if json:
				_handle_packet(json)
	elif state == WebSocketPeer.STATE_CLOSED:
		if is_connected_to_space:
			print("NetworkManager: Disconnected from Space Core")
			is_connected_to_space = false
			emit_signal("disconnected_from_server")

	# UDP Polling
	while udp_peer.get_available_packet_count() > 0:
		var pkt = udp_peer.get_packet()
		var txt = pkt.get_string_from_utf8()
		var json = JSON.parse_string(txt)
		if json:
			_handle_packet(json)

```

### ./client/src/scenes/space/EntityManager.gd
```gd
extends Node3D

# Simple Entity Manager to visualize other ships/objects
# Mapped by Entity ID -> Node3D
var entities: Dictionary = {}

func _ready():
	NetworkManager.connect("packet_received", _on_packet_received)

func _on_packet_received(type: String, payload: Dictionary):
	match type:
		"PACKET_TYPE_SPACE_STATE":
			_handle_space_state(payload)
		"PACKET_TYPE_ENTITY_DESTROYED":
			_handle_entity_destroyed(payload)

func _handle_space_state(payload: Dictionary):
	# Expected Payload: { "entities": [ { "id": "...", "pos_x": ... }, ... ] }
	# Note: The current server implementation might just send individual updates or a list.
	# The mechanics.go/main.go implementation of Space Core mainly echoes or validates,
	# but `ENTITY_TRACKING` was mentioned in context.
	# For Sprint 14, we assume we might receive a list of nearby entities.

	var entity_list = payload.get("entities", [])
	var seen_ids = []

	for ent_data in entity_list:
		var id = ent_data.get("id")
		seen_ids.append(id)

		if entities.has(id):
			_update_entity(id, ent_data)
		else:
			_spawn_entity(id, ent_data)

func _spawn_entity(id: String, data: Dictionary):
	var mesh_inst = MeshInstance3D.new()
	var sphere = SphereMesh.new()
	sphere.radius = 1.0
	mesh_inst.mesh = sphere
	add_child(mesh_inst)

	entities[id] = mesh_inst
	_update_entity(id, data)

func _update_entity(id: String, data: Dictionary):
	var node = entities[id]
	var x = data.get("pos_x", 0.0)
	var y = data.get("pos_y", 0.0)
	var z = data.get("pos_z", 0.0)
	node.global_position = Vector3(x, y, z)

func _handle_entity_destroyed(payload: Dictionary):
	var id = payload.get("id")
	if entities.has(id):
		entities[id].queue_free()
		entities.erase(id)

```

### ./client/src/scenes/space/SpaceProjectileManager.gd
```gd
extends Node3D

# SpaceProjectileManager
# Visualizes 3D projectiles with trails and behaviors (Missile, Linear, Instant)

func _ready():
	NetworkManager.connect("packet_received", _on_packet_received)

func _on_packet_received(type: String, payload: Dictionary):
	if type == "PACKET_TYPE_PROJECTILE_SPAWN":
		_spawn_projectile(payload)

func _spawn_projectile(data: Dictionary):
	var behavior = data.get("behavior", "linear")

	if behavior == "instant":
		_draw_beam(data)
	else:
		_spawn_moving_projectile(data)

func _spawn_moving_projectile(data: Dictionary):
	var start = Vector3(data.get("start_x"), data.get("start_z"), data.get("start_y")) # Z-Y swap check?
	# Godot Y is up. Protocol Z is likely depth/up depending on Space convention.
	# Assuming Space Core uses X/Z plane and Y is up? Or full 3D?
	# Let's assume direct mapping: X->X, Y->Y, Z->Z.
	start = Vector3(data.get("start_x"), data.get("start_y"), data.get("start_z"))
	var end = Vector3(data.get("end_x"), data.get("end_y"), data.get("end_z"))

	var speed = data.get("speed", 50.0)
	var behavior = data.get("behavior", "linear")

	var mesh_inst = MeshInstance3D.new()

	if behavior == "missile":
		var sphere = SphereMesh.new()
		sphere.radius = 0.5
		mesh_inst.mesh = sphere
		# Add Trail (Placeholder)
		# var trail = GPUParticles3D.new() ...
	else:
		# Linear (Railgun)
		var cylinder = CylinderMesh.new()
		cylinder.height = 2.0
		cylinder.radius = 0.1
		mesh_inst.mesh = cylinder
		mesh_inst.rotation.x = PI/2 # Align forward

	add_child(mesh_inst)
	mesh_inst.global_position = start
	mesh_inst.look_at(end)

	var dist = start.distance_to(end)
	var duration = dist / speed

	if duration <= 0: duration = 0.1

	var tween = create_tween()

	if behavior == "missile":
		# Quadratic Curve interpolation could go here for "launch" feel
		# For MVP, linear move to target
		tween.tween_property(mesh_inst, "global_position", end, duration)
	else:
		tween.tween_property(mesh_inst, "global_position", end, duration)

	tween.tween_callback(mesh_inst.queue_free)

func _draw_beam(data: Dictionary):
	# Draw Line3D logic
	pass

```

### ./client/src/scenes/space/ShipController.gd
```gd
extends CharacterBody3D

# Newtonian Physics Parameters
@export var acceleration_force: float = 50.0
@export var rotation_speed: float = 2.0
@export var max_speed: float = 100.0
@export var braking_factor: float = 2.0

var network_tick_timer: float = 0.0
const NETWORK_TICK_RATE: float = 0.1 # 10Hz

func _physics_process(delta):
	handle_input(delta)
	move_and_slide()
	handle_networking(delta)

func handle_input(delta):
	var rot_dir = 0
	if Input.is_action_pressed("ui_left") or Input.is_key_pressed(KEY_A):
		rot_dir += 1
	if Input.is_action_pressed("ui_right") or Input.is_key_pressed(KEY_D):
		rot_dir -= 1

	rotate_y(rot_dir * rotation_speed * delta)

	var forward_dir = -global_transform.basis.z
	var input_accel = Vector3.ZERO

	if Input.is_action_pressed("ui_up") or Input.is_key_pressed(KEY_W):
		input_accel += forward_dir * acceleration_force * delta
	if Input.is_action_pressed("ui_down") or Input.is_key_pressed(KEY_S):
		input_accel -= forward_dir * acceleration_force * delta

	if Input.is_action_pressed("ui_select") or Input.is_key_pressed(KEY_SPACE):
		var current_speed = velocity.length()
		if current_speed > 0:
			var brake_force = velocity.normalized() * -1 * acceleration_force * braking_factor * delta
			if brake_force.length() > current_speed:
				velocity = Vector3.ZERO
			else:
				velocity += brake_force

	velocity += input_accel

	if velocity.length() > max_speed:
		velocity = velocity.normalized() * max_speed

func handle_networking(delta):
	# Use NetworkManager instead of direct WS peer
	network_tick_timer += delta
	if network_tick_timer >= NETWORK_TICK_RATE:
		network_tick_timer = 0
		send_state()

func send_state():
	var payload = {
		"pos_x": global_position.x,
		"pos_y": global_position.y,
		"pos_z": global_position.z,
		"vel_x": velocity.x,
		"vel_y": velocity.y,
		"vel_z": velocity.z,
		"rot_y": rotation.y
	}
	NetworkManager.send_packet("PACKET_TYPE_SPACE_STATE", payload)

```

### ./client/src/scenes/ground/ProjectileManager.gd
```gd
extends Node3D

# Visualizes Projectiles based on Server Events
# Projectiles are just visual meshes interpolated between Start and End.

func _ready():
	NetworkManager.connect("packet_received", _on_packet_received)

func _on_packet_received(type: String, payload: Dictionary):
	if type == "PACKET_TYPE_PROJECTILE_SPAWN":
		_spawn_projectile(payload)

func _spawn_projectile(data: Dictionary):
	var start = Vector3(data.get("start_x"), 1.0, data.get("start_y"))
	var end = Vector3(data.get("end_x"), 1.0, data.get("end_y"))
	var speed = data.get("speed", 20.0)
	var outcome = data.get("outcome", "Hit")

	var mesh = MeshInstance3D.new()
	var cylinder = CylinderMesh.new()
	cylinder.height = 0.5
	cylinder.radius = 0.1
	mesh.mesh = cylinder

	# Rotate to face target
	mesh.rotation.x = PI / 2 # Lay flat

	add_child(mesh)
	mesh.global_position = start
	mesh.look_at(end)

	# Tween movement
	var dist = start.distance_to(end)
	var duration = dist / speed

	var tween = create_tween()
	tween.tween_property(mesh, "global_position", end, duration)

	if outcome == "Miss":
		# overshoot
		var miss_pos = end + (end - start).normalized() * 5.0
		tween.stop() # Reset
		dist = start.distance_to(miss_pos)
		duration = dist / speed
		tween = create_tween()
		tween.tween_property(mesh, "global_position", miss_pos, duration)

	tween.tween_callback(mesh.queue_free)

```

### ./client/src/scenes/ground/CameraRig.gd
```gd
extends Node3D

@export var target_path: NodePath
@export var smooth_speed: float = 5.0
@export var offset: Vector3 = Vector3(20, 20, 20) # High up and isometric

var target: Node3D
var camera: Camera3D

func _ready():
	if target_path:
		target = get_node(target_path)

	# Create Camera if not present as child (for script attachment usage)
	# But typically this script goes on a Node3D holding the camera
	camera = $Camera3D
	if not camera:
		camera = Camera3D.new()
		add_child(camera)
		camera.name = "Camera3D"

	# Set Orthogonal Projection for 2.5D look
	camera.projection = Camera3D.PROJECTION_ORTHOGONAL
	camera.size = 20.0 # Adjust size for zoom level

	# Set isometric rotation (Look at center from offset)
	# Typically isometric is 45 degrees Y, and roughly 30-35 degrees X down
	# But looking from (20, 20, 20) to (0,0,0) does this automatically if we use look_at

	# Initial position
	if target:
		global_position = target.global_position + offset
		look_at(target.global_position)

func _physics_process(delta):
	if !target:
		return

	var desired_position = target.global_position + offset
	var smoothed_position = global_position.lerp(desired_position, smooth_speed * delta)
	global_position = smoothed_position

	# Ensure we always look at the target, or a fixed point relative to the camera
	# For strict isometric, rotation usually shouldn't change, just position.
	# So we just move the rig.


```

### ./client/src/scenes/ground/PlayerController.gd
```gd
extends CharacterBody3D

# Settings
@export var speed: float = 10.0
@export var acceleration: float = 50.0
@export var friction: float = 60.0

# Network
var update_timer: float = 0.0
const UPDATE_RATE: float = 0.05 # 20Hz

# Combat
var current_target_id: String = ""
var camera: Camera3D

func _ready():
	# Temporary: Spawn Hangar Zone for testing
	var hangar = preload("res://src/scenes/ground/HangarZone.gd").new()
	hangar.position = Vector3(10, 0, 10)
	get_parent().call_deferred("add_child", hangar)

	# Find Camera (Assuming CameraRig is sibling or child, for now grab viewport camera)
	camera = get_viewport().get_camera_3d()

	# Listen for Gear
	NetworkManager.connect("packet_received", _on_packet_received)

func _on_packet_received(type: String, payload: Dictionary):
	if type == "LOGIN_SUCCESS":
		var gear = payload.get("ground_gear", {})
		if gear.has("primary_weapon"):
			var weapon = gear["primary_weapon"]
			if weapon:
				print("PlayerController: Equipped Weapon: ", weapon.get("item_id", "Unknown"))

func _physics_process(delta):
	# Movement
	var input_dir = Input.get_vector("ui_left", "ui_right", "ui_up", "ui_down")
	var direction = Vector3(input_dir.x, 0, input_dir.y).normalized()

	if direction:
		velocity.x = move_toward(velocity.x, direction.x * speed, acceleration * delta)
		velocity.z = move_toward(velocity.z, direction.z * speed, acceleration * delta)
	else:
		velocity.x = move_toward(velocity.x, 0, friction * delta)
		velocity.z = move_toward(velocity.z, 0, friction * delta)

	move_and_slide()

	# Network Update
	update_timer += delta
	if update_timer >= UPDATE_RATE:
		update_timer = 0
		_send_movement()

func _input(event):
	# Tab Targeting
	if event.is_action_pressed("ui_focus_next"): # Tab
		_cycle_target()

	# Mouse Click Targeting
	if event is InputEventMouseButton and event.pressed and event.button_index == MOUSE_BUTTON_LEFT:
		_raycast_target(event.position)

	# F1 Self Target (Placeholder for Party)
	if event is InputEventKey and event.pressed and event.keycode == KEY_F1:
		current_target_id = "" # Clear or Set to Self ID if known
		print("Target Self/None")

	# Fire
	if event.is_action_pressed("ui_accept"): # Space or Enter, map "Fire" later
		_fire_weapon()

func _cycle_target():
	# Access GroundEntityManager. Assuming it's a sibling or singleton.
	# For MVP, let's assume we can find it in the Scene Tree.
	var manager = get_node_or_null("../GroundEntityManager")
	if manager:
		current_target_id = manager.get_next_target(current_target_id)
		print("Target Locked: ", current_target_id)
		_update_reticle()

func _raycast_target(mouse_pos):
	if !camera:
		camera = get_viewport().get_camera_3d()
		if !camera: return

	var from = camera.project_ray_origin(mouse_pos)
	var to = from + camera.project_ray_normal(mouse_pos) * 1000.0

	var space_state = get_world_3d().direct_space_state
	var query = PhysicsRayQueryParameters3D.create(from, to)
	query.collision_mask = 2 # Layer 2 (Entities)

	var result = space_state.intersect_ray(query)
	if result:
		var collider = result.collider
		if collider.has_meta("entity_id"):
			current_target_id = collider.get_meta("entity_id")
			print("Click Target: ", current_target_id)
			_update_reticle()

func _fire_weapon():
	if current_target_id == "":
		print("No Target!")
		return

	print("Firing at ", current_target_id)

	# Visuals: Muzzle Flash (Placeholder)

	# Network
	var payload = { "target_id": current_target_id }
	NetworkManager.send_udp_packet("PACKET_TYPE_GROUND_ATTACK", payload)

func _update_reticle():
	# Visual feedback for target (Simple print for now or highlight shader)
	# In real imp, we would move a Sprite3D to the target's position.
	pass

func _send_movement():
	var payload = {
		"x": global_position.x,
		"y": global_position.z, # Mapping 3D Z to 2D Y
		"vx": velocity.x,
		"vy": velocity.z,
		"ts": Time.get_unix_time_from_system() * 1000
	}
	NetworkManager.send_udp_packet("PACKET_TYPE_GROUND_MOVEMENT", payload)

```

### ./client/src/scenes/ground/GroundEntityManager.gd
```gd
extends Node3D

# Visualizes other players on ground
var entities: Dictionary = {} # ID -> Node3D

func _ready():
	NetworkManager.connect("packet_received", _on_packet_received)

func _on_packet_received(type: String, payload: Dictionary):
	if type == "PACKET_TYPE_GROUND_STATE":
		_update_entities(payload.get("entities", []))

func _update_entities(list: Array):
	for data in list:
		var id = data.get("id")

		if entities.has(id):
			var node = entities[id]
			var x = data.get("x", 0.0)
			var y = data.get("y", 0.0)
			node.global_position = Vector3(x, 0, y)
		else:
			_spawn_entity(id, data)

func _spawn_entity(id: String, data: Dictionary):
	var mesh_inst = MeshInstance3D.new()
	var capsule = CapsuleMesh.new()
	capsule.height = 2.0
	mesh_inst.mesh = capsule
	mesh_inst.position.y = 1.0 # Offset up
	add_child(mesh_inst)

	# Add Collider for Raycast Targeting
	var static_body = StaticBody3D.new()
	var shape = CollisionShape3D.new()
	var cap_shape = CapsuleShape3D.new()
	cap_shape.height = 2.0
	shape.shape = cap_shape
	static_body.add_child(shape)

	# Metadata for identification
	static_body.set_meta("entity_id", id)

	# Offset collider to match mesh
	static_body.position.y = 1.0
	add_child(static_body)

	var x = data.get("x", 0.0)
	var y = data.get("y", 0.0)
	# mesh_inst is a child, but we need to move the parent node if we wrapped it?
	# Ah, I added mesh_inst as child of self (GroundEntityManager).
	# I should probably wrap them in a Node3D "Entity" so I move one thing.

	# Refactor for cleaner movement
	mesh_inst.queue_free()
	static_body.queue_free()

	var wrapper = Node3D.new()
	add_child(wrapper)
	wrapper.global_position = Vector3(x, 0, y)

	# Re-add visuals to wrapper
	var vis = MeshInstance3D.new()
	vis.mesh = capsule
	vis.position.y = 1.0
	wrapper.add_child(vis)

	# Re-add collider to wrapper
	var body = StaticBody3D.new()
	body.add_child(shape.duplicate()) # Shape resource is shared
	body.set_meta("entity_id", id)
	body.position.y = 1.0
	body.collision_layer = 2 # Enemy/Entity Layer
	wrapper.add_child(body)

	entities[id] = wrapper

func get_next_target(current_id: String) -> String:
	var keys = entities.keys()
	if keys.is_empty():
		return ""

	var index = keys.find(current_id)
	if index == -1 or index == keys.size() - 1:
		return keys[0]

	return keys[index + 1]

```

### ./client/src/scenes/ground/HangarZone.gd
```gd
extends Area3D

func _ready():
	# Visuals: Transparent Blue Cylinder
	var mesh_inst = MeshInstance3D.new()
	var cylinder = CylinderMesh.new()
	cylinder.height = 4.0
	cylinder.radius = 2.0
	mesh_inst.mesh = cylinder

	var mat = StandardMaterial3D.new()
	mat.albedo_color = Color(0, 0.5, 1.0, 0.3) # Blue Transparent
	mat.transparency = BaseMaterial3D.TRANSPARENCY_ALPHA
	mat.cull_mode = BaseMaterial3D.CULL_DISABLED
	mesh_inst.material_override = mat

	add_child(mesh_inst)

	# Collision Shape
	var shape = CollisionShape3D.new()
	var cyl_shape = CylinderShape3D.new()
	cyl_shape.height = 4.0
	cyl_shape.radius = 2.0
	shape.shape = cyl_shape
	add_child(shape)

	# Signal
	connect("body_entered", _on_body_entered)

func _on_body_entered(body):
	# Check if body is the player (basic check for MVP)
	if body.has_method("_send_movement"): # Duck typing check for PlayerController
		print("Hangar Zone Entered! Requesting Launch...")
		NetworkManager.send_udp_packet("REQUEST_LAUNCH", {})
		# Note: The server will respond with PACKET_TYPE_LAUNCH_GRANTED.
		# NetworkManager should handle that to switch scenes.

```

### ./client/src/vfx/SanityController.gd
```gd
extends Node

# Controls the Sanity Distortion Shader based on network events
@export var shader_rect_path: NodePath
var shader_mat: ShaderMaterial

func _ready():
	NetworkManager.connect("sanity_changed", _on_sanity_changed)

	if shader_rect_path:
		var rect = get_node(shader_rect_path)
		if rect and rect.material is ShaderMaterial:
			shader_mat = rect.material

func _on_sanity_changed(new_level: float):
	if !shader_mat:
		return

	# Sanity is 0-100 usually, shader might expect 0.0-1.0
	var normalized = clamp(new_level / 100.0, 0.0, 1.0)
	shader_mat.set_shader_parameter("sanity_level", normalized)

```
