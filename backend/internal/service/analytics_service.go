package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"jewellery-billing/internal/domain"
)

// AnalyticsService handles analytics and dashboard operations.
type AnalyticsService struct {
	analyticsRepo domain.AnalyticsRepository
}

func NewAnalyticsService(repo domain.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{analyticsRepo: repo}
}

func (s *AnalyticsService) GetDashboard(ctx context.Context, orgID uuid.UUID) (*domain.DashboardData, error) {
	now := time.Now()
	today := now.Format("2006-01-02")
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	lastDayOfMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	metrics, err := s.analyticsRepo.GetDashboardMetrics(ctx, orgID, today, firstDayOfMonth)
	if err != nil {
		return nil, err
	}

	trends, err := s.analyticsRepo.GetMonthlyTrends(ctx, orgID, firstDayOfMonth, lastDayOfMonth)
	if err != nil {
		return nil, err
	}

	return &domain.DashboardData{
		Metrics: *metrics,
		Trends:  trends,
	}, nil
}
