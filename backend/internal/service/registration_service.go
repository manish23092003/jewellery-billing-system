package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"jewellery-billing/internal/domain"
)

// RegistrationService handles the business registration flow.
type RegistrationService struct {
	orgRepo     domain.OrganizationRepository
	userRepo    domain.UserRepository
	tokenRepo   domain.TokenRepository
	settingRepo domain.SettingRepository
	authService *AuthService
	emailSender EmailSender
}

// NewRegistrationService creates a RegistrationService with all required dependencies.
func NewRegistrationService(
	orgRepo domain.OrganizationRepository,
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	settingRepo domain.SettingRepository,
	authService *AuthService,
	emailSender EmailSender,
) *RegistrationService {
	return &RegistrationService{
		orgRepo:     orgRepo,
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		settingRepo: settingRepo,
		authService: authService,
		emailSender: emailSender,
	}
}

// Register creates a new organization + admin user and returns auth tokens.
func (s *RegistrationService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// ── Validate ──────────────────────────────────────────────────────
	if req.BusinessName == "" {
		return nil, fmt.Errorf("business name is required")
	}
	if req.OwnerName == "" {
		return nil, fmt.Errorf("owner name is required")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if len(req.Password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}
	if req.Password != req.ConfirmPassword {
		return nil, fmt.Errorf("passwords do not match")
	}

	// Check if email is already registered (user table)
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, fmt.Errorf("email is already registered")
	}

	// Check if organization email is already taken
	existingOrg, _ := s.orgRepo.GetByEmail(ctx, req.Email)
	if existingOrg != nil {
		return nil, fmt.Errorf("an organization with this email already exists")
	}

	// ── Create Organization ───────────────────────────────────────────
	org := &domain.Organization{
		BusinessName: req.BusinessName,
		OwnerName:    req.OwnerName,
		Email:        req.Email,
		Phone:        req.Phone,
	}
	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// ── Create Admin User ─────────────────────────────────────────────
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		OrganizationID: org.ID,
		Name:           req.OwnerName,
		Email:          req.Email,
		PasswordHash:   string(hash),
		Role:           domain.RoleAdmin,
		IsActive:       true,
		EmailVerified:  true, // Auto-verified since email sending is disabled
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	// ── Create Default Settings ───────────────────────────────────────
	defaultSettings := &domain.ShopSettings{
		OrganizationID: org.ID,
		ShopName:       req.BusinessName,
		Phone:          req.Phone,
		InvoicePrefix:  "INV",
	}
	if err := s.settingRepo.Upsert(ctx, defaultSettings); err != nil {
		// Non-fatal — log and continue
		fmt.Printf("Warning: failed to create default settings: %v\n", err)
	}

	// ── Generate Auth Tokens ──────────────────────────────────────────
	authResponse, err := s.authService.GenerateAuthResponse(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth tokens: %w", err)
	}

	return authResponse, nil
}
