package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ── Bill Entity ────────────────────────────────────────────────────────

type Bill struct {
	ID            uuid.UUID  `json:"id"`
	InvoiceNumber string     `json:"invoice_number"`
	InvoiceDate   string     `json:"invoice_date"` // YYYY-MM-DD
	CustomerName  string     `json:"customer_name,omitempty"`
	CustomerPhone string     `json:"customer_phone,omitempty"`
	Subtotal      float64    `json:"subtotal"`
	GSTAmount     float64    `json:"gst_amount"`
	GrandTotal    float64    `json:"grand_total"`
	PaymentMethod string     `json:"payment_method"`
	Notes         string     `json:"notes,omitempty"`
	CreatedBy     uuid.UUID  `json:"created_by"`
	CreatedByName string     `json:"created_by_name,omitempty"` // joined from users
	CreatedAt     time.Time  `json:"created_at"`
	Items         []BillItem `json:"items,omitempty"`
}

// ── BillItem Entity ────────────────────────────────────────────────────

type BillItem struct {
	ID            uuid.UUID `json:"id"`
	BillID        uuid.UUID `json:"bill_id"`
	ItemName      string    `json:"item_name"`
	MetalType     string    `json:"metal_type"`
	Purity        string    `json:"purity"`
	Weight        float64   `json:"weight"`
	RatePerGram   float64   `json:"rate_per_gram"`
	MakingCharge  float64   `json:"making_charge"`
	GSTPercentage float64   `json:"gst_percentage"`
	Quantity      int       `json:"quantity"`
	LineTotal     float64   `json:"line_total"`
	CreatedAt     time.Time `json:"created_at"`
}

// ── Request DTOs ───────────────────────────────────────────────────────

type CreateBillRequest struct {
	InvoiceDate   string              `json:"invoice_date"`
	CustomerName  string              `json:"customer_name"`
	CustomerPhone string              `json:"customer_phone"`
	PaymentMethod string              `json:"payment_method"`
	Notes         string              `json:"notes"`
	Items         []CreateBillItemDTO `json:"items"`
}

type CreateBillItemDTO struct {
	ItemName      string  `json:"item_name"`
	MetalType     string  `json:"metal_type"`
	Purity        string  `json:"purity"`
	Weight        float64 `json:"weight"`
	RatePerGram   float64 `json:"rate_per_gram"`
	MakingCharge  float64 `json:"making_charge"`
	GSTPercentage float64 `json:"gst_percentage"`
	Quantity      int     `json:"quantity"`
}

// ── Filter DTO ─────────────────────────────────────────────────────────

type BillFilter struct {
	Search        string // search invoice_number or customer_name
	PaymentMethod string
	DateFrom      string // YYYY-MM-DD
	DateTo        string // YYYY-MM-DD
	Page          int
	PerPage       int
}

// ── Payment Methods ────────────────────────────────────────────────────

var ValidPaymentMethods = map[string]bool{
	"cash":          true,
	"card":          true,
	"upi":           true,
	"bank_transfer": true,
}

// ── Repository Interface ───────────────────────────────────────────────

type BillRepository interface {
	Create(ctx context.Context, bill *Bill) error
	GetByID(ctx context.Context, id uuid.UUID) (*Bill, error)
	GetAll(ctx context.Context, filter BillFilter) ([]Bill, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetNextInvoiceNumber(ctx context.Context, prefix string) (string, error)
}
