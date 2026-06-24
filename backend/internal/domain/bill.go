package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ── Bill Entity ────────────────────────────────────────────────────────

type Bill struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	InvoiceNumber  string     `json:"invoice_number"`
	InvoiceDate    string          `json:"invoice_date"` // YYYY-MM-DD
	Type           string          `json:"type"`   // 'invoice' or 'estimate'
	Status         string          `json:"status"` // 'pending', 'completed', 'cancelled'
	CustomerName   string          `json:"customer_name,omitempty"`
	CustomerPhone  string          `json:"customer_phone,omitempty"`
	Subtotal       float64         `json:"subtotal"`
	GSTAmount      float64         `json:"gst_amount"`
	OldGoldAmount  float64         `json:"old_gold_amount"`
	AdvanceAmount  float64         `json:"advance_amount"`
	GrandTotal     float64         `json:"grand_total"` // Before old gold & advance
	BalanceDue     float64         `json:"balance_due"` // Final amount payable
	PaymentMethod  string          `json:"payment_method"`
	Notes          string          `json:"notes,omitempty"`
	CreatedBy      uuid.UUID       `json:"created_by"`
	CreatedByName  string          `json:"created_by_name,omitempty"` // joined from users
	CreatedAt      time.Time       `json:"created_at"`
	Items          []BillItem      `json:"items,omitempty"`
	OldGoldItems       []BillOldGold   `json:"old_gold_items,omitempty"`
	Payments           []BillPayment   `json:"payments,omitempty"`
	VerificationToken  uuid.UUID       `json:"verification_token,omitempty"`
	InvoiceHash        string          `json:"invoice_hash,omitempty"`
	VerificationStatus string          `json:"verification_status,omitempty"`
}

// ── BillPayment Entity ──────────────────────────────────────────────────

type BillPayment struct {
	ID            uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	BillID        uuid.UUID `json:"bill_id"`
	Amount        float64   `json:"amount"`
	PaymentDate   time.Time `json:"payment_date"`
}

// ── BillItem Entity ────────────────────────────────────────────────────

type BillItem struct {
	ID            uuid.UUID `json:"id"`
	BillID        uuid.UUID `json:"bill_id"`
	ItemName      string    `json:"item_name"`
	HSNCode       string    `json:"hsn_code,omitempty"`
	MetalType     string    `json:"metal_type"`
	Purity        string    `json:"purity"`
	Weight        float64   `json:"weight"`
	RatePerGram   float64   `json:"rate_per_gram"`
	MakingCharge  float64   `json:"making_charge"`
	GSTPercentage float64   `json:"gst_percentage"`
	Quantity      int              `json:"quantity"`
	LineTotal     float64          `json:"line_total"`
	CreatedAt     time.Time        `json:"created_at"`
	Charges       []BillItemCharge `json:"charges,omitempty"`
}

// ── BillItemCharge Entity ──────────────────────────────────────────────

type BillItemCharge struct {
	ID         uuid.UUID `json:"id"`
	BillItemID uuid.UUID `json:"bill_item_id"`
	ChargeName string    `json:"charge_name"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

// ── BillOldGold Entity ──────────────────────────────────────────────────

type BillOldGold struct {
	ID                    uuid.UUID `json:"id"`
	BillID                uuid.UUID `json:"bill_id"`
	Name                  string    `json:"name"`
	Weight                float64   `json:"weight"`
	Purity                string    `json:"purity"`
	MeltingLossPercentage float64   `json:"melting_loss_percentage"`
	RatePerGram           float64   `json:"rate_per_gram"`
	TotalValue            float64   `json:"total_value"`
	CreatedAt             time.Time `json:"created_at"`
}

// ── VerificationLog Entity ──────────────────────────────────────────────

type VerificationLog struct {
	ID            uuid.UUID `json:"id"`
	Token         string    `json:"token"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	IsValid       bool      `json:"is_valid"`
	FailureReason string    `json:"failure_reason,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// ── HTTP DTOs ──────────────────────────────────────────────────────────

type PublicVerificationResponse struct {
	ShopName           string  `json:"shop_name"`
	InvoiceNumber      string  `json:"invoice_number"`
	InvoiceDate        string  `json:"invoice_date"`
	GrandTotal         float64 `json:"grand_total"`
	BalanceDue         float64 `json:"balance_due"`
	VerificationStatus string  `json:"verification_status"`
}

type CreateBillRequest struct {
	InvoiceDate   string                   `json:"invoice_date"`
	Type          string                   `json:"type"`   // 'invoice' or 'estimate'
	Status        string                   `json:"status"` // 'pending', 'completed'
	AdvanceAmount float64                  `json:"advance_amount"`
	CustomerName  string                   `json:"customer_name"`
	CustomerPhone string                   `json:"customer_phone"`
	PaymentMethod string                   `json:"payment_method"`
	Notes         string                   `json:"notes"`
	Items         []CreateBillItemDTO      `json:"items"`
	OldGoldItems  []CreateBillOldGoldDTO   `json:"old_gold_items"`
	ConvertFromID string                   `json:"convert_from_id,omitempty"`
}

type AddPaymentRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type CreateBillItemDTO struct {
	ItemName      string  `json:"item_name"`
	HSNCode       string  `json:"hsn_code"`
	MetalType     string  `json:"metal_type"`
	Purity        string  `json:"purity"`
	Weight        float64 `json:"weight"`
	RatePerGram   float64 `json:"rate_per_gram"`
	MakingCharge  float64                 `json:"making_charge"`
	GSTPercentage float64                 `json:"gst_percentage"`
	Quantity      int                     `json:"quantity"`
	Charges       []CreateBillItemChargeDTO `json:"charges"`
}

type CreateBillItemChargeDTO struct {
	ChargeName string  `json:"charge_name"`
	Amount     float64 `json:"amount"`
}

type CreateBillOldGoldDTO struct {
	Name                  string  `json:"name"`
	Weight                float64 `json:"weight"`
	Purity                string  `json:"purity"`
	MeltingLossPercentage float64 `json:"melting_loss_percentage"`
	RatePerGram           float64 `json:"rate_per_gram"`
}

// ── Filter DTO ─────────────────────────────────────────────────────────

type BillFilter struct {
	Search        string // search invoice_number or customer_name
	Type          string // 'invoice', 'estimate'
	Status        string
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
	GetByID(ctx context.Context, orgID, id uuid.UUID) (*Bill, error)
	GetAll(ctx context.Context, orgID uuid.UUID, filter BillFilter) ([]Bill, int64, error)
	Delete(ctx context.Context, orgID, id uuid.UUID) error
	GetNextInvoiceNumber(ctx context.Context, orgID uuid.UUID, prefix string) (string, error)
	UpdatePayment(ctx context.Context, orgID, billID uuid.UUID, paymentAmount, advanceAmount, balanceDue float64, status string) error
	GetByVerificationToken(ctx context.Context, token string) (*Bill, error)
}
