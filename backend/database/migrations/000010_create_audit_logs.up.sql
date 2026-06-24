-- ============================================
-- Migration 000010: Create audit logs
-- ============================================

CREATE TABLE audit_logs (
    id              UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID          NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id         UUID          REFERENCES users(id) ON DELETE SET NULL,
    action          VARCHAR(100)  NOT NULL,
    entity_type     VARCHAR(50)   NOT NULL,
    entity_id       VARCHAR(100),
    details         JSONB,
    ip_address      VARCHAR(45),
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_org       ON audit_logs(organization_id);
CREATE INDEX idx_audit_time      ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_entity    ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_user      ON audit_logs(user_id);
