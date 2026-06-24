ALTER TABLE bills ADD COLUMN type VARCHAR(20) NOT NULL DEFAULT 'invoice';
ALTER TABLE bills ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'completed';
ALTER TABLE bills ADD COLUMN old_gold_amount DECIMAL(14,2) NOT NULL DEFAULT 0;
ALTER TABLE bills ADD COLUMN advance_amount DECIMAL(14,2) NOT NULL DEFAULT 0;
ALTER TABLE bills ADD COLUMN balance_due DECIMAL(14,2) NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS bill_old_gold (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bill_id UUID NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    weight DECIMAL(10,3) NOT NULL,
    purity VARCHAR(10) NOT NULL,
    melting_loss_percentage DECIMAL(5,2) NOT NULL DEFAULT 0,
    rate_per_gram DECIMAL(12,2) NOT NULL,
    total_value DECIMAL(14,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bill_old_gold_bill_id ON bill_old_gold(bill_id);
