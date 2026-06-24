package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"jewellery-billing/internal/domain"
)

type CustomerService struct {
	customerRepo domain.CustomerRepository
}

func NewCustomerService(repo domain.CustomerRepository) *CustomerService {
	return &CustomerService{customerRepo: repo}
}

func (s *CustomerService) Create(ctx context.Context, orgID uuid.UUID, req domain.CreateCustomerRequest) (*domain.Customer, error) {
	if req.Name == "" || req.Phone == "" {
		return nil, fmt.Errorf("name and phone are required")
	}

	req.Phone = normalizePhone(req.Phone)

	// Check if customer with same phone already exists
	existing, _ := s.customerRepo.GetByPhone(ctx, orgID, req.Phone)
	if existing != nil {
		return nil, fmt.Errorf("customer with phone %s already exists", req.Phone)
	}

	customer := &domain.Customer{
		OrganizationID: orgID,
		Name:           req.Name,
		Phone:          req.Phone,
		Email:          req.Email,
		Address:        req.Address,
		TotalPurchases: 0,
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) GetByID(ctx context.Context, orgID, id uuid.UUID) (*domain.Customer, error) {
	return s.customerRepo.GetByID(ctx, orgID, id)
}

func (s *CustomerService) Update(ctx context.Context, orgID, id uuid.UUID, req domain.UpdateCustomerRequest) (*domain.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.Phone != "" {
		req.Phone = normalizePhone(req.Phone)
		// Check phone uniqueness if changed
		if req.Phone != customer.Phone {
			existing, _ := s.customerRepo.GetByPhone(ctx, orgID, req.Phone)
			if existing != nil {
				return nil, fmt.Errorf("customer with phone %s already exists", req.Phone)
			}
		}
		customer.Phone = req.Phone
	}
	customer.Email = req.Email
	customer.Address = req.Address

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) Delete(ctx context.Context, orgID, id uuid.UUID) error {
	return s.customerRepo.Delete(ctx, orgID, id)
}

func (s *CustomerService) GetAll(ctx context.Context, orgID uuid.UUID, search string, limit, offset int) ([]domain.Customer, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.customerRepo.GetAll(ctx, orgID, search, limit, offset)
}
