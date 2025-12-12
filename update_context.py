import os
import datetime

OUTPUT_FILE = "PROJECT_CONTEXT.md"
IGNORE_DIRS = {".git", "bin", "__pycache__", ".godot", ".import", "tmp", "vendor"}
IGNORE_FILES = {".DS_Store", "go.sum", "go.mod"} # minimal ignore
INCLUDE_EXTENSIONS = {".go", ".gd", ".md", ".json", ".sql", ".py", ".gdshader", ".tscn", ".tres"}

def get_tree_structure(path):
    tree_str = ""
    for root, dirs, files in os.walk(path):
        dirs[:] = [d for d in dirs if d not in IGNORE_DIRS]
        level = root.replace(path, '').count(os.sep)
        indent = '    ' * level
        tree_str += '{}{}/\n'.format(indent, os.path.basename(root))
        for f in files:
            if f not in IGNORE_FILES:
                tree_str += '{}{}\n'.format(indent + '    ', f)
    return tree_str

def get_file_census(path):
    counts = {}
    for root, dirs, files in os.walk(path):
        dirs[:] = [d for d in dirs if d not in IGNORE_DIRS]
        for f in files:
            ext = os.path.splitext(f)[1]
            counts[ext] = counts.get(ext, 0) + 1
    return counts

def get_mermaid_mindmap():
    return """mindmap
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
"""

def get_mermaid_critical_path():
    return """graph TD
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
"""

def main():
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    with open(OUTPUT_FILE, "w", encoding="utf-8") as out:
        out.write(f"# PROJECT CONTEXT\n\n")
        out.write(f"**Last Updated:** {timestamp}\n\n")

        # 0. PERSONA ROSTER
        out.write("## 🤖 AI Persona Roster\n")
        out.write("* **The Architect:** System Design, Database Schema, Network Topology. (Use for: Infrastructure)\n")
        out.write("* **The Void Engineer:** Go Backend, Physics, Concurrency. (Use for: Server Logic)\n")
        out.write("* **The Operator:** Godot Engine, GDScript, Shaders, UI. (Use for: Client)\n")
        out.write("* **The Biologist:** Game Design, Balancing, Lore, Item Configs. (Use for: Mechanics)\n")
        out.write("* **The Scribe:** Documentation, Context Management, Verification. (Use for: Organization)\n\n")

        # 1. VISUALS
        out.write("## Project Visuals\n\n")
        out.write("### System Architecture (Mindmap)\n")
        out.write("```mermaid\n")
        out.write(get_mermaid_mindmap())
        out.write("```\n\n")

        out.write("### Critical Path (Dependency Graph)\n")
        out.write("```mermaid\n")
        out.write(get_mermaid_critical_path())
        out.write("```\n\n")

        # 2. STATUS (Epics)
        out.write("## Project Status (Epics)\n\n")
        out.write("### Completed Sprints: 12-17 (Foundation of Bio-Economy & Twin-Core)\n")
        out.write("- [x] **Sprint 12 (Data):** Refactored Inventory (`ItemStack` metadata), Ship Layout (Tiered Slots), and Ground Gear.\n")
        out.write("- [x] **Sprint 13 (Logic):** Implemented `Mechanics` core (Stat Aggregation, Bio-Load) and integrated `Main` game loop.\n")
        out.write("- [x] **Sprint 14 (Client):** Restored Architecture (`NetworkManager`, `SanityController`) and built HUD (Bio-Feedback).\n")
        out.write("- [x] **Sprint 15 (Integration):** Wired `MarketService` and `SurgeryService` to frontend UI (`MarketWindow`, `SurgeryWindow`).\n")
        out.write("- [x] **Sprint 16 (Ground):** Established Ground Gameplay Loop (20Hz UDP, Validation, Broadcasting).\n")
        out.write("- [x] **Sprint 17 (Tether):** Implemented Ground Persistence and Hangar Handoff trigger.\n")
        out.write("- [x] **Sprint 18 (Expansion):** Implemented Gatekeeper Service and SystemID persistence.\n")
        out.write("- [x] **Sprint 18 (Mechanics):** Implemented Passive Ship Simulation (Shield/Capacitor Regen).\n")
        out.write("- [x] **Sprint 19 (Ground):** Tactical Physics & Turn-Rate Movement.\n")
        out.write("- [x] **Sprint 19.2 (Ground):** Client-Side Raycast Fog of War (Shadow System).\n")
        out.write("- [x] **Sprint 19.3 (Ground):** Action State Machine (Cast Point logic).\n")
        out.write("- [x] **Sprint 19.4 (Ground):** Ground Damage Resolution & Respawn.\n")
        out.write("- [x] **Sprint 20 (Data):** Implemented Industrial Stats & 10-Class Roster.\n")
        out.write("- [ ] **Sprint 21 (Data):** The Motherboard (Circuit Grid & Chips).\n")
        out.write("- [ ] **Sprint 22 (Logic):** Stat Wiring (Torque/Compute/Synapse/Flux).\n")
        out.write("- [ ] **Sprint 23 (Client):** The Configurator UI.\n")
        out.write("- [ ] **Sprint 24 (Server):** Server Authority (Action State Validation).\n")
        out.write("- [ ] **Sprint 25 (Logic):** Ability System (Breacher, Null-Walker, Sapper skills).\n\n")

        out.write("## The Macro-Scale Architecture (Planned)\n")
        out.write("* **Zone Sharding:** The universe is split into `Systems`. Each System can be hosted on a different physical server node. The `IGatekeeper` interface will manage routing.\n")
        out.write("* **Gatekeeper:** A dedicated service to track which server hosts which system and facilitate handoffs (e.g., `Packet_JumpGate`).\n")
        out.write("* **Planet Hazards:** Planets will have JSON-defined environmental variables (`Gravity`, `Atmosphere`, `ThermalRating`) that the Physics Engine must respect.\n")
        out.write("* **Social (Guilds):** Guild ownership is baked into the Player/Entity relationship via `guild_id`. This allows robust permission checks (e.g. Door Access).\n")
        out.write("* **Simulation (Status Effects):** Active Effects are handled via a generic `JSONB` container to allow for infinite expansion (Buffs, Diseases) without schema changes.\n")
        out.write("* **Legacy (Account Separation):** The `accounts` table stores permanent data (Cosmetics, Legacy Currency) separate from the wipeable `players` table.\n\n")

        out.write("### Technical Debt & Future Focus\n")
        out.write("- [x] **Ground Core Lag:** The `GroundGear` data structures exist on the server but are not used by the Client or Ground Core networking.\n")
        out.write("- [x] **Inventory UI:** `InventoryUI.gd` and `SurgeryWindow.gd` now support full Drag-and-Drop interaction.\n")
        out.write("- [x] **Space Core Combat Logic Integration:** Wired up Angular Ballistics and Signature Analysis to the main game loop.\n")
        out.write("- [ ] **Gatekeeper Real-Implementation:** Gatekeeper currently uses a mocked routing table; needs Redis backing.\n")
        out.write("- [ ] **Ground Core Refactor:** Implement RTS Turn-Rate Physics & Raycast Vision.\n")
        out.write("- [ ] **Next Goal:** Sprint 21: The Motherboard (Perk Architecture).\n\n")

        out.write("## Directory Tree\n\n```\n")
        out.write(get_tree_structure("."))
        out.write("```\n\n")

        out.write("## File Census\n\n")
        counts = get_file_census(".")
        for ext, count in counts.items():
            out.write(f"- **{ext}**: {count}\n")
        out.write("\n")

        # 3. CONTENT
        out.write("## File Contents\n\n")
        for root, dirs, files in os.walk("."):
            dirs[:] = [d for d in dirs if d not in IGNORE_DIRS]
            for f in files:
                if f == OUTPUT_FILE or f == "update_context.py" or f in IGNORE_FILES:
                    continue

                ext = os.path.splitext(f)[1]
                path = os.path.join(root, f)

                if ext in INCLUDE_EXTENSIONS:
                    out.write(f"### {path}\n")
                    try:
                        with open(path, "r", encoding="utf-8") as infile:
                            content = infile.read()
                            # Wrap in code block
                            lang = ext.replace(".", "")
                            if lang == "gdshader": lang = "glsl"
                            out.write(f"```{lang}\n{content}\n```\n\n")
                    except Exception as e:
                        out.write(f"*Error reading file: {e}*\n\n")
                else:
                     out.write(f"### {path}\n*Binary/Asset File*\n\n")

if __name__ == "__main__":
    main()
