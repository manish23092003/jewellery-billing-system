package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"jewellery-billing/internal/domain"
)

type PostgresCustomerRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCustomerRepository(db *pgxpool.Pool) *PostgresCustomerRepository {
	return &PostgresCustomerRepository{db: db}
}

func (r *PostgresCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (organization_id, name, phone, email, address, total_purchases)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		customer.OrganizationID, customer.Name, customer.Phone, customer.Email, customer.Address, customer.TotalPurchases,
	).Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert customer: %w", err)
	}
	return nil
}

func (r *PostgresCustomerRepository) GetByID(ctx context.Context, orgID, id uuid.UUID) (*domain.Customer, error) {
	query := `
		SELECT id, organization_id, name, phone, email, address, total_purchases, created_at, updated_at
		FROM customers
		WHERE organization_id = $1 AND id = $2
	`
	var c domain.Customer
	err := r.db.QueryRow(ctx, query, orgID, id).Scan(
		&c.ID, &c.OrganizationID, &c.Name, &c.Phone, &c.Email, &c.Address, &c.TotalPurchases, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	return &c, nil
}

func (r *PostgresCustomerRepository) GetByPhone(ctx context.Context, orgID uuid.UUID, phone string) (*domain.Customer, error) {
	query := `
		SELECT id, organization_id, name, phone, email, address, total_purchases, created_at, updated_at
		FROM customers
		WHERE organization_id = $1 AND phone = $2
		LIMIT 1
	`
	var c domain.Customer
	err := r.db.QueryRow(ctx, query, orgID, phone).Scan(
		&c.ID, &c.OrganizationID, &c.Name, &c.Phone, &c.Email, &c.Address, &c.TotalPurchases, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("customer not found by phone: %w", err)
	}
	return &c, nil
}

func (r *PostgresCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	query := `
		UPDATE customers
		SET name = $1, phone = $2, email = $3, address = $4, total_purchases = $5, updated_at = NOW()
		WHERE organization_id = $6 AND id = $7
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		customer.Name, customer.Phone, customer.Email, customer.Address, customer.TotalPurchases,
		customer.OrganizationID, customer.ID,
	).Scan(&customer.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

func (r *PostgresCustomerRepository) Delete(ctx context.Context, orgID, id uuid.UUID) error {
	query := `DELETE FROM customers WHERE organization_id = $1 AND id = $2`
	tag, err := r.db.Exec(ctx, query, orgID, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customer not found")
	}
	return nil
}

func (r *PostgresCustomerRepository) GetAll(ctx context.Context, orgID uuid.UUID, search string, limit, offset int) ([]domain.Customer, int64, error) {
	searchPattern := "%" + search + "%"

	countQuery := `
		SELECT COUNT(*) FROM customers
		WHERE organization_id = $1 AND (name ILIKE $2 OR phone ILIKE $2)
	`
	var total int64
	err := r.db.QueryRow(ctx, countQuery, orgID, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}

	query := `
		SELECT id, organization_id, name, phone, email, address, total_purchases, created_at, updated_at
		FROM customers
		WHERE organization_id = $1 AND (name ILIKE $2 OR phone ILIKE $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.Query(ctx, query, orgID, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query customers: %w", err)
	}
	defer rows.Close()

	var customers []domain.Customer
	for rows.Next() {
		var c domain.Customer
		if err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.Name, &c.Phone, &c.Email, &c.Address, &c.TotalPurchases, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		customers = append(customers, c)
	}
	return customers, total, nil
}
