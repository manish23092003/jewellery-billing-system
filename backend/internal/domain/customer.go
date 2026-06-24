package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Customer represents a client of the jewellery shop.
type Customer struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Address        string    `json:"address"`
	TotalPurchases float64   `json:"total_purchases"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateCustomerRequest is the DTO for creating a new customer.
type CreateCustomerRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// UpdateCustomerRequest is the DTO for updating an existing customer.
type UpdateCustomerRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// CustomerRepository defines the database operations for customers.
type CustomerRepository interface {
	Create(ctx context.Context, customer *Customer) error
	GetByID(ctx context.Context, orgID, id uuid.UUID) (*Customer, error)
	GetByPhone(ctx context.Context, orgID uuid.UUID, phone string) (*Customer, error)
	Update(ctx context.Context, customer *Customer) error
	Delete(ctx context.Context, orgID, id uuid.UUID) error

	// GetAll returns paginated customers for an organization, optionally searching by name/phone.
	GetAll(ctx context.Context, orgID uuid.UUID, search string, limit, offset int) ([]Customer, int64, error)
}
