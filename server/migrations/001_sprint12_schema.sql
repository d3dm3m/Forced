-- Migration for Sprint 12: Bio-Grafting and Ground Gear

-- Add ground_gear column to players table
ALTER TABLE players ADD COLUMN IF NOT EXISTS ground_gear JSONB DEFAULT '{}'::jsonb;

-- Note: inventory and ship_layout columns are already JSONB.
-- The application logic will handle the structure change.
-- For a real production migration, we might need to transform existing data,
-- but for Sprint 12 dev, we assume backwards compatibility or reset.
