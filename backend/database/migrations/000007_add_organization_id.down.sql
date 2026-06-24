-- Reverse organization_id additions

-- shop_settings
DROP INDEX IF EXISTS idx_shop_settings_organization;
ALTER TABLE shop_settings DROP CONSTRAINT IF EXISTS fk_shop_settings_organization;
ALTER TABLE shop_settings DROP COLUMN IF EXISTS organization_id;
ALTER TABLE shop_settings ADD CONSTRAINT singleton_check CHECK (id = 1);

-- metal_rates
DROP INDEX IF EXISTS idx_metal_rates_unique;
CREATE UNIQUE INDEX idx_metal_rates_unique ON metal_rates(metal_type, purity, effective_date);
DROP INDEX IF EXISTS idx_metal_rates_organization;
ALTER TABLE metal_rates DROP CONSTRAINT IF EXISTS fk_metal_rates_organization;
ALTER TABLE metal_rates DROP COLUMN IF EXISTS organization_id;

-- expenses
DROP INDEX IF EXISTS idx_expenses_organization;
ALTER TABLE expenses DROP CONSTRAINT IF EXISTS fk_expenses_organization;
ALTER TABLE expenses DROP COLUMN IF EXISTS organization_id;

-- bill_items
DROP INDEX IF EXISTS idx_bill_items_organization;
ALTER TABLE bill_items DROP CONSTRAINT IF EXISTS fk_bill_items_organization;
ALTER TABLE bill_items DROP COLUMN IF EXISTS organization_id;

-- bills
DROP INDEX IF EXISTS idx_bills_organization;
ALTER TABLE bills DROP CONSTRAINT IF EXISTS fk_bills_organization;
ALTER TABLE bills DROP COLUMN IF EXISTS organization_id;

-- users
DROP INDEX IF EXISTS idx_users_organization;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_organization;
ALTER TABLE users DROP COLUMN IF EXISTS organization_id;
ALTER TABLE users DROP COLUMN IF EXISTS is_active;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified;

-- Remove default organization
DELETE FROM organizations WHERE id = '00000000-0000-0000-0000-000000000001';
