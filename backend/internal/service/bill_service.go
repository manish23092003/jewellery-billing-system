package service

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"

	"jewellery-billing/internal/domain"
)

// BillService handles bill creation with server-side calculation validation.
type BillService struct {
	billRepo domain.BillRepository
}

func NewBillService(billRepo domain.BillRepository) *BillService {
	return &BillService{billRepo: billRepo}
}

// ── Calculation Engine ─────────────────────────────────────────────────

// calculateItemTotal computes a single line item's total using the formula:
//
//	metal_value  = weight × rate_per_gram
//	subtotal     = metal_value + making_charge
//	gst_amount   = subtotal × (gst_percentage / 100)
//	item_total   = subtotal + gst_amount
//	line_total   = item_total × quantity
func calculateItemTotal(item *domain.BillItem) {
	metalValue := item.Weight * item.RatePerGram
	subtotal := metalValue + item.MakingCharge
	gstAmount := subtotal * (item.GSTPercentage / 100.0)
	itemTotal := subtotal + gstAmount
	item.LineTotal = round2(itemTotal * float64(item.Quantity))
}

// round2 rounds a float to 2 decimal places (standard currency rounding).
func round2(val float64) float64 {
	return math.Round(val*100) / 100
}

// ── Create Bill ────────────────────────────────────────────────────────

func (s *BillService) Create(ctx context.Context, req domain.CreateBillRequest, createdBy uuid.UUID, invoicePrefix string) (*domain.Bill, error) {
	// Validate items
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("at least one item is required")
	}

	// Validate payment method
	if !domain.ValidPaymentMethods[req.PaymentMethod] {
		req.PaymentMethod = "cash"
	}

	// Generate invoice number
	invoiceNumber, err := s.billRepo.GetNextInvoiceNumber(ctx, invoicePrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice number: %w", err)
	}

	// Build and calculate items
	var items []domain.BillItem
	var totalSubtotal, totalGST, grandTotal float64

	for i, itemReq := range req.Items {
		// Validate each item
		if itemReq.ItemName == "" {
			return nil, fmt.Errorf("item %d: name is required", i+1)
		}
		if itemReq.Weight <= 0 {
			return nil, fmt.Errorf("item %d: weight must be positive", i+1)
		}
		if itemReq.RatePerGram <= 0 {
			return nil, fmt.Errorf("item %d: rate per gram must be positive", i+1)
		}
		if itemReq.Quantity <= 0 {
			itemReq.Quantity = 1
		}
		if itemReq.GSTPercentage < 0 {
			return nil, fmt.Errorf("item %d: GST percentage cannot be negative", i+1)
		}

		item := domain.BillItem{
			ItemName:      itemReq.ItemName,
			MetalType:     itemReq.MetalType,
			Purity:        itemReq.Purity,
			Weight:        itemReq.Weight,
			RatePerGram:   itemReq.RatePerGram,
			MakingCharge:  itemReq.MakingCharge,
			GSTPercentage: itemReq.GSTPercentage,
			Quantity:      itemReq.Quantity,
		}

		// Server-side calculation — overrides any client-supplied total
		calculateItemTotal(&item)

		// Accumulate invoice totals
		metalValue := item.Weight * item.RatePerGram
		subtotal := metalValue + item.MakingCharge
		gstAmount := subtotal * (item.GSTPercentage / 100.0)

		totalSubtotal += round2(subtotal * float64(item.Quantity))
		totalGST += round2(gstAmount * float64(item.Quantity))
		grandTotal += item.LineTotal

		items = append(items, item)
	}

	// Set invoice date (default to today if empty)
	invoiceDate := req.InvoiceDate
	if invoiceDate == "" {
		invoiceDate = "now()"
	}

	bill := &domain.Bill{
		InvoiceNumber: invoiceNumber,
		InvoiceDate:   invoiceDate,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		Subtotal:      round2(totalSubtotal),
		GSTAmount:     round2(totalGST),
		GrandTotal:    round2(grandTotal),
		PaymentMethod: req.PaymentMethod,
		Notes:         req.Notes,
		CreatedBy:     createdBy,
		Items:         items,
	}

	if err := s.billRepo.Create(ctx, bill); err != nil {
		return nil, fmt.Errorf("failed to create bill: %w", err)
	}

	return bill, nil
}

// ── Read Operations ────────────────────────────────────────────────────

func (s *BillService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Bill, error) {
	return s.billRepo.GetByID(ctx, id)
}

func (s *BillService) GetAll(ctx context.Context, filter domain.BillFilter) ([]domain.Bill, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PerPage <= 0 {
		filter.PerPage = 20
	}
	return s.billRepo.GetAll(ctx, filter)
}

func (s *BillService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.billRepo.Delete(ctx, id)
}



