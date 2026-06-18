-- ============================================
-- Migration 000001: Create users table
-- ============================================
-- Supports two roles: admin (full access) and staff (limited access).
-- Passwords are stored as bcrypt hashes — never plaintext.

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role AS ENUM ('admin', 'staff');

CREATE TABLE users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(100)  NOT NULL,
    email         VARCHAR(255)  UNIQUE NOT NULL,
    password_hash VARCHAR(255)  NOT NULL,
    role          user_role     NOT NULL DEFAULT 'staff',
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- Fast lookup by email during login
CREATE INDEX idx_users_email ON users(email);

-- Filter users by role on the admin panel
CREATE INDEX idx_users_role ON users(role);
