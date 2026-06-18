package domain

import "context"

// DashboardMetrics represents the high-level summary cards.
type DashboardMetrics struct {
	TodaySales      float64 `json:"today_sales"`
	TodayExpenses   float64 `json:"today_expenses"`
	TodayProfit     float64 `json:"today_profit"`
	MonthlySales    float64 `json:"monthly_sales"`
	MonthlyExpenses float64 `json:"monthly_expenses"`
	MonthlyProfit   float64 `json:"monthly_profit"`
}

// DailyTrend represents a single point on the trend chart.
type DailyTrend struct {
	Date     string  `json:"date"`
	Sales    float64 `json:"sales"`
	Expenses float64 `json:"expenses"`
	Profit   float64 `json:"profit"`
}

// DashboardData encapsulates the entire dashboard response.
type DashboardData struct {
	Metrics DashboardMetrics `json:"metrics"`
	Trends  []DailyTrend     `json:"trends"`
}

// AnalyticsRepository handles complex aggregation queries.
type AnalyticsRepository interface {
	GetDashboardMetrics(ctx context.Context, today, firstDayOfMonth string) (*DashboardMetrics, error)
	GetMonthlyTrends(ctx context.Context, firstDayOfMonth, lastDayOfMonth string) ([]DailyTrend, error)
}
