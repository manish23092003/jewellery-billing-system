package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ── User Roles ─────────────────────────────────────────────────────────

// UserRole represents the access level of a user in the system.
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleStaff UserRole = "staff"
)

// IsValid checks whether the role is one of the allowed values.
func (r UserRole) IsValid() bool {
	return r == RoleAdmin || r == RoleStaff
}

// ── User Entity ────────────────────────────────────────────────────────

// User is the core domain entity representing a system user.
type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never serialized to JSON
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToResponse strips sensitive fields and returns a safe projection.
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ── Request / Response DTOs ────────────────────────────────────────────

type CreateUserRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}

type UpdateUserRequest struct {
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email,omitempty"`
	Password string   `json:"password,omitempty"`
	Role     UserRole `json:"role,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// ── Repository Interface ───────────────────────────────────────────────

// UserRepository defines the data-access contract for user persistence.
// Implementations live in the repository package (PostgreSQL).
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
}
