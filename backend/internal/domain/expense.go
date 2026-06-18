package domain

import (
	"context"

	"github.com/google/uuid"
)

// Expense represents a daily shop expenditure.
type Expense struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Category    string    `json:"category" db:"category"`
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	ExpenseDate string    `json:"expense_date" db:"expense_date"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt   string    `json:"created_at" db:"created_at"`
}

// CreateExpenseRequest is the payload to create an expense.
type CreateExpenseRequest struct {
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	ExpenseDate string  `json:"expense_date"`
}

// ExpenseFilter for pagination and search.
type ExpenseFilter struct {
	Category string
	DateFrom string
	DateTo   string
	Page     int
	PerPage  int
}

// ExpenseRepository handles persistent storage operations for expenses.
type ExpenseRepository interface {
	Create(ctx context.Context, expense *Expense) error
	GetAll(ctx context.Context, filter ExpenseFilter) ([]Expense, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
