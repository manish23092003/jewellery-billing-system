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
	queryInsertBillMultiTenant = `
		INSERT INTO bills (organization_id, invoice_number, invoice_date, type, status, customer_name, customer_phone,
		                   subtotal, gst_amount, old_gold_amount, advance_amount, grand_total, balance_due, payment_method, notes, created_by,
						   verification_token, invoice_hash, verification_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id, created_at`

	queryInsertBillItemMultiTenant = `
		INSERT INTO bill_items (organization_id, bill_id, item_name, hsn_code, metal_type, purity, weight,
		                        rate_per_gram, making_charge, gst_percentage, quantity, line_total)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at`

	queryGetBillByIDMultiTenant = `
		SELECT b.id, b.organization_id, b.invoice_number, TO_CHAR(b.invoice_date, 'YYYY-MM-DD'),
		       b.type, b.status,
		       COALESCE(b.customer_name, ''), COALESCE(b.customer_phone, ''),
		       b.subtotal, b.gst_amount, b.old_gold_amount, b.advance_amount,
		       b.grand_total, b.balance_due, b.payment_method,
		       COALESCE(b.notes, ''), b.created_by, b.created_at,
		       COALESCE(u.name, 'Unknown'),
		       b.verification_token, COALESCE(b.invoice_hash, ''), COALESCE(b.verification_status, 'ACTIVE')
		FROM bills b
		LEFT JOIN users u ON u.id = b.created_by
		WHERE b.id = $1 AND b.organization_id = $2`

	queryGetBillByVerificationToken = `
		SELECT b.id, b.organization_id, b.invoice_number, TO_CHAR(b.invoice_date, 'YYYY-MM-DD'),
		       b.type, b.status,
		       COALESCE(b.customer_name, ''), COALESCE(b.customer_phone, ''),
		       b.subtotal, b.gst_amount, b.old_gold_amount, b.advance_amount,
		       b.grand_total, b.balance_due, b.payment_method,
		       COALESCE(b.notes, ''), b.created_by, b.created_at,
		       COALESCE(u.name, 'Unknown'),
		       b.verification_token, COALESCE(b.invoice_hash, ''), COALESCE(b.verification_status, 'ACTIVE')
		FROM bills b
		LEFT JOIN users u ON u.id = b.created_by
		WHERE b.verification_token = $1`

	queryGetBillItemsMultiTenant = `
		SELECT id, bill_id, item_name, COALESCE(hsn_code, ''), metal_type, purity, weight,
		       rate_per_gram, making_charge, gst_percentage, quantity, line_total, created_at
		FROM bill_items
		WHERE bill_id = $1
		ORDER BY created_at`

	queryDeleteBillMultiTenant = `DELETE FROM bills WHERE id = $1 AND organization_id = $2`

	queryLatestInvoiceMultiTenant = `
		SELECT invoice_number FROM bills
		WHERE organization_id = $1 AND invoice_number LIKE $2
		ORDER BY invoice_number DESC
		LIMIT 1`

	queryInsertBillItemCharge = `
		INSERT INTO bill_item_charges (bill_item_id, charge_name, amount)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	queryGetBillItemCharges = `
		SELECT id, bill_item_id, charge_name, amount, created_at
		FROM bill_item_charges
		WHERE bill_item_id = $1
		ORDER BY created_at`

	queryInsertBillOldGold = `
		INSERT INTO bill_old_gold (bill_id, name, weight, purity, melting_loss_percentage, rate_per_gram, total_value)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	queryGetBillOldGold = `
		SELECT id, bill_id, name, weight, purity, melting_loss_percentage, rate_per_gram, total_value, created_at
		FROM bill_old_gold
		WHERE bill_id = $1
		ORDER BY created_at`

	queryGetBillPayments = `
		SELECT id, organization_id, bill_id, amount, payment_date
		FROM bill_payments
		WHERE bill_id = $1
		ORDER BY payment_date ASC`
)

// ── Create (Transactional) ─────────────────────────────────────────────

func (r *billRepository) Create(ctx context.Context, bill *domain.Bill) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert bill header
	err = tx.QueryRow(ctx, queryInsertBillMultiTenant,
		bill.OrganizationID, bill.InvoiceNumber, bill.InvoiceDate, bill.Type, bill.Status,
		bill.CustomerName, bill.CustomerPhone,
		bill.Subtotal, bill.GSTAmount, bill.OldGoldAmount, bill.AdvanceAmount,
		bill.GrandTotal, bill.BalanceDue, bill.PaymentMethod, bill.Notes, bill.CreatedBy,
		bill.VerificationToken, bill.InvoiceHash, bill.VerificationStatus,
	).Scan(&bill.ID, &bill.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert bill: %w", err)
	}

	// Insert initial advance amount into bill_payments if > 0
	if bill.AdvanceAmount > 0 {
		insertPaymentQuery := `
			INSERT INTO bill_payments (organization_id, bill_id, amount)
			VALUES ($1, $2, $3)
		`
		_, err = tx.Exec(ctx, insertPaymentQuery, bill.OrganizationID, bill.ID, bill.AdvanceAmount)
		if err != nil {
			return fmt.Errorf("failed to insert initial payment: %w", err)
		}
	}

	// Insert each line item
	for i := range bill.Items {
		bill.Items[i].BillID = bill.ID
		err = tx.QueryRow(ctx, queryInsertBillItemMultiTenant,
			bill.OrganizationID, bill.Items[i].BillID, bill.Items[i].ItemName, bill.Items[i].HSNCode,
			bill.Items[i].MetalType, bill.Items[i].Purity,
			bill.Items[i].Weight, bill.Items[i].RatePerGram,
			bill.Items[i].MakingCharge, bill.Items[i].GSTPercentage,
			bill.Items[i].Quantity, bill.Items[i].LineTotal,
		).Scan(&bill.Items[i].ID, &bill.Items[i].CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert bill item %d: %w", i+1, err)
		}

		for j := range bill.Items[i].Charges {
			bill.Items[i].Charges[j].BillItemID = bill.Items[i].ID
			err = tx.QueryRow(ctx, queryInsertBillItemCharge,
				bill.Items[i].Charges[j].BillItemID,
				bill.Items[i].Charges[j].ChargeName,
				bill.Items[i].Charges[j].Amount,
			).Scan(&bill.Items[i].Charges[j].ID, &bill.Items[i].Charges[j].CreatedAt)
			if err != nil {
				return fmt.Errorf("failed to insert charge %d for item %d: %w", j+1, i+1, err)
			}
		}
	}

	// Insert Old Gold Items
	for i := range bill.OldGoldItems {
		bill.OldGoldItems[i].BillID = bill.ID
		err = tx.QueryRow(ctx, queryInsertBillOldGold,
			bill.OldGoldItems[i].BillID, bill.OldGoldItems[i].Name,
			bill.OldGoldItems[i].Weight, bill.OldGoldItems[i].Purity,
			bill.OldGoldItems[i].MeltingLossPercentage, bill.OldGoldItems[i].RatePerGram,
			bill.OldGoldItems[i].TotalValue,
		).Scan(&bill.OldGoldItems[i].ID, &bill.OldGoldItems[i].CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert old gold item %d: %w", i+1, err)
		}
	}

	return tx.Commit(ctx)
}

// ── Get By ID ──────────────────────────────────────────────────────────

func (r *billRepository) GetByID(ctx context.Context, orgID, id uuid.UUID) (*domain.Bill, error) {
	bill := &domain.Bill{}
	err := r.db.QueryRow(ctx, queryGetBillByIDMultiTenant, id, orgID).Scan(
		&bill.ID, &bill.OrganizationID, &bill.InvoiceNumber, &bill.InvoiceDate,
		&bill.Type, &bill.Status,
		&bill.CustomerName, &bill.CustomerPhone,
		&bill.Subtotal, &bill.GSTAmount, &bill.OldGoldAmount, &bill.AdvanceAmount,
		&bill.GrandTotal, &bill.BalanceDue, &bill.PaymentMethod,
		&bill.Notes, &bill.CreatedBy,
		&bill.CreatedAt, &bill.CreatedByName,
		&bill.VerificationToken, &bill.InvoiceHash, &bill.VerificationStatus,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("bill not found")
		}
		return nil, fmt.Errorf("failed to get bill: %w", err)
	}

	// Fetch line items
	rows, err := r.db.Query(ctx, queryGetBillItemsMultiTenant, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bill items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.BillItem
		if err := rows.Scan(
			&item.ID, &item.BillID, &item.ItemName, &item.HSNCode, &item.MetalType,
			&item.Purity, &item.Weight, &item.RatePerGram,
			&item.MakingCharge, &item.GSTPercentage, &item.Quantity,
			&item.LineTotal, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan bill item: %w", err)
		}
		bill.Items = append(bill.Items, item)
	}
	rows.Close() // Close item rows before running nested queries

	// Fetch charges for each item
	for i := range bill.Items {
		cRows, err := r.db.Query(ctx, queryGetBillItemCharges, bill.Items[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get charges for item %s: %w", bill.Items[i].ID, err)
		}

		for cRows.Next() {
			var charge domain.BillItemCharge
			if err := cRows.Scan(
				&charge.ID, &charge.BillItemID, &charge.ChargeName,
				&charge.Amount, &charge.CreatedAt,
			); err != nil {
				cRows.Close()
				return nil, fmt.Errorf("failed to scan charge: %w", err)
			}
			bill.Items[i].Charges = append(bill.Items[i].Charges, charge)
		}
		cRows.Close()
	}

	// Fetch Old Gold items
	ogRows, err := r.db.Query(ctx, queryGetBillOldGold, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get old gold items: %w", err)
	}

	for ogRows.Next() {
		var og domain.BillOldGold
		if err := ogRows.Scan(
			&og.ID, &og.BillID, &og.Name, &og.Weight, &og.Purity,
			&og.MeltingLossPercentage, &og.RatePerGram, &og.TotalValue, &og.CreatedAt,
		); err != nil {
			ogRows.Close()
			return nil, fmt.Errorf("failed to scan old gold item: %w", err)
		}
		bill.OldGoldItems = append(bill.OldGoldItems, og)
	}
	ogRows.Close()

	// Fetch Payments
	payRows, err := r.db.Query(ctx, queryGetBillPayments, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bill payments: %w", err)
	}

	for payRows.Next() {
		var pay domain.BillPayment
		if err := payRows.Scan(
			&pay.ID, &pay.OrganizationID, &pay.BillID, &pay.Amount, &pay.PaymentDate,
		); err != nil {
			payRows.Close()
			return nil, fmt.Errorf("failed to scan bill payment: %w", err)
		}
		bill.Payments = append(bill.Payments, pay)
	}
	payRows.Close()

	return bill, nil
}

// ── Get By Verification Token ──────────────────────────────────────────

func (r *billRepository) GetByVerificationToken(ctx context.Context, token string) (*domain.Bill, error) {
	bill := &domain.Bill{}
	err := r.db.QueryRow(ctx, queryGetBillByVerificationToken, token).Scan(
		&bill.ID, &bill.OrganizationID, &bill.InvoiceNumber, &bill.InvoiceDate,
		&bill.Type, &bill.Status,
		&bill.CustomerName, &bill.CustomerPhone,
		&bill.Subtotal, &bill.GSTAmount, &bill.OldGoldAmount, &bill.AdvanceAmount,
		&bill.GrandTotal, &bill.BalanceDue, &bill.PaymentMethod,
		&bill.Notes, &bill.CreatedBy,
		&bill.CreatedAt, &bill.CreatedByName,
		&bill.VerificationToken, &bill.InvoiceHash, &bill.VerificationStatus,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invalid verification token")
		}
		return nil, fmt.Errorf("failed to get bill by verification token: %w", err)
	}

	// We don't fetch full items here unless needed, but for validation we only need GrandTotal, InvoiceNumber, InvoiceDate, etc.
	// Since the requirement is just Shop Name, Invoice Number, Date, Grand Total, and Status, we don't need items.
	// But let's fetch shop name instead of createdByName for this purpose if needed, or we can leave it as is 
	// and fetch settings in the service layer using organization_id.

	return bill, nil
}

// ── Get All (Paginated + Filtered) ─────────────────────────────────────

func (r *billRepository) GetAll(ctx context.Context, orgID uuid.UUID, filter domain.BillFilter) ([]domain.Bill, int64, error) {
	baseSelect := `SELECT b.id, b.organization_id, b.invoice_number, TO_CHAR(b.invoice_date, 'YYYY-MM-DD'),
	                      b.type, b.status,
	                      COALESCE(b.customer_name, ''), COALESCE(b.customer_phone, ''),
	                      b.subtotal, b.gst_amount, b.old_gold_amount, b.advance_amount,
	                      b.grand_total, b.balance_due, b.payment_method,
	                      COALESCE(b.notes, ''), b.created_by, b.created_at,
	                      COALESCE(u.name, 'Unknown')
	               FROM bills b
	               LEFT JOIN users u ON u.id = b.created_by`
	countSelect := `SELECT COUNT(*) FROM bills b`

	// Always filter by organization
	conditions := []string{fmt.Sprintf("b.organization_id = $1")}
	args := []interface{}{orgID}
	argIdx := 2

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(b.invoice_number ILIKE $%d OR b.customer_name ILIKE $%d)", argIdx, argIdx,
		))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}
	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("b.type = $%d", argIdx))
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("b.status = $%d", argIdx))
		args = append(args, filter.Status)
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

	where := " WHERE " + strings.Join(conditions, " AND ")

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
			&b.ID, &b.OrganizationID, &b.InvoiceNumber, &b.InvoiceDate,
			&b.Type, &b.Status,
			&b.CustomerName, &b.CustomerPhone,
			&b.Subtotal, &b.GSTAmount, &b.OldGoldAmount, &b.AdvanceAmount,
			&b.GrandTotal, &b.BalanceDue, &b.PaymentMethod,
			&b.Notes, &b.CreatedBy,
			&b.CreatedAt, &b.CreatedByName,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan bill: %w", err)
		}
		bills = append(bills, b)
	}

	return bills, total, rows.Err()
}

// ── Delete ─────────────────────────────────────────────────────────────

func (r *billRepository) Delete(ctx context.Context, orgID, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, queryDeleteBillMultiTenant, id, orgID)
	if err != nil {
		return fmt.Errorf("failed to delete bill: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("bill not found")
	}
	return nil
}

// ── Invoice Number Generator ───────────────────────────────────────────

func (r *billRepository) GetNextInvoiceNumber(ctx context.Context, orgID uuid.UUID, customPrefix string) (string, error) {
	now := time.Now()
	if customPrefix == "" {
		customPrefix = "INV"
	}
	prefix := fmt.Sprintf("%s-%s-", customPrefix, now.Format("200601"))
	pattern := prefix + "%"

	var lastInvoice string
	err := r.db.QueryRow(ctx, queryLatestInvoiceMultiTenant, orgID, pattern).Scan(&lastInvoice)
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

func (r *billRepository) UpdatePayment(ctx context.Context, orgID, billID uuid.UUID, paymentAmount, advanceAmount, balanceDue float64, status string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update bill
	updateQuery := `
		UPDATE bills
		SET advance_amount = $1, balance_due = $2, status = $3
		WHERE organization_id = $4 AND id = $5
	`
	tag, err := tx.Exec(ctx, updateQuery, advanceAmount, balanceDue, status, orgID, billID)
	if err != nil {
		return fmt.Errorf("failed to update bill payment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bill not found or no permission")
	}

	// Insert payment record
	insertQuery := `
		INSERT INTO bill_payments (organization_id, bill_id, amount)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, insertQuery, orgID, billID, paymentAmount)
	if err != nil {
		return fmt.Errorf("failed to insert payment record: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit payment transaction: %w", err)
	}

	return nil
}
