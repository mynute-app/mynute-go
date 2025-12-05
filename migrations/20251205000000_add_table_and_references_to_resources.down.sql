-- Migration: add_table_and_references_to_resources (DOWN)
-- Created at: 20251205000000
-- Description: Remove table and references columns from resources table

ALTER TABLE public.resources 
DROP COLUMN IF EXISTS "table",
DROP COLUMN IF EXISTS "references";
