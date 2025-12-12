-- Migration for Sprint 21: The Motherboard (Industrial Stats & Inventory)

-- Add motherboard column to players table (JSONB)
ALTER TABLE players ADD COLUMN IF NOT EXISTS motherboard JSONB DEFAULT '{"slots": 0, "chips": [], "stats": {"torque": 0, "compute": 0, "synapse": 0, "flux": 0}}'::jsonb;
