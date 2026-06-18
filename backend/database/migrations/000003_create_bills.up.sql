-- ============================================
-- Migration 000003: Create bills + bill_items tables
-- ============================================
-- Core billing tables. Items are entered inline — no inventory reference.
-- All monetary calculations are stored for audit trail.

CREATE TABLE bills (
    id              UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number  VARCHAR(20)   UNIQUE NOT NULL,
    invoice_date    DATE          NOT NULL DEFAULT CURRENT_DATE,
    customer_name   VARCHAR(100),
    customer_phone  VARCHAR(20),
    subtotal        DECIMAL(14,2) NOT NULL DEFAULT 0,
    gst_amount      DECIMAL(14,2) NOT NULL DEFAULT 0,
    grand_total     DECIMAL(14,2) NOT NULL DEFAULT 0,
    payment_method  VARCHAR(20)   NOT NULL DEFAULT 'cash',
    notes           TEXT,
    created_by      UUID          REFERENCES users(id),
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_payment_method CHECK (
        payment_method IN ('cash', 'card', 'upi', 'bank_transfer')
    )
);

CREATE TABLE bill_items (
    id              UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    bill_id         UUID          NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
    item_name       VARCHAR(100)  NOT NULL,
    metal_type      VARCHAR(10)   NOT NULL,
    purity          VARCHAR(10)   NOT NULL,
    weight          DECIMAL(10,3) NOT NULL,
    rate_per_gram   DECIMAL(12,2) NOT NULL,
    making_charge   DECIMAL(12,2) NOT NULL DEFAULT 0,
    gst_percentage  DECIMAL(5,2)  NOT NULL DEFAULT 3.00,
    quantity        INT           NOT NULL DEFAULT 1,
    line_total      DECIMAL(14,2) NOT NULL,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- Indexes for common query patterns
CREATE INDEX idx_bills_invoice_number ON bills(invoice_number);
CREATE INDEX idx_bills_invoice_date   ON bills(invoice_date DESC);
CREATE INDEX idx_bills_created_by     ON bills(created_by);
CREATE INDEX idx_bills_payment_method ON bills(payment_method);
CREATE INDEX idx_bill_items_bill_id   ON bill_items(bill_id);
