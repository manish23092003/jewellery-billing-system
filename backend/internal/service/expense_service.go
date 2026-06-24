package service

import (
	"context"

	"github.com/google/uuid"

	"jewellery-billing/internal/domain"
)

// ExpenseService handles expense CRUD operations.
type ExpenseService struct {
	expenseRepo domain.ExpenseRepository
}

func NewExpenseService(expenseRepo domain.ExpenseRepository) *ExpenseService {
	return &ExpenseService{expenseRepo: expenseRepo}
}

func (s *ExpenseService) Create(ctx context.Context, orgID uuid.UUID, req domain.CreateExpenseRequest, createdBy uuid.UUID) (*domain.Expense, error) {
	expense := &domain.Expense{
		OrganizationID: orgID,
		Category:       req.Category,
		Amount:         req.Amount,
		Description:    req.Description,
		ExpenseDate:    req.ExpenseDate,
		CreatedBy:      createdBy,
	}

	if err := s.expenseRepo.Create(ctx, expense); err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *ExpenseService) GetAll(ctx context.Context, orgID uuid.UUID, filter domain.ExpenseFilter) ([]domain.Expense, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PerPage <= 0 {
		filter.PerPage = 20
	}
	return s.expenseRepo.GetAll(ctx, orgID, filter)
}

func (s *ExpenseService) Delete(ctx context.Context, orgID, id uuid.UUID) error {
	return s.expenseRepo.Delete(ctx, orgID, id)
}
