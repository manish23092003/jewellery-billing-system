package service

import (
	"context"
	"time"

	"jewellery-billing/internal/domain"
)

// AnalyticsService handles dashboard data aggregation.
type AnalyticsService struct {
	repo domain.AnalyticsRepository
}

func NewAnalyticsService(repo domain.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

// GetDashboard returns the full dashboard payload including metrics and trend charts.
func (s *AnalyticsService) GetDashboard(ctx context.Context) (*domain.DashboardData, error) {
	now := time.Now()
	today := now.Format("2006-01-02")
	
	// First day of current month
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	
	// Last day of current month
	lastDay := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	metrics, err := s.repo.GetDashboardMetrics(ctx, today, firstDay)
	if err != nil {
		return nil, err
	}

	trends, err := s.repo.GetMonthlyTrends(ctx, firstDay, lastDay)
	if err != nil {
		return nil, err
	}

	return &domain.DashboardData{
		Metrics: *metrics,
		Trends:  trends,
	}, nil
}
