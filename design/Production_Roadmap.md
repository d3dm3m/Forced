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


## Phase 4: The Industrial Engine (Months 7-8)
**Goal:** Transform the MVP into a "Survival Strategy RPG" by wiring stats to gameplay logic and enforcing server authority.

*   **Sprint 21 (Weeks 21-22): The Motherboard**
    *   Implement the "Circuit Grid" inventory backend.
    *   Update `GroundGear` to support socketed "Chips" (Items with stat modifiers).
*   **Sprint 22 (Weeks 23-24): Stat Wiring**
    *   Connect Torque, Compute, Synapse, Flux to game logic.
        *   **Torque:** Inventory Mass Limit / Recoil.
        *   **Compute:** Crafting Speed / Drone Count.
        *   **Synapse:** Turn Rate / Cast Point speed.
        *   **Flux:** Shield Regen / Ability Cooldowns.
*   **Sprint 23 (Weeks 25-26): The Configurator**
    *   Build the Client UI for the Motherboard.
    *   Drag-and-drop Chips into Sockets to change stats dynamically.
*   **Sprint 24 (Weeks 27-28): Server Authority**
    *   Implement server-side validation for the Action State Machine.
    *   Reject packets if the player is in Backswing or Windup.
*   **Sprint 25 (Weeks 29-30): Ability System**
    *   Create a generic `AbilityManager` on the server.
    *   Implement the first 3 unique skills:
        *   **Breacher:** Phalanx Shield (Directional damage reduction).
        *   **Null-Walker:** Phase Shift (Teleport/Disjoint).
        *   **Sapper:** Deploy Turret (Spawn Entity).
