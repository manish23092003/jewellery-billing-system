package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ── Audit Log Entity ──────────────────────────────────────────────────

// AuditLog records a business-critical action performed within an organization.
type AuditLog struct {
	ID             uuid.UUID              `json:"id"`
	OrganizationID uuid.UUID              `json:"organization_id"`
	UserID         *uuid.UUID             `json:"user_id,omitempty"`
	Action         string                 `json:"action"`
	EntityType     string                 `json:"entity_type"`
	EntityID       string                 `json:"entity_id,omitempty"`
	Details        map[string]interface{} `json:"details,omitempty"`
	IPAddress      string                 `json:"ip_address,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}

// ── Repository Interface ──────────────────────────────────────────────

// AuditRepository defines the data-access contract for audit log persistence.
type AuditRepository interface {
	Create(ctx context.Context, log *AuditLog) error
	GetByOrganization(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]AuditLog, int64, error)
}
