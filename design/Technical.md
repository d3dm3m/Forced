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
