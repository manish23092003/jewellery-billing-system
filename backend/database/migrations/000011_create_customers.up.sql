CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255) DEFAULT '',
    address TEXT DEFAULT '',
    total_purchases DECIMAL(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_customers_organization_id ON customers(organization_id);
CREATE INDEX idx_customers_phone ON customers(phone);

-- Add customer_id to bills
ALTER TABLE bills ADD COLUMN customer_id UUID REFERENCES customers(id) ON DELETE SET NULL;
