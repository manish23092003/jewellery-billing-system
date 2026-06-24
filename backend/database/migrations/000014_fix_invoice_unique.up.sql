-- ============================================
-- Migration 000014: Fix Invoice Unique Constraint
-- ============================================

-- Drop the global unique constraint on invoice_number
ALTER TABLE bills DROP CONSTRAINT IF EXISTS bills_invoice_number_key;

-- Create a unique constraint that is organization-scoped
ALTER TABLE bills ADD CONSTRAINT bills_org_invoice_unique UNIQUE (organization_id, invoice_number);
