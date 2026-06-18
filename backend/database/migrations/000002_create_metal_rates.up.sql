-- ============================================
-- Migration 000002: Create metal_rates table
-- ============================================
-- Stores daily rates for gold (24K, 22K, 18K) and silver.
-- One rate per metal_type + purity + effective_date.

CREATE TABLE metal_rates (
    id             UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    metal_type     VARCHAR(10)   NOT NULL,
    purity         VARCHAR(10)   NOT NULL,
    rate_per_gram  DECIMAL(12,2) NOT NULL,
    effective_date DATE          NOT NULL,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_metal_type CHECK (metal_type IN ('gold', 'silver')),
    CONSTRAINT chk_purity     CHECK (purity IN ('24K', '22K', '18K', 'pure'))
);

-- Fast lookups for "current rates" query
CREATE INDEX idx_metal_rates_effective_date ON metal_rates(effective_date DESC);

-- Prevent duplicate rates for the same metal/purity/date combination
CREATE UNIQUE INDEX idx_metal_rates_unique ON metal_rates(metal_type, purity, effective_date);

-- Seed today's sample rates
INSERT INTO metal_rates (metal_type, purity, rate_per_gram, effective_date) VALUES
    ('gold',   '24K',  7250.00, CURRENT_DATE),
    ('gold',   '22K',  6650.00, CURRENT_DATE),
    ('gold',   '18K',  5440.00, CURRENT_DATE),
    ('silver', 'pure',   92.50, CURRENT_DATE);
