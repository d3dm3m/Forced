-- Migration for Sprint 13 Phase 2: Add Current Health

-- Add current_health column to players table (default 1000.0)
ALTER TABLE players ADD COLUMN IF NOT EXISTS current_health DOUBLE PRECISION DEFAULT 1000.0;
