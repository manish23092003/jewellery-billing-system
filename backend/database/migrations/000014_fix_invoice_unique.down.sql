-- ============================================
-- Migration 000014: Fix Invoice Unique Constraint
-- ============================================

ALTER TABLE bills DROP CONSTRAINT IF EXISTS bills_org_invoice_unique;
ALTER TABLE bills ADD CONSTRAINT bills_invoice_number_key UNIQUE (invoice_number);
