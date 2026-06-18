package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"jewellery-billing/internal/config"
	"jewellery-billing/internal/domain"
)

// JWTClaims extends the standard JWT claims with application-specific fields.
type JWTClaims struct {
	UserID uuid.UUID       `json:"user_id"`
	Email  string          `json:"email"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// AuthService handles login, token generation, and token validation.
type AuthService struct {
	userRepo domain.UserRepository
	config   *config.Config
}

// NewAuthService creates an AuthService with the given dependencies.
func NewAuthService(userRepo domain.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// Login verifies credentials and returns a pair of JWTs (access + refresh).
func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	// Look up the user by email.
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Deliberately vague to prevent email enumeration.
		return nil, fmt.Errorf("invalid email or password")
	}

	// Compare the supplied password against the stored bcrypt hash.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate token pair.
	accessToken, err := s.generateToken(user, time.Duration(s.config.JWTAccessExpiryMin)*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(user, time.Duration(s.config.JWTRefreshExpiryDay)*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	}, nil
}

// RefreshToken validates an existing refresh token and issues a new pair.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Re-fetch the user to ensure they still exist and their role hasn't changed.
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user no longer exists")
	}

	newAccess, err := s.generateToken(user, time.Duration(s.config.JWTAccessExpiryMin)*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefresh, err := s.generateToken(user, time.Duration(s.config.JWTRefreshExpiryDay)*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &domain.AuthResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
		User:         user.ToResponse(),
	}, nil
}

// ValidateToken parses and validates a JWT string, returning its claims.
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC (prevent algorithm-switching attacks).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// generateToken creates a signed JWT with the given expiry duration.
func (s *AuthService) generateToken(user *domain.User, expiry time.Duration) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "jewellery-billing",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}
