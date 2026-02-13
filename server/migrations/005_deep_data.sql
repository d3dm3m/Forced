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
