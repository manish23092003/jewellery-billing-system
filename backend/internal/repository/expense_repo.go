package repository

import (
	"context"
	"fmt"
	"strings"

	"jewellery-billing/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ExpenseRepo implements domain.ExpenseRepository using pgx.
type ExpenseRepo struct {
	db *pgxpool.Pool
}

func NewExpenseRepository(db *pgxpool.Pool) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

func (r *ExpenseRepo) Create(ctx context.Context, e *domain.Expense) error {
	query := `
		INSERT INTO expenses (organization_id, category, amount, description, expense_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
	`
	return r.db.QueryRow(ctx, query,
		e.OrganizationID, e.Category, e.Amount, e.Description, e.ExpenseDate, e.CreatedBy,
	).Scan(&e.ID, &e.CreatedAt)
}

func (r *ExpenseRepo) GetAll(ctx context.Context, orgID uuid.UUID, filter domain.ExpenseFilter) ([]domain.Expense, int64, error) {
	// Always filter by organization
	conditions := []string{"organization_id = $1"}
	args := []interface{}{orgID}
	argID := 2

	if filter.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argID))
		args = append(args, filter.Category)
		argID++
	}
	if filter.DateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("expense_date >= $%d", argID))
		args = append(args, filter.DateFrom)
		argID++
	}
	if filter.DateTo != "" {
		conditions = append(conditions, fmt.Sprintf("expense_date <= $%d", argID))
		args = append(args, filter.DateTo)
		argID++
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// Count total
	countQuery := "SELECT COUNT(*) FROM expenses " + whereClause
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch paginated rows
	perPage := filter.PerPage
	if perPage <= 0 {
		perPage = 20
	}
	offset := (filter.Page - 1) * perPage
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT id, organization_id, category, amount, description,
		       TO_CHAR(expense_date, 'YYYY-MM-DD'), created_by,
		       TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM expenses
		%s
		ORDER BY expense_date DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argID, argID+1)

	args = append(args, perPage, offset)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []domain.Expense
	for rows.Next() {
		var e domain.Expense
		var dt string
		if err := rows.Scan(
			&e.ID, &e.OrganizationID, &e.Category, &e.Amount, &e.Description,
			&dt, &e.CreatedBy, &e.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		e.ExpenseDate = dt
		expenses = append(expenses, e)
	}

	return expenses, total, nil
}

func (r *ExpenseRepo) Delete(ctx context.Context, orgID, id uuid.UUID) error {
	query := `DELETE FROM expenses WHERE id = $1 AND organization_id = $2`
	ct, err := r.db.Exec(ctx, query, id, orgID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("expense not found")
	}
	return nil
}
