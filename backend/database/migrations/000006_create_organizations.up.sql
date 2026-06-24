-- ============================================
-- Migration 000006: Create organizations table
-- ============================================
-- Each jewellery shop registers as an organization. All business data
-- (users, bills, expenses, rates, settings) belongs to one organization.

CREATE TABLE organizations (
    id                  UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    business_name       VARCHAR(255)  NOT NULL,
    owner_name          VARCHAR(100)  NOT NULL,
    email               VARCHAR(255)  UNIQUE NOT NULL,
    phone               VARCHAR(20)   NOT NULL DEFAULT '',
    gstin               VARCHAR(50)   NOT NULL DEFAULT '',
    address             TEXT          NOT NULL DEFAULT '',
    subscription_status VARCHAR(20)   NOT NULL DEFAULT 'active',
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_subscription_status CHECK (
        subscription_status IN ('active', 'inactive', 'trial', 'suspended')
    )
);

-- Fast lookup by email during registration uniqueness check
CREATE INDEX idx_organizations_email ON organizations(email);

-- Subscription status filtering for admin dashboards
CREATE INDEX idx_organizations_status ON organizations(subscription_status);
