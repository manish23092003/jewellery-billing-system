package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ── Organization Entity ───────────────────────────────────────────────

// Organization represents a jewellery business (tenant) in the SaaS platform.
type Organization struct {
	ID                 uuid.UUID `json:"id"`
	BusinessName       string    `json:"business_name"`
	OwnerName          string    `json:"owner_name"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	GSTIN              string    `json:"gstin"`
	Address            string    `json:"address"`
	SubscriptionStatus string    `json:"subscription_status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ToResponse returns a safe projection of the organization.
func (o *Organization) ToResponse() OrganizationResponse {
	return OrganizationResponse{
		ID:                 o.ID,
		BusinessName:       o.BusinessName,
		OwnerName:          o.OwnerName,
		Email:              o.Email,
		Phone:              o.Phone,
		GSTIN:              o.GSTIN,
		Address:            o.Address,
		SubscriptionStatus: o.SubscriptionStatus,
		CreatedAt:          o.CreatedAt,
		UpdatedAt:          o.UpdatedAt,
	}
}

// ── Request / Response DTOs ───────────────────────────────────────────

// RegisterRequest is the payload for registering a new business (organization + owner).
type RegisterRequest struct {
	BusinessName    string `json:"business_name"`
	OwnerName       string `json:"owner_name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// OrganizationResponse is the safe projection of an organization.
type OrganizationResponse struct {
	ID                 uuid.UUID `json:"id"`
	BusinessName       string    `json:"business_name"`
	OwnerName          string    `json:"owner_name"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	GSTIN              string    `json:"gstin"`
	Address            string    `json:"address"`
	SubscriptionStatus string    `json:"subscription_status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// UpdateOrganizationRequest is the payload for updating an organization.
type UpdateOrganizationRequest struct {
	BusinessName string `json:"business_name,omitempty"`
	OwnerName    string `json:"owner_name,omitempty"`
	Phone        string `json:"phone,omitempty"`
	GSTIN        string `json:"gstin,omitempty"`
	Address      string `json:"address,omitempty"`
}

// ── Repository Interface ──────────────────────────────────────────────

// OrganizationRepository defines the data-access contract for organizations.
type OrganizationRepository interface {
	Create(ctx context.Context, org *Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*Organization, error)
	GetByEmail(ctx context.Context, email string) (*Organization, error)
	Update(ctx context.Context, org *Organization) error
}
