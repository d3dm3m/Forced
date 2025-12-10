-- Migration for Sprint 12 Phase 2: Player Economy and Skills

-- Add solium column to players table (default 0)
ALTER TABLE players ADD COLUMN IF NOT EXISTS solium INT DEFAULT 0;

-- Add skills column to players table (JSONB)
ALTER TABLE players ADD COLUMN IF NOT EXISTS skills JSONB DEFAULT '{}'::jsonb;
