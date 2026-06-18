package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"jewellery-billing/internal/domain"
)

// MetalRateService handles business logic for metal rate management.
type MetalRateService struct {
	rateRepo domain.MetalRateRepository
}

func NewMetalRateService(rateRepo domain.MetalRateRepository) *MetalRateService {
	return &MetalRateService{rateRepo: rateRepo}
}

// Create validates and persists a new metal rate.
func (s *MetalRateService) Create(ctx context.Context, req domain.CreateMetalRateRequest) (*domain.MetalRate, error) {
	// Normalize inputs
	mt := domain.MetalType(strings.ToLower(strings.TrimSpace(string(req.MetalType))))
	req.MetalType = mt

	p := strings.ToUpper(strings.TrimSpace(req.Purity))
	if mt == domain.MetalGold {
		if p == "24" || p == "22" || p == "18" {
			p += "K"
		}
	} else if mt == domain.MetalSilver {
		if p == "PURE" {
			p = "pure"
		}
	}
	req.Purity = p

	if !req.MetalType.IsValid() {
		return nil, fmt.Errorf("invalid metal type — must be 'gold' or 'silver'")
	}
	if !domain.IsValidPurity(req.MetalType, req.Purity) {
		return nil, fmt.Errorf("invalid purity '%s' for metal type '%s'", req.Purity, req.MetalType)
	}
	if req.RatePerGram <= 0 {
		return nil, fmt.Errorf("rate per gram must be positive")
	}
	if req.EffectiveDate == "" {
		req.EffectiveDate = "now()" // Postgres fallback
	}

	rate := &domain.MetalRate{
		MetalType:     req.MetalType,
		Purity:        req.Purity,
		RatePerGram:   req.RatePerGram,
		EffectiveDate: req.EffectiveDate,
	}

	if err := s.rateRepo.Create(ctx, rate); err != nil {
		return nil, fmt.Errorf("failed to create rate: %w", err)
	}
	return rate, nil
}

// GetCurrentRates returns the latest rate for each metal/purity combination.
func (s *MetalRateService) GetCurrentRates(ctx context.Context) ([]domain.MetalRate, error) {
	return s.rateRepo.GetCurrentRates(ctx)
}

// GetHistory returns paginated rate history with optional metal type filter.
func (s *MetalRateService) GetHistory(ctx context.Context, metalType string, page, perPage int) ([]domain.MetalRate, int64, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	return s.rateRepo.GetHistory(ctx, metalType, perPage, offset)
}

// Update modifies an existing rate's value and/or effective date.
func (s *MetalRateService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateMetalRateRequest) (*domain.MetalRate, error) {
	rate, err := s.rateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.RatePerGram > 0 {
		rate.RatePerGram = req.RatePerGram
	}
	if req.EffectiveDate != "" {
		rate.EffectiveDate = req.EffectiveDate
	}

	if err := s.rateRepo.Update(ctx, rate); err != nil {
		return nil, err
	}
	return rate, nil
}

// Delete removes a metal rate entry.
func (s *MetalRateService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.rateRepo.Delete(ctx, id)
}
