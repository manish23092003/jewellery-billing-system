package domain

import "context"

// ShopSettings represents the core store configuration.
type ShopSettings struct {
	ID            int    `json:"id" db:"id"`
	ShopName      string `json:"shop_name" db:"shop_name"`
	GSTIN         string `json:"gstin" db:"gstin"`
	Phone         string `json:"phone" db:"phone"`
	Address       string `json:"address" db:"address"`
	LogoPath      string `json:"logo_path" db:"logo_path"`
	InvoicePrefix string `json:"invoice_prefix" db:"invoice_prefix"`
	UpdatedAt     string `json:"updated_at" db:"updated_at"`
}

// UpdateShopSettingsRequest is the payload for updating the settings.
type UpdateShopSettingsRequest struct {
	ShopName      string `json:"shop_name"`
	GSTIN         string `json:"gstin"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	InvoicePrefix string `json:"invoice_prefix"`
	// Logo is handled separately via multipart form
}

// SettingRepository handles persistent storage operations for settings.
type SettingRepository interface {
	Get(ctx context.Context) (*ShopSettings, error)
	Update(ctx context.Context, settings *ShopSettings) error
	UpdateLogo(ctx context.Context, logoPath string) error
}
