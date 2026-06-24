package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"jewellery-billing/internal/domain"
)

// auditRepository is the PostgreSQL implementation of domain.AuditRepository.
type auditRepository struct {
	db *pgxpool.Pool
}

// NewAuditRepository returns a production-ready AuditRepository backed by pgx.
func NewAuditRepository(db *pgxpool.Pool) domain.AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	detailsJSON, _ := json.Marshal(log.Details)

	query := `
		INSERT INTO audit_logs (organization_id, user_id, action, entity_type, entity_id, details, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`
	return r.db.QueryRow(ctx, query,
		log.OrganizationID, log.UserID, log.Action, log.EntityType,
		log.EntityID, detailsJSON, log.IPAddress,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *auditRepository) GetByOrganization(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]domain.AuditLog, int64, error) {
	// Count total
	var total int64
	if err := r.db.QueryRow(ctx,
		"SELECT COUNT(*) FROM audit_logs WHERE organization_id = $1", orgID,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	query := `
		SELECT id, organization_id, user_id, action, entity_type, entity_id, details, ip_address, created_at
		FROM audit_logs
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.AuditLog
	for rows.Next() {
		var log domain.AuditLog
		var detailsJSON []byte
		if err := rows.Scan(
			&log.ID, &log.OrganizationID, &log.UserID, &log.Action,
			&log.EntityType, &log.EntityID, &detailsJSON, &log.IPAddress, &log.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit log: %w", err)
		}
		if detailsJSON != nil {
			_ = json.Unmarshal(detailsJSON, &log.Details)
		}
		logs = append(logs, log)
	}

	return logs, total, rows.Err()
}
