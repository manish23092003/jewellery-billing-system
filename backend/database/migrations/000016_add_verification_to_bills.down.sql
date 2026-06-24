-- ============================================
-- Migration 000016: Add verification to bills
-- ============================================

DROP INDEX IF EXISTS idx_bills_invoice_hash;
DROP INDEX IF EXISTS idx_bills_verification_token;

ALTER TABLE bills
DROP COLUMN verification_status,
DROP COLUMN invoice_hash,
DROP COLUMN verification_token;
