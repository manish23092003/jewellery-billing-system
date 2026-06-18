package repository

import (
	"context"

	"jewellery-billing/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SettingRepo implements domain.SettingRepository using pgx.
type SettingRepo struct {
	db *pgxpool.Pool
}

func NewSettingRepository(db *pgxpool.Pool) *SettingRepo {
	return &SettingRepo{db: db}
}

func (r *SettingRepo) Get(ctx context.Context) (*domain.ShopSettings, error) {
	query := `
		SELECT id, shop_name, gstin, phone, address, logo_path, invoice_prefix, TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM shop_settings
		WHERE id = 1
	`
	var s domain.ShopSettings
	err := r.db.QueryRow(ctx, query).Scan(
		&s.ID, &s.ShopName, &s.GSTIN, &s.Phone, &s.Address, &s.LogoPath, &s.InvoicePrefix, &s.UpdatedAt,
	)
	if err != nil {
		// If no row exists, return a default empty settings object so frontend doesn't crash
		if err.Error() == "no rows in result set" {
			return &domain.ShopSettings{
				ID:            1,
				ShopName:      "My Jewellery Shop",
				GSTIN:         "",
				Phone:         "",
				Address:       "",
				LogoPath:      "",
				InvoicePrefix: "INV",
			}, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *SettingRepo) Update(ctx context.Context, settings *domain.ShopSettings) error {
	// UPSERT (Insert if not exists, else Update)
	query := `
		INSERT INTO shop_settings (id, shop_name, gstin, phone, address, invoice_prefix, updated_at)
		VALUES (1, $1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE 
		SET shop_name = EXCLUDED.shop_name, 
		    gstin = EXCLUDED.gstin, 
		    phone = EXCLUDED.phone, 
		    address = EXCLUDED.address, 
		    invoice_prefix = EXCLUDED.invoice_prefix, 
		    updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(ctx, query, settings.ShopName, settings.GSTIN, settings.Phone, settings.Address, settings.InvoicePrefix)
	return err
}

func (r *SettingRepo) UpdateLogo(ctx context.Context, logoPath string) error {
	query := `
		UPDATE shop_settings
		SET logo_path = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = 1
	`
	_, err := r.db.Exec(ctx, query, logoPath)
	return err
}
