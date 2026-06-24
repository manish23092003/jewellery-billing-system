-- ============================================
-- Migration 000007: Add organization_id to all data tables
-- ============================================
-- Retrofits multi-tenancy onto existing tables. Creates a "Default Organization"
-- and assigns all existing rows to it for backward compatibility.

-- Step 1: Create a default organization for existing data
INSERT INTO organizations (id, business_name, owner_name, email)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Default Organization',
    'Administrator',
    'admin@jewellery.com'
);

-- Step 2: Add columns to users
ALTER TABLE users
    ADD COLUMN organization_id UUID,
    ADD COLUMN is_active       BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN email_verified  BOOLEAN NOT NULL DEFAULT false;

-- Backfill existing users
UPDATE users SET organization_id = '00000000-0000-0000-0000-000000000001';

-- Make NOT NULL + add FK
ALTER TABLE users
    ALTER COLUMN organization_id SET NOT NULL,
    ADD CONSTRAINT fk_users_organization
        FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

CREATE INDEX idx_users_organization ON users(organization_id);

-- Step 3: Add organization_id to bills
ALTER TABLE bills ADD COLUMN organization_id UUID;
UPDATE bills SET organization_id = '00000000-0000-0000-0000-000000000001';
ALTER TABLE bills
    ALTER COLUMN organization_id SET NOT NULL,
    ADD CONSTRAINT fk_bills_organization
        FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_bills_organization ON bills(organization_id);

-- Step 4: Add organization_id to bill_items
ALTER TABLE bill_items ADD COLUMN organization_id UUID;
UPDATE bill_items SET organization_id = '00000000-0000-0000-0000-000000000001';
ALTER TABLE bill_items
    ALTER COLUMN organization_id SET NOT NULL,
    ADD CONSTRAINT fk_bill_items_organization
        FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_bill_items_organization ON bill_items(organization_id);

-- Step 5: Add organization_id to expenses
ALTER TABLE expenses ADD COLUMN organization_id UUID;
UPDATE expenses SET organization_id = '00000000-0000-0000-0000-000000000001';
ALTER TABLE expenses
    ALTER COLUMN organization_id SET NOT NULL,
    ADD CONSTRAINT fk_expenses_organization
        FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_expenses_organization ON expenses(organization_id);

-- Step 6: Add organization_id to metal_rates
ALTER TABLE metal_rates ADD COLUMN organization_id UUID;
UPDATE metal_rates SET organization_id = '00000000-0000-0000-0000-000000000001';
ALTER TABLE metal_rates
    ALTER COLUMN organization_id SET NOT NULL,
    ADD CONSTRAINT fk_metal_rates_organization
        FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;
CREATE INDEX idx_metal_rates_organization ON metal_rates(organization_id);

-- Update unique constraint to be per-organization
DROP INDEX IF EXISTS idx_metal_rates_unique;
CREATE UNIQUE INDEX idx_metal_rates_unique ON metal_rates(organization_id, metal_type, purity, effective_date);

-- Step 7: Rework shop_settings for multi-tenancy
-- Remove the singleton constraint
ALTER TABLE shop_settings DROP CONSTRAINT IF EXISTS singleton_check;

-- Add organization_id
ALTER TABLE shop_settings ADD COLUMN organization_id UUID;
UPDATE shop_settings SET organization_id = '00000000-0000-0000-0000-000000000001';
ALTER TABLE shop_settings
    ALTER COLUMN organization_id SET NOT NULL,
    ADD CONSTRAINT fk_shop_settings_organization
        FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- Replace singleton with per-org uniqueness
CREATE UNIQUE INDEX idx_shop_settings_organization ON shop_settings(organization_id);
