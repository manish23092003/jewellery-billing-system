package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"jewellery-billing/internal/domain"
)

// MetalRateService handles metal rate CRUD operations.
type MetalRateService struct {
	metalRateRepo domain.MetalRateRepository
}

func NewMetalRateService(repo domain.MetalRateRepository) *MetalRateService {
	return &MetalRateService{metalRateRepo: repo}
}

func (s *MetalRateService) Create(ctx context.Context, orgID uuid.UUID, req domain.CreateMetalRateRequest) (*domain.MetalRate, error) {
	if !req.MetalType.IsValid() {
		return nil, fmt.Errorf("invalid metal type: %s", req.MetalType)
	}
	if !domain.IsValidPurity(req.MetalType, req.Purity) {
		return nil, fmt.Errorf("invalid purity '%s' for metal type '%s'", req.Purity, req.MetalType)
	}
	if req.RatePerGram <= 0 {
		return nil, fmt.Errorf("rate per gram must be positive")
	}
	if req.EffectiveDate == "" {
		return nil, fmt.Errorf("effective date is required")
	}

	rate := &domain.MetalRate{
		OrganizationID: orgID,
		MetalType:      req.MetalType,
		Purity:         req.Purity,
		RatePerGram:    req.RatePerGram,
		EffectiveDate:  req.EffectiveDate,
	}

	if err := s.metalRateRepo.Create(ctx, rate); err != nil {
		return nil, fmt.Errorf("failed to create metal rate: %w", err)
	}

	return rate, nil
}

func (s *MetalRateService) GetCurrentRates(ctx context.Context, orgID uuid.UUID) ([]domain.MetalRate, error) {
	return s.metalRateRepo.GetCurrentRates(ctx, orgID)
}

func (s *MetalRateService) GetHistory(ctx context.Context, orgID uuid.UUID, metalType string, limit, offset int) ([]domain.MetalRate, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.metalRateRepo.GetHistory(ctx, orgID, metalType, limit, offset)
}

func (s *MetalRateService) Delete(ctx context.Context, orgID, id uuid.UUID) error {
	return s.metalRateRepo.Delete(ctx, orgID, id)
}
