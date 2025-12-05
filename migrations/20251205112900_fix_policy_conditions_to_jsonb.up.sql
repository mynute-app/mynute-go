-- Fix PolicyRule.conditions column type from TEXT to JSONB
-- The conditions column should store JSON data for policy conditions

ALTER TABLE public.policy_rules 
ALTER COLUMN "conditions" TYPE JSONB USING "conditions"::JSONB;
