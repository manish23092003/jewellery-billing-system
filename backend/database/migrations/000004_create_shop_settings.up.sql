CREATE TABLE IF NOT EXISTS shop_settings (
    id SERIAL PRIMARY KEY,
    shop_name VARCHAR(255) NOT NULL DEFAULT 'Aura Jewels',
    gstin VARCHAR(50) NOT NULL DEFAULT '',
    phone VARCHAR(20) NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    logo_path VARCHAR(500) NOT NULL DEFAULT '',
    invoice_prefix VARCHAR(20) NOT NULL DEFAULT 'INV',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Ensure only one row exists using a check constraint
ALTER TABLE shop_settings ADD CONSTRAINT singleton_check CHECK (id = 1);

-- Insert the default singleton row
INSERT INTO shop_settings (id) VALUES (1) ON CONFLICT DO NOTHING;
