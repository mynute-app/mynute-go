-- Fix Resource.references column type from TEXT to JSONB
-- The references column should store JSON data for resource references

ALTER TABLE public.resources 
ALTER COLUMN "references" TYPE JSONB USING "references"::JSONB;
