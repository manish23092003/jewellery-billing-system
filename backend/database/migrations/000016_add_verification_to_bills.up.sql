-- ============================================
-- Migration 000016: Add verification to bills
-- ============================================

ALTER TABLE bills
ADD COLUMN verification_token UUID UNIQUE DEFAULT gen_random_uuid(),
ADD COLUMN invoice_hash VARCHAR(64) UNIQUE,
ADD COLUMN verification_status VARCHAR(20) DEFAULT 'ACTIVE';

CREATE INDEX idx_bills_verification_token ON bills(verification_token);
CREATE INDEX idx_bills_invoice_hash ON bills(invoice_hash);
