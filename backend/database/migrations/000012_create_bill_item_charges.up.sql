CREATE TABLE IF NOT EXISTS bill_item_charges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bill_item_id UUID NOT NULL REFERENCES bill_items(id) ON DELETE CASCADE,
    charge_name VARCHAR(100) NOT NULL,
    amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bill_item_charges_item_id ON bill_item_charges(bill_item_id);
