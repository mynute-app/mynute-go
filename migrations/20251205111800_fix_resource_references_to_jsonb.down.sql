-- Rollback Resource.references column type from JSONB to TEXT

ALTER TABLE public.resources 
ALTER COLUMN "references" TYPE TEXT USING "references"::TEXT;
