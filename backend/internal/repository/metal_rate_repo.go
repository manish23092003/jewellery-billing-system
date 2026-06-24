package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"jewellery-billing/internal/domain"
)

type metalRateRepository struct {
	db *pgxpool.Pool
}

func NewMetalRateRepository(db *pgxpool.Pool) domain.MetalRateRepository {
	return &metalRateRepository{db: db}
}

// ── Queries ────────────────────────────────────────────────────────────

const (
	queryInsertRateMultiTenant = `
		INSERT INTO metal_rates (organization_id, metal_type, purity, rate_per_gram, effective_date)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (organization_id, metal_type, purity, effective_date)
		DO UPDATE SET rate_per_gram = EXCLUDED.rate_per_gram
		RETURNING id, created_at`

	queryGetRateByIDMultiTenant = `
		SELECT id, organization_id, metal_type, purity, rate_per_gram,
		       TO_CHAR(effective_date, 'YYYY-MM-DD') AS effective_date, created_at
		FROM metal_rates WHERE id = $1 AND organization_id = $2`

	queryCurrentRatesMultiTenant = `
		SELECT DISTINCT ON (metal_type, purity)
		       id, organization_id, metal_type, purity, rate_per_gram,
		       TO_CHAR(effective_date, 'YYYY-MM-DD') AS effective_date, created_at
		FROM metal_rates
		WHERE organization_id = $1 AND effective_date <= CURRENT_DATE
		ORDER BY metal_type, purity, effective_date DESC`

	queryUpdateRateMultiTenant = `
		UPDATE metal_rates
		SET rate_per_gram = $3, effective_date = $4
		WHERE id = $1 AND organization_id = $2
		RETURNING created_at`

	queryDeleteRateMultiTenant = `DELETE FROM metal_rates WHERE id = $1 AND organization_id = $2`
)

// ── Implementation ─────────────────────────────────────────────────────

func (r *metalRateRepository) Create(ctx context.Context, rate *domain.MetalRate) error {
	return r.db.QueryRow(ctx, queryInsertRateMultiTenant,
		rate.OrganizationID, rate.MetalType, rate.Purity, rate.RatePerGram, rate.EffectiveDate,
	).Scan(&rate.ID, &rate.CreatedAt)
}

func (r *metalRateRepository) GetByID(ctx context.Context, orgID, id uuid.UUID) (*domain.MetalRate, error) {
	rate := &domain.MetalRate{}
	err := r.db.QueryRow(ctx, queryGetRateByIDMultiTenant, id, orgID).Scan(
		&rate.ID, &rate.OrganizationID, &rate.MetalType, &rate.Purity, &rate.RatePerGram,
		&rate.EffectiveDate, &rate.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("metal rate not found")
		}
		return nil, fmt.Errorf("failed to get metal rate: %w", err)
	}
	return rate, nil
}

func (r *metalRateRepository) GetCurrentRates(ctx context.Context, orgID uuid.UUID) ([]domain.MetalRate, error) {
	rows, err := r.db.Query(ctx, queryCurrentRatesMultiTenant, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to query current rates: %w", err)
	}
	defer rows.Close()

	var rates []domain.MetalRate
	for rows.Next() {
		var rate domain.MetalRate
		if err := rows.Scan(
			&rate.ID, &rate.OrganizationID, &rate.MetalType, &rate.Purity, &rate.RatePerGram,
			&rate.EffectiveDate, &rate.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan rate: %w", err)
		}
		rates = append(rates, rate)
	}
	return rates, rows.Err()
}

func (r *metalRateRepository) GetHistory(ctx context.Context, orgID uuid.UUID, metalType string, limit, offset int) ([]domain.MetalRate, int64, error) {
	baseQuery := `SELECT id, organization_id, metal_type, purity, rate_per_gram,
	                     TO_CHAR(effective_date, 'YYYY-MM-DD') AS effective_date, created_at
	              FROM metal_rates`
	countQuery := `SELECT COUNT(*) FROM metal_rates`

	args := []interface{}{orgID}
	where := " WHERE organization_id = $1"
	argIdx := 2

	if metalType != "" {
		where += fmt.Sprintf(" AND metal_type = $%d", argIdx)
		args = append(args, metalType)
		argIdx++
	}

	// Count total
	var total int64
	if err := r.db.QueryRow(ctx, countQuery+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count rates: %w", err)
	}

	// Fetch page
	query := baseQuery + where + fmt.Sprintf(
		" ORDER BY effective_date DESC, metal_type, purity LIMIT $%d OFFSET $%d",
		argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query rate history: %w", err)
	}
	defer rows.Close()

	var rates []domain.MetalRate
	for rows.Next() {
		var rate domain.MetalRate
		if err := rows.Scan(
			&rate.ID, &rate.OrganizationID, &rate.MetalType, &rate.Purity, &rate.RatePerGram,
			&rate.EffectiveDate, &rate.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan rate: %w", err)
		}
		rates = append(rates, rate)
	}
	return rates, total, rows.Err()
}

func (r *metalRateRepository) Update(ctx context.Context, rate *domain.MetalRate) error {
	result, err := r.db.Exec(ctx, queryUpdateRateMultiTenant,
		rate.ID, rate.OrganizationID, rate.RatePerGram, rate.EffectiveDate,
	)
	if err != nil {
		return fmt.Errorf("failed to update rate: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("metal rate not found")
	}
	return nil
}

func (r *metalRateRepository) Delete(ctx context.Context, orgID, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, queryDeleteRateMultiTenant, id, orgID)
	if err != nil {
		return fmt.Errorf("failed to delete rate: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("metal rate not found")
	}
	return nil
}
