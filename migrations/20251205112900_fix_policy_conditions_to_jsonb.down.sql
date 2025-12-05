-- Rollback PolicyRule.conditions column type from JSONB to TEXT

ALTER TABLE public.policy_rules 
ALTER COLUMN "conditions" TYPE TEXT USING "conditions"::TEXT;
