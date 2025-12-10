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
