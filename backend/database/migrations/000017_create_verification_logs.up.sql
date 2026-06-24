-- ============================================
-- Migration 000017: Create verification_logs table
-- ============================================

CREATE TABLE IF NOT EXISTS verification_logs (
    id             UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    token          VARCHAR(255)  NOT NULL,
    ip_address     VARCHAR(45),
    user_agent     TEXT,
    is_valid       BOOLEAN       NOT NULL,
    failure_reason VARCHAR(255),
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_verification_logs_token ON verification_logs(token);
