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
		INSERT INTO expenses (category, amount, description, expense_date, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
	`
	return r.db.QueryRow(ctx, query, e.Category, e.Amount, e.Description, e.ExpenseDate, e.CreatedBy).
		Scan(&e.ID, &e.CreatedAt)
}

func (r *ExpenseRepo) GetAll(ctx context.Context, filter domain.ExpenseFilter) ([]domain.Expense, int64, error) {
	var conditions []string
	var args []interface{}
	argID := 1

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

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM expenses " + whereClause
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch paginated rows
	offset := (filter.Page - 1) * filter.PerPage
	query := fmt.Sprintf(`
		SELECT id, category, amount, description, TO_CHAR(expense_date, 'YYYY-MM-DD'), created_by, TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM expenses
		%s
		ORDER BY expense_date DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argID, argID+1)

	args = append(args, filter.PerPage, offset)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []domain.Expense
	for rows.Next() {
		var e domain.Expense
		var dt string
		if err := rows.Scan(&e.ID, &e.Category, &e.Amount, &e.Description, &dt, &e.CreatedBy, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		e.ExpenseDate = dt
		expenses = append(expenses, e)
	}

	return expenses, total, nil
}

func (r *ExpenseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM expenses WHERE id = $1`
	ct, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("expense not found")
	}
	return nil
}
