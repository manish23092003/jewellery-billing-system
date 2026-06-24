package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"jewellery-billing/internal/domain"
)

// SettingService handles shop settings CRUD operations.
type SettingService struct {
	settingRepo domain.SettingRepository
}

func NewSettingService(settingRepo domain.SettingRepository) *SettingService {
	return &SettingService{settingRepo: settingRepo}
}

func (s *SettingService) Get(ctx context.Context, orgID uuid.UUID) (*domain.ShopSettings, error) {
	return s.settingRepo.Get(ctx, orgID)
}

func (s *SettingService) Update(ctx context.Context, orgID uuid.UUID, req domain.UpdateShopSettingsRequest) (*domain.ShopSettings, error) {
	settings := &domain.ShopSettings{
		OrganizationID: orgID,
		ShopName:       req.ShopName,
		GSTIN:          req.GSTIN,
		Phone:          req.Phone,
		Address:        req.Address,
		InvoicePrefix:  req.InvoicePrefix,
	}

	if err := s.settingRepo.Upsert(ctx, settings); err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err)
	}

	// Re-fetch to get the full record
	return s.settingRepo.Get(ctx, orgID)
}

func (s *SettingService) UpdateLogo(ctx context.Context, orgID uuid.UUID, logoPath string) error {
	return s.settingRepo.UpdateLogo(ctx, orgID, logoPath)
}

func (s *SettingService) GetInvoicePrefix(ctx context.Context, orgID uuid.UUID) string {
	settings, err := s.settingRepo.Get(ctx, orgID)
	if err != nil || settings.InvoicePrefix == "" {
		return "INV"
	}
	return settings.InvoicePrefix
}
