-- ============================================
-- Migration 000015: Create bill_payments table
-- ============================================

CREATE TABLE IF NOT EXISTS bill_payments (
    id              UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID          NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    bill_id         UUID          NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
    amount          DECIMAL(14,2) NOT NULL,
    payment_date    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bill_payments_bill_id ON bill_payments(bill_id);
