package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"jewellery-billing/internal/domain"
)

// ExpenseService handles business logic for expenses.
type ExpenseService struct {
	repo domain.ExpenseRepository
}

func NewExpenseService(repo domain.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

// Create validates and persists a new expense.
func (s *ExpenseService) Create(ctx context.Context, req domain.CreateExpenseRequest, createdBy uuid.UUID) (*domain.Expense, error) {
	if req.Category == "" {
		return nil, fmt.Errorf("category is required")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}

	expenseDate := req.ExpenseDate
	if expenseDate == "" {
		expenseDate = "now()"
	}

	expense := &domain.Expense{
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
		ExpenseDate: expenseDate,
		CreatedBy:   createdBy,
	}

	if err := s.repo.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("failed to create expense: %w", err)
	}

	return expense, nil
}

// GetAll returns a paginated list of expenses.
func (s *ExpenseService) GetAll(ctx context.Context, filter domain.ExpenseFilter) ([]domain.Expense, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PerPage <= 0 {
		filter.PerPage = 20
	}
	return s.repo.GetAll(ctx, filter)
}

// Delete removes an expense (typically admin only).
func (s *ExpenseService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
