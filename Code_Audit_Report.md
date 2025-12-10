# Code Audit Report: Phase 3 Beta Readiness

## 1. Architecture Integrity Check

*   **Twin-Core Model:** The project structure reflects the Twin-Core design with `server/cmd/ground/main.go` existing. However, the Space Core (`server/cmd/space/main.go`) is currently missing from the codebase, despite being mentioned in `docker-compose.yml` (Sprint 10) and the original plan.
    *   **Verdict:** **Failed**. The Space Core implementation is incomplete or missing.
*   **Handoff Protocol (Redis):** While `docker-compose.yml` includes a Redis service, there is **no code** in the `server/` directory that imports or uses a Redis client (verified via grep). The Handoff Protocol described in `Technical.md` (saving state to Redis, token exchange) is unimplemented.
    *   **Verdict:** **Failed**. The persistence layer relies solely on PostgreSQL, likely causing data race issues during server switching if the Space Core existed.
*   **Headless Validator:** There is no "Headless Godot" integration. The project uses simple server-side Go logic for validation (e.g., `ValidateVector` in `physics.go` from Sprint 5 history, but file check failed). The "Anti-Cheat" is rudimentary logic, not a true headless simulation.
    *   **Verdict:** **Partial/Placeholder**.

## 2. Feature Gap Analysis (The 'Forgotten' List)

Comparing `Mechanics.md` to the codebase:

*   **Bio-Grafting Rejection:** **Missing**. No logic for "Bleeding," "Spasms," or "Necrosis" exists in `player.go` or `gamedata.go`. The `Player` struct has `ShipLayout`, but no organ tier tracking or rejection mechanics.
*   **Gate Sickness:** **Missing**. No "Sanity Meter" or debuff logic tied to FTL travel.
*   **Sanity System:** **Partial**. `Sanity` field added to `Player` struct (Sprint 7), and basic decay logic exists in `player.go`. However, the "Phantom" injection and deep mechanic integration are minimal.
*   **Economy (Meat Market):** **Missing**. No implementation of "Solium," "Biomass" refining, or "Stable Tissue."
*   **Territory Control:** **Missing**. No "Claim Flags" or "Vulnerability Window" logic found.
*   **Ground Combat:** **Basic**. Simple `Projectile` and `DamageNumber` logic exists, but the "Light/Dark" stealth mechanic is missing.

## 3. Technical Debt & Risk

*   **Concurrency:**
    *   `ChatManager` in `server/internal/game/chat.go` uses `sync.RWMutex`, which is good.
    *   However, `ChatManager` sessions map usage needs careful review to ensure no race conditions during rapid connect/disconnect cycles.
*   **Hardcoded Values:**
    *   `server/internal/db/db.go` contains hardcoded database credentials (`postgres://user:password@localhost...`) as a fallback. This is a security risk.
    *   `server/internal/game/gamedata.go` uses hardcoded `Classes` and `Items` definitions instead of loading from the JSON files created in Sprint 9 (`ships.json`, `classes.json` were created but the loader code provided in Sprint 9 might have been overwritten or not fully integrated as the read of `gamedata.go` shows hardcoded values).
*   **SQL Sanitization:**
    *   `server/internal/game/player.go` uses `pgx` parameterized queries (`$1`, `$2`...), which is **Safe**.
    *   **Critical Bug:** The `player.go` code passes `[]byte` (from `json.Marshal`) directly to `jsonb` columns. PostgreSQL will likely reject this with a type mismatch error (`bytea` vs `jsonb`). These should be cast to string or used with a driver-compatible JSON wrapper.

## 4. Asset Inventory Status

*   **JSON Files:** `server/assets/data/ships.json` exists and contains the requested Kestrel, Omen, and Hauler data.
*   **Missing Files:** `server/assets/data/classes.json` was not found during the file listing, despite being part of the Sprint 9 plan.
*   **Loader Status:** `server/internal/game/gamedata.go` currently populates data via hardcoded Go structs, ignoring the `ships.json` file. The "Content Injection" task is effectively incomplete in the active codebase.

## 5. Remediation Plan (Must Fix for RC1)

1.  **Fix Persistence Layer (Critical):** Update `server/internal/game/player.go` to properly cast JSON byte slices to strings/types compatible with PostgreSQL `jsonb` columns to prevent runtime errors.
2.  **Integrate JSON Loaders:** Modify `server/internal/game/gamedata.go` to actually parse `ships.json` and the missing `classes.json` instead of using hardcoded values.
3.  **Implement Handoff Protocol:** Add Redis client support to the server and implement the session token exchange defined in `Technical.md` to support the Twin-Core architecture.
4.  **Create Space Core:** Implement the missing `server/cmd/space/main.go` entry point.
5.  **Remove Hardcoded Secrets:** Switch `server/internal/db/db.go` to strictly use environment variables and remove the fallback hardcoded credentials.
6.  **Mechanic Implementation:** Select at least one core "Bio-Horror" mechanic (e.g., basic Organ Rejection stats) and implement it server-side to justify the "Bio-Horror" genre tag.
