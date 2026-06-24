package repository

import (
	"context"

	"jewellery-billing/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SettingRepo implements domain.SettingRepository using pgx.
type SettingRepo struct {
	db *pgxpool.Pool
}

func NewSettingRepository(db *pgxpool.Pool) *SettingRepo {
	return &SettingRepo{db: db}
}

func (r *SettingRepo) Get(ctx context.Context, orgID uuid.UUID) (*domain.ShopSettings, error) {
	query := `
		SELECT id, organization_id, shop_name, gstin, phone, address, logo_path, invoice_prefix,
		       TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM shop_settings
		WHERE organization_id = $1
	`
	var s domain.ShopSettings
	err := r.db.QueryRow(ctx, query, orgID).Scan(
		&s.ID, &s.OrganizationID, &s.ShopName, &s.GSTIN, &s.Phone,
		&s.Address, &s.LogoPath, &s.InvoicePrefix, &s.UpdatedAt,
	)
	if err != nil {
		// If no row exists, return a default empty settings object so frontend doesn't crash
		if err.Error() == "no rows in result set" {
			return &domain.ShopSettings{
				OrganizationID: orgID,
				ShopName:       "My Jewellery Shop",
				GSTIN:          "",
				Phone:          "",
				Address:        "",
				LogoPath:       "",
				InvoicePrefix:  "INV",
			}, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *SettingRepo) Upsert(ctx context.Context, settings *domain.ShopSettings) error {
	query := `
		INSERT INTO shop_settings (organization_id, shop_name, gstin, phone, address, invoice_prefix, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		ON CONFLICT (organization_id) DO UPDATE
		SET shop_name = EXCLUDED.shop_name,
		    gstin = EXCLUDED.gstin,
		    phone = EXCLUDED.phone,
		    address = EXCLUDED.address,
		    invoice_prefix = EXCLUDED.invoice_prefix,
		    updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(ctx, query,
		settings.OrganizationID, settings.ShopName, settings.GSTIN,
		settings.Phone, settings.Address, settings.InvoicePrefix,
	)
	return err
}

func (r *SettingRepo) UpdateLogo(ctx context.Context, orgID uuid.UUID, logoPath string) error {
	query := `
		UPDATE shop_settings
		SET logo_path = $1, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $2
	`
	_, err := r.db.Exec(ctx, query, logoPath, orgID)
	return err
}
