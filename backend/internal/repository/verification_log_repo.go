package repository

import (
	"context"
	"fmt"

	"jewellery-billing/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VerificationLogRepository interface {
	LogVerification(ctx context.Context, log *domain.VerificationLog) error
}

type verificationLogRepository struct {
	db *pgxpool.Pool
}

func NewVerificationLogRepository(db *pgxpool.Pool) VerificationLogRepository {
	return &verificationLogRepository{db: db}
}

const queryInsertVerificationLog = `
	INSERT INTO verification_logs (token, ip_address, user_agent, is_valid, failure_reason)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at`

func (r *verificationLogRepository) LogVerification(ctx context.Context, log *domain.VerificationLog) error {
	err := r.db.QueryRow(ctx, queryInsertVerificationLog,
		log.Token, log.IPAddress, log.UserAgent, log.IsValid, log.FailureReason,
	).Scan(&log.ID, &log.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert verification log: %w", err)
	}
	return nil
}
