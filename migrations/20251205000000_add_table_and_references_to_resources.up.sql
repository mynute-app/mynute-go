-- Migration: add_table_and_references_to_resources
-- Created at: 20251205000000
-- Description: Add missing table and references columns to resources table

ALTER TABLE public.resources 
ADD COLUMN IF NOT EXISTS "table" VARCHAR(255),
ADD COLUMN IF NOT EXISTS "references" JSONB;
