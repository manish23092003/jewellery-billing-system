-- ============================================
-- Migration 000008: Create email verification tokens
-- ============================================

CREATE TABLE email_verification_tokens (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_evt_token ON email_verification_tokens(token);
CREATE INDEX idx_evt_user  ON email_verification_tokens(user_id);
