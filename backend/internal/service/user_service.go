package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"jewellery-billing/internal/domain"
)

// UserService handles user CRUD operations (admin-only).
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService creates a UserService with the given repository.
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Create registers a new user within an organization after validating uniqueness and hashing the password.
func (s *UserService) Create(ctx context.Context, orgID uuid.UUID, req domain.CreateUserRequest) (*domain.UserResponse, error) {
	// Guard against duplicate emails.
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, fmt.Errorf("email already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		OrganizationID: orgID,
		Name:           req.Name,
		Email:          req.Email,
		PasswordHash:   string(hash),
		Role:           req.Role,
		IsActive:       true,
		EmailVerified:  false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	resp := user.ToResponse()
	return &resp, nil
}

// GetByID returns a single user by ID (safe projection — no password hash).
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := user.ToResponse()
	return &resp, nil
}

// GetAll returns every user in the organization (admin-only endpoint).
func (s *UserService) GetAll(ctx context.Context, orgID uuid.UUID) ([]domain.UserResponse, error) {
	users, err := s.userRepo.GetAllByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.UserResponse, len(users))
	for i, u := range users {
		responses[i] = u.ToResponse()
	}
	return responses, nil
}

// Update modifies an existing user. Empty fields in the request are skipped.
func (s *UserService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		// Prevent stealing another user's email.
		existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
		if existing != nil && existing.ID != id {
			return nil, fmt.Errorf("email already in use")
		}
		user.Email = req.Email
	}

	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hash)
	}

	if req.Role != "" && req.Role.IsValid() {
		user.Role = req.Role
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	resp := user.ToResponse()
	return &resp, nil
}

// Delete removes a user by ID.
func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
