package repository

import (
	"context"

	"jewellery-billing/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AnalyticsRepo implements domain.AnalyticsRepository using pgx.
type AnalyticsRepo struct {
	db *pgxpool.Pool
}

func NewAnalyticsRepository(db *pgxpool.Pool) *AnalyticsRepo {
	return &AnalyticsRepo{db: db}
}

func (r *AnalyticsRepo) GetDashboardMetrics(ctx context.Context, orgID uuid.UUID, today, firstDayOfMonth string) (*domain.DashboardMetrics, error) {
	var metrics domain.DashboardMetrics

	// Calculate Today's Sales
	queryTodaySales := `SELECT COALESCE(SUM(grand_total), 0) FROM bills WHERE organization_id = $1 AND DATE(invoice_date) = $2`
	r.db.QueryRow(ctx, queryTodaySales, orgID, today).Scan(&metrics.TodaySales)

	// Calculate Today's Expenses
	queryTodayExpenses := `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE organization_id = $1 AND expense_date = $2`
	r.db.QueryRow(ctx, queryTodayExpenses, orgID, today).Scan(&metrics.TodayExpenses)

	metrics.TodayProfit = metrics.TodaySales - metrics.TodayExpenses

	// Calculate Monthly Sales
	queryMonthlySales := `SELECT COALESCE(SUM(grand_total), 0) FROM bills WHERE organization_id = $1 AND DATE(invoice_date) >= $2`
	r.db.QueryRow(ctx, queryMonthlySales, orgID, firstDayOfMonth).Scan(&metrics.MonthlySales)

	// Calculate Monthly Expenses
	queryMonthlyExpenses := `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE organization_id = $1 AND expense_date >= $2`
	r.db.QueryRow(ctx, queryMonthlyExpenses, orgID, firstDayOfMonth).Scan(&metrics.MonthlyExpenses)

	metrics.MonthlyProfit = metrics.MonthlySales - metrics.MonthlyExpenses

	return &metrics, nil
}

func (r *AnalyticsRepo) GetMonthlyTrends(ctx context.Context, orgID uuid.UUID, startDate, endDate string) ([]domain.DailyTrend, error) {
	query := `
		WITH dates AS (
			SELECT generate_series($1::DATE, $2::DATE, '1 day'::interval)::DATE AS date
		),
		daily_sales AS (
			SELECT DATE(invoice_date) AS date, SUM(grand_total) AS total_sales
			FROM bills
			WHERE organization_id = $3 AND DATE(invoice_date) BETWEEN $1 AND $2
			GROUP BY DATE(invoice_date)
		),
		daily_expenses AS (
			SELECT expense_date AS date, SUM(amount) AS total_expenses
			FROM expenses
			WHERE organization_id = $3 AND expense_date BETWEEN $1 AND $2
			GROUP BY expense_date
		)
		SELECT
			TO_CHAR(d.date, 'YYYY-MM-DD'),
			COALESCE(s.total_sales, 0) AS sales,
			COALESCE(e.total_expenses, 0) AS expenses,
			COALESCE(s.total_sales, 0) - COALESCE(e.total_expenses, 0) AS profit
		FROM dates d
		LEFT JOIN daily_sales s ON d.date = s.date
		LEFT JOIN daily_expenses e ON d.date = e.date
		ORDER BY d.date ASC;
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []domain.DailyTrend
	for rows.Next() {
		var t domain.DailyTrend
		if err := rows.Scan(&t.Date, &t.Sales, &t.Expenses, &t.Profit); err != nil {
			return nil, err
		}
		trends = append(trends, t)
	}

	return trends, nil
}
