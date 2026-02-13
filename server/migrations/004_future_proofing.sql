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
