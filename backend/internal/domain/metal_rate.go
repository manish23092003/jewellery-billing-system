package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ── Metal Types ────────────────────────────────────────────────────────

type MetalType string

const (
	MetalGold   MetalType = "gold"
	MetalSilver MetalType = "silver"
)

func (m MetalType) IsValid() bool {
	return m == MetalGold || m == MetalSilver
}

// ValidPurities returns the allowed purities for a given metal type.
func ValidPurities(mt MetalType) []string {
	switch mt {
	case MetalGold:
		return []string{"24K", "22K", "18K"}
	case MetalSilver:
		return []string{"pure"}
	default:
		return nil
	}
}

func IsValidPurity(mt MetalType, purity string) bool {
	for _, p := range ValidPurities(mt) {
		if p == purity {
			return true
		}
	}
	return false
}

// ── MetalRate Entity ───────────────────────────────────────────────────

type MetalRate struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	MetalType      MetalType `json:"metal_type"`
	Purity         string    `json:"purity"`
	RatePerGram    float64   `json:"rate_per_gram"`
	EffectiveDate  string    `json:"effective_date"` // YYYY-MM-DD
	CreatedAt      time.Time `json:"created_at"`
}

// ── Request DTOs ───────────────────────────────────────────────────────

type CreateMetalRateRequest struct {
	MetalType     MetalType `json:"metal_type"`
	Purity        string    `json:"purity"`
	RatePerGram   float64   `json:"rate_per_gram"`
	EffectiveDate string    `json:"effective_date"`
}

type UpdateMetalRateRequest struct {
	RatePerGram   float64 `json:"rate_per_gram"`
	EffectiveDate string  `json:"effective_date,omitempty"`
}

// ── Repository Interface ───────────────────────────────────────────────

type MetalRateRepository interface {
	Create(ctx context.Context, rate *MetalRate) error
	GetByID(ctx context.Context, orgID, id uuid.UUID) (*MetalRate, error)
	Update(ctx context.Context, rate *MetalRate) error
	Delete(ctx context.Context, orgID, id uuid.UUID) error

	// GetCurrentRates returns the latest rate for each metal_type + purity
	// combination where effective_date <= today, scoped to the organization.
	GetCurrentRates(ctx context.Context, orgID uuid.UUID) ([]MetalRate, error)

	// GetHistory returns paginated rate history, optionally filtered by metal type.
	GetHistory(ctx context.Context, orgID uuid.UUID, metalType string, limit, offset int) ([]MetalRate, int64, error)
}
