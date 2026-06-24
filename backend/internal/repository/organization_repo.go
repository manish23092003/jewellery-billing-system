package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"jewellery-billing/internal/domain"
)

// organizationRepository is the PostgreSQL implementation of domain.OrganizationRepository.
type organizationRepository struct {
	db *pgxpool.Pool
}

// NewOrganizationRepository returns a production-ready OrganizationRepository backed by pgx.
func NewOrganizationRepository(db *pgxpool.Pool) domain.OrganizationRepository {
	return &organizationRepository{db: db}
}

// ── Queries ────────────────────────────────────────────────────────────

const (
	queryInsertOrg = `
		INSERT INTO organizations (business_name, owner_name, email, phone, gstin, address)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, subscription_status, created_at, updated_at`

	queryGetOrgByID = `
		SELECT id, business_name, owner_name, email, phone, gstin, address,
		       subscription_status, created_at, updated_at
		FROM organizations WHERE id = $1`

	queryGetOrgByEmail = `
		SELECT id, business_name, owner_name, email, phone, gstin, address,
		       subscription_status, created_at, updated_at
		FROM organizations WHERE email = $1`

	queryUpdateOrg = `
		UPDATE organizations
		SET business_name = $2, owner_name = $3, phone = $4, gstin = $5,
		    address = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`
)

// ── Implementation ─────────────────────────────────────────────────────

func (r *organizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	return r.db.QueryRow(ctx, queryInsertOrg,
		org.BusinessName, org.OwnerName, org.Email, org.Phone, org.GSTIN, org.Address,
	).Scan(&org.ID, &org.SubscriptionStatus, &org.CreatedAt, &org.UpdatedAt)
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	org := &domain.Organization{}
	err := r.db.QueryRow(ctx, queryGetOrgByID, id).Scan(
		&org.ID, &org.BusinessName, &org.OwnerName, &org.Email, &org.Phone,
		&org.GSTIN, &org.Address, &org.SubscriptionStatus, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return org, nil
}

func (r *organizationRepository) GetByEmail(ctx context.Context, email string) (*domain.Organization, error) {
	org := &domain.Organization{}
	err := r.db.QueryRow(ctx, queryGetOrgByEmail, email).Scan(
		&org.ID, &org.BusinessName, &org.OwnerName, &org.Email, &org.Phone,
		&org.GSTIN, &org.Address, &org.SubscriptionStatus, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return org, nil
}

func (r *organizationRepository) Update(ctx context.Context, org *domain.Organization) error {
	err := r.db.QueryRow(ctx, queryUpdateOrg,
		org.ID, org.BusinessName, org.OwnerName, org.Phone, org.GSTIN, org.Address,
	).Scan(&org.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("organization not found")
		}
		return fmt.Errorf("failed to update organization: %w", err)
	}
	return nil
}
