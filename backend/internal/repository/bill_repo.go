package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"jewellery-billing/internal/domain"
)

type billRepository struct {
	db *pgxpool.Pool
}

func NewBillRepository(db *pgxpool.Pool) domain.BillRepository {
	return &billRepository{db: db}
}

// ── Queries ────────────────────────────────────────────────────────────

const (
	queryInsertBill = `
		INSERT INTO bills (invoice_number, invoice_date, customer_name, customer_phone,
		                   subtotal, gst_amount, grand_total, payment_method, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`

	queryInsertBillItem = `
		INSERT INTO bill_items (bill_id, item_name, metal_type, purity, weight,
		                        rate_per_gram, making_charge, gst_percentage, quantity, line_total)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`

	queryGetBillByID = `
		SELECT b.id, b.invoice_number, TO_CHAR(b.invoice_date, 'YYYY-MM-DD'),
		       COALESCE(b.customer_name, ''), COALESCE(b.customer_phone, ''),
		       b.subtotal, b.gst_amount, b.grand_total, b.payment_method,
		       COALESCE(b.notes, ''), b.created_by, b.created_at,
		       COALESCE(u.name, 'Unknown')
		FROM bills b
		LEFT JOIN users u ON u.id = b.created_by
		WHERE b.id = $1`

	queryGetBillItems = `
		SELECT id, bill_id, item_name, metal_type, purity, weight,
		       rate_per_gram, making_charge, gst_percentage, quantity, line_total, created_at
		FROM bill_items
		WHERE bill_id = $1
		ORDER BY created_at`

	queryDeleteBill = `DELETE FROM bills WHERE id = $1`

	// Gets the latest invoice number for the current month to determine next sequence
	queryLatestInvoice = `
		SELECT invoice_number FROM bills
		WHERE invoice_number LIKE $1
		ORDER BY invoice_number DESC
		LIMIT 1`
)

// ── Create (Transactional) ─────────────────────────────────────────────

func (r *billRepository) Create(ctx context.Context, bill *domain.Bill) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert bill header
	err = tx.QueryRow(ctx, queryInsertBill,
		bill.InvoiceNumber, bill.InvoiceDate, bill.CustomerName, bill.CustomerPhone,
		bill.Subtotal, bill.GSTAmount, bill.GrandTotal,
		bill.PaymentMethod, bill.Notes, bill.CreatedBy,
	).Scan(&bill.ID, &bill.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert bill: %w", err)
	}

	// Insert each line item
	for i := range bill.Items {
		bill.Items[i].BillID = bill.ID
		err = tx.QueryRow(ctx, queryInsertBillItem,
			bill.Items[i].BillID, bill.Items[i].ItemName,
			bill.Items[i].MetalType, bill.Items[i].Purity,
			bill.Items[i].Weight, bill.Items[i].RatePerGram,
			bill.Items[i].MakingCharge, bill.Items[i].GSTPercentage,
			bill.Items[i].Quantity, bill.Items[i].LineTotal,
		).Scan(&bill.Items[i].ID, &bill.Items[i].CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert bill item %d: %w", i+1, err)
		}
	}

	return tx.Commit(ctx)
}

// ── Get By ID ──────────────────────────────────────────────────────────

func (r *billRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Bill, error) {
	bill := &domain.Bill{}
	err := r.db.QueryRow(ctx, queryGetBillByID, id).Scan(
		&bill.ID, &bill.InvoiceNumber, &bill.InvoiceDate,
		&bill.CustomerName, &bill.CustomerPhone,
		&bill.Subtotal, &bill.GSTAmount, &bill.GrandTotal,
		&bill.PaymentMethod, &bill.Notes, &bill.CreatedBy,
		&bill.CreatedAt, &bill.CreatedByName,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("bill not found")
		}
		return nil, fmt.Errorf("failed to get bill: %w", err)
	}

	// Fetch line items
	rows, err := r.db.Query(ctx, queryGetBillItems, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bill items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.BillItem
		if err := rows.Scan(
			&item.ID, &item.BillID, &item.ItemName, &item.MetalType,
			&item.Purity, &item.Weight, &item.RatePerGram,
			&item.MakingCharge, &item.GSTPercentage, &item.Quantity,
			&item.LineTotal, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan bill item: %w", err)
		}
		bill.Items = append(bill.Items, item)
	}

	return bill, rows.Err()
}

// ── Get All (Paginated + Filtered) ─────────────────────────────────────

func (r *billRepository) GetAll(ctx context.Context, filter domain.BillFilter) ([]domain.Bill, int64, error) {
	baseSelect := `SELECT b.id, b.invoice_number, TO_CHAR(b.invoice_date, 'YYYY-MM-DD'),
	                      COALESCE(b.customer_name, ''), COALESCE(b.customer_phone, ''),
	                      b.subtotal, b.gst_amount, b.grand_total, b.payment_method,
	                      COALESCE(b.notes, ''), b.created_by, b.created_at,
	                      COALESCE(u.name, 'Unknown')
	               FROM bills b
	               LEFT JOIN users u ON u.id = b.created_by`
	countSelect := `SELECT COUNT(*) FROM bills b`

	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(b.invoice_number ILIKE $%d OR b.customer_name ILIKE $%d)", argIdx, argIdx,
		))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}
	if filter.PaymentMethod != "" {
		conditions = append(conditions, fmt.Sprintf("b.payment_method = $%d", argIdx))
		args = append(args, filter.PaymentMethod)
		argIdx++
	}
	if filter.DateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("b.invoice_date >= $%d::date", argIdx))
		args = append(args, filter.DateFrom)
		argIdx++
	}
	if filter.DateTo != "" {
		conditions = append(conditions, fmt.Sprintf("b.invoice_date <= $%d::date", argIdx))
		args = append(args, filter.DateTo)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Count
	var total int64
	if err := r.db.QueryRow(ctx, countSelect+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count bills: %w", err)
	}

	// Paginated query
	limit := filter.PerPage
	if limit <= 0 {
		limit = 20
	}
	offset := (filter.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	query := baseSelect + where + fmt.Sprintf(
		" ORDER BY b.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query bills: %w", err)
	}
	defer rows.Close()

	var bills []domain.Bill
	for rows.Next() {
		var b domain.Bill
		if err := rows.Scan(
			&b.ID, &b.InvoiceNumber, &b.InvoiceDate,
			&b.CustomerName, &b.CustomerPhone,
			&b.Subtotal, &b.GSTAmount, &b.GrandTotal,
			&b.PaymentMethod, &b.Notes, &b.CreatedBy,
			&b.CreatedAt, &b.CreatedByName,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan bill: %w", err)
		}
		bills = append(bills, b)
	}

	return bills, total, rows.Err()
}

// ── Delete ─────────────────────────────────────────────────────────────

func (r *billRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, queryDeleteBill, id)
	if err != nil {
		return fmt.Errorf("failed to delete bill: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("bill not found")
	}
	return nil
}

// ── Invoice Number Generator ───────────────────────────────────────────

// GetNextInvoiceNumber generates the next invoice number in the format PREFIX-YYYYMM-XXXX.
// It queries the latest invoice for the current month and increments the sequence.
func (r *billRepository) GetNextInvoiceNumber(ctx context.Context, customPrefix string) (string, error) {
	now := time.Now()
	if customPrefix == "" {
		customPrefix = "INV"
	}
	prefix := fmt.Sprintf("%s-%s-", customPrefix, now.Format("200601"))
	pattern := prefix + "%"

	var lastInvoice string
	err := r.db.QueryRow(ctx, queryLatestInvoice, pattern).Scan(&lastInvoice)
	if err != nil {
		if err == pgx.ErrNoRows {
			return prefix + "0001", nil
		}
		return "", fmt.Errorf("failed to query latest invoice: %w", err)
	}

	// Parse the sequence number from "INV-YYYYMM-XXXX"
	parts := strings.Split(lastInvoice, "-")
	if len(parts) != 3 {
		return prefix + "0001", nil
	}

	seq, err := strconv.Atoi(parts[2])
	if err != nil {
		return prefix + "0001", nil
	}

	return fmt.Sprintf("%s%04d", prefix, seq+1), nil
}



