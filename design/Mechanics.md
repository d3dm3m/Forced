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

## 8. Space Combat Physics (The Math of War)

### Angular Ballistics
Space combat is not twitch-based; it is calculated based on physics and angular velocities.
*   **Turret Tracking:** Hit chance is determined by the Weapon Slew Rate versus the Target's Angular Velocity relative to the shooter.
    *   **Transversal Movement:** High angular velocity (orbiting). Safe.
    *   **Radial Movement:** Zero angular velocity (burning straight at/away). Dead.
*   **Formula:** Hit Chance degrades as Angular Velocity exceeds Tracking Speed.

### Volumetric Detonation
Missiles operate on a different paradigm. They always hit (if in range), but damage is applied based on the explosion's ability to catch the target.
*   **Explosion Radius vs Signature Radius:** Small targets take less damage from big explosions.
*   **Explosion Velocity vs Target Velocity:** Fast targets outrun the shockwave.

### Signature Analysis
Target locking is an active sensor process.
*   **Sensor Cross-Section (SCS):** The "size" of the ship on radar.
*   **Scan Resolution:** The speed of the targeting sensors.
*   **Lock Time:** Defined by `Scanner Resolution / Target SCS`. Active modules (MWD, Jammers) bloom the SCS, making the ship faster to lock.

### Capacitor Warfare
Energy is life. The Capacitor powers shields, weapons, and propulsion.
*   **Recharge:** Non-linear. Recharge rate peaks at ~30% capacity and drops off at 0% and 100%.
*   **Warfare:** Energy Neutralizers can drain enemy caps, leaving them dead in space.

### Tackling
Preventing escape is a dedicated role.
*   **Warp Jammers:** Prevent the target from entering warp.
*   **Webifiers:** Artificial gravity drag. Reduces target speed, which lowers their Transversal, making them easier to hit.

### Layered Mitigation
*   **Shields:** Regenerating, weak to EM. First line of defense.
*   **Armor:** Static HP, high Kinetic resistance. Reduces speed when heavy plates are installed.
*   **Hull:** The structure. No resistances. When this hits 0, the ship explodes.

## 9. Tactical Sensor-Link (Ground Combat v2.0)

### Entity Physics
Movement is "Weighty" and deliberate, moving away from twitch-shooters to RTS-style tactical positioning.
*   **Turn-Rate:** Characters cannot move instantly in a new direction. They must rotate (Turn Rate) to face the target vector before Translation begins.
*   **Movement Threshold:** Entities only begin moving once facing is within ~15 degrees of the target vector. This makes "kiting" difficult for heavy frames.

### Sensor Vision (Fog of War)
Vision is calculated via Raycast, not simple distance checks.
*   **Layers:** Ground (0), Catwalk (1), Obstruction (2).
*   **High Ground Advantage:** High ground sees Low ground freely. Low ground cannot see up to High ground unless they have a spotter or active sensor sweep.
*   **Occlusion:** Obstacles block vision rays, creating dynamic shadows where enemies can hide.

### Action State Machine
All abilities and attacks follow a rigorous state machine to prevent animation canceling exploits and enforce commitment.
1.  **Idle:** Ready to act.
2.  **Windup (Cast Point):** The preparation phase. Can be canceled to bait enemies. Resource is not yet consumed.
3.  **Active:** The effect occurs (Projectile fired, Heal applied). Resource is burned.
4.  **Backswing:** The recovery phase. Animation lock prevents moving or attacking immediately, but can be canceled by Move commands in some agile frames.

### Ballistics & Mitigation
*   **Projectile Entities:** Ranged attacks are physical entities that travel through space. They can be dodged or disjointed (e.g., blinking/teleporting).
*   **High Ground Defense:** Attacking High Ground from Low Ground incurs a 25% Miss Chance (Uphill Battle).

### Environmental Interaction
*   **Scrap Piles:** Destructible cover elements scattered in the world.
    *   **Tactical:** They block vision and pathing.
    *   **Strategic:** Sapper/Biologist classes can salvage them for "Nanite Repair" or resources.
