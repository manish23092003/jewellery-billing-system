package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"jewellery-billing/internal/domain"
)

// BillService handles bill creation with server-side calculation validation.
type BillService struct {
	billRepo     domain.BillRepository
	customerRepo domain.CustomerRepository
}

func NewBillService(billRepo domain.BillRepository, customerRepo domain.CustomerRepository) *BillService {
	return &BillService{billRepo: billRepo, customerRepo: customerRepo}
}

// ── Calculation Engine ─────────────────────────────────────────────────

func calculateItemTotal(item *domain.BillItem) (float64, float64) {
	metalValue := item.Weight * item.RatePerGram
	
	chargesSum := 0.0
	for _, c := range item.Charges {
		chargesSum += c.Amount
	}

	subtotal := metalValue + item.MakingCharge + chargesSum
	gstAmount := subtotal * (item.GSTPercentage / 100.0)
	itemTotal := subtotal + gstAmount
	item.LineTotal = round2(itemTotal * float64(item.Quantity))
	return subtotal, gstAmount
}

func round2(val float64) float64 {
	return math.Round(val*100) / 100
}

// ── Create Bill ────────────────────────────────────────────────────────

func (s *BillService) Create(ctx context.Context, orgID uuid.UUID, req domain.CreateBillRequest, createdBy uuid.UUID, invoicePrefix string) (*domain.Bill, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("at least one item is required")
	}

	if !domain.ValidPaymentMethods[req.PaymentMethod] {
		req.PaymentMethod = "cash"
	}

	prefix := invoicePrefix
	if req.Type == "estimate" {
		prefix = "EST"
	}

	var items []domain.BillItem
	var totalSubtotal, totalGST, grandTotal float64

	for i, itemReq := range req.Items {
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

		var charges []domain.BillItemCharge
		for _, c := range itemReq.Charges {
			if c.ChargeName != "" && c.Amount > 0 {
				charges = append(charges, domain.BillItemCharge{
					ChargeName: c.ChargeName,
					Amount:     c.Amount,
				})
			}
		}

		item := domain.BillItem{
			ItemName:      itemReq.ItemName,
			HSNCode:       itemReq.HSNCode,
			MetalType:     itemReq.MetalType,
			Purity:        itemReq.Purity,
			Weight:        itemReq.Weight,
			RatePerGram:   itemReq.RatePerGram,
			MakingCharge:  itemReq.MakingCharge,
			GSTPercentage: itemReq.GSTPercentage,
			Quantity:      itemReq.Quantity,
			Charges:       charges,
		}

		subtotal, gstAmount := calculateItemTotal(&item)

		totalSubtotal += round2(subtotal * float64(item.Quantity))
		totalGST += round2(gstAmount * float64(item.Quantity))
		grandTotal += item.LineTotal

		items = append(items, item)
	}

	// Calculate Old Gold
	var oldGoldItems []domain.BillOldGold
	var totalOldGold float64

	for i, ogReq := range req.OldGoldItems {
		if ogReq.Weight <= 0 || ogReq.RatePerGram <= 0 {
			return nil, fmt.Errorf("old gold item %d: invalid weight or rate", i+1)
		}
		grossVal := ogReq.Weight * ogReq.RatePerGram
		deduction := grossVal * (ogReq.MeltingLossPercentage / 100.0)
		netVal := round2(grossVal - deduction)

		totalOldGold += netVal
		oldGoldItems = append(oldGoldItems, domain.BillOldGold{
			Name:                  ogReq.Name,
			Weight:                ogReq.Weight,
			Purity:                ogReq.Purity,
			MeltingLossPercentage: ogReq.MeltingLossPercentage,
			RatePerGram:           ogReq.RatePerGram,
			TotalValue:            netVal,
		})
	}

	invoiceDate := req.InvoiceDate
	if invoiceDate == "" {
		invoiceDate = "now()"
	}

	billType := "invoice"
	if req.Type == "estimate" {
		billType = "estimate"
	}
	
	status := "completed"
	if billType == "estimate" {
		status = "pending"
	}
	if req.Status != "" {
		status = req.Status
	}

	balanceDue := grandTotal - totalOldGold - req.AdvanceAmount
	if balanceDue < 0 {
		balanceDue = 0 // Or handle refunds if necessary
	}

	var bill *domain.Bill
	var saveErr error

	for attempt := 0; attempt < 3; attempt++ {
		invoiceNumber, err := s.billRepo.GetNextInvoiceNumber(ctx, orgID, prefix)
		if err != nil {
			return nil, fmt.Errorf("failed to generate invoice number: %w", err)
		}

		bill = &domain.Bill{
			OrganizationID: orgID,
			InvoiceNumber:  invoiceNumber,
			InvoiceDate:    invoiceDate,
			Type:           billType,
			Status:         status,
			CustomerName:   req.CustomerName,
			CustomerPhone:  req.CustomerPhone,
			Subtotal:       round2(totalSubtotal),
			GSTAmount:      round2(totalGST),
			OldGoldAmount:  totalOldGold,
			AdvanceAmount:  req.AdvanceAmount,
			GrandTotal:     round2(grandTotal),
			BalanceDue:     round2(balanceDue),
			PaymentMethod:  req.PaymentMethod,
			Notes:          req.Notes,
			CreatedBy:      createdBy,
			Items:              items,
			OldGoldItems:       oldGoldItems,
			VerificationToken:  uuid.New(),
			VerificationStatus: "ACTIVE",
		}

		hashDate := bill.InvoiceDate
		if hashDate == "now()" {
			hashDate = time.Now().Format("2006-01-02")
		} else if len(hashDate) >= 10 {
			hashDate = hashDate[:10]
		}

		secret := os.Getenv("VERIFICATION_SECRET")
		hashData := fmt.Sprintf("%s|%.2f|%s", bill.InvoiceNumber, bill.GrandTotal, hashDate)
		bill.InvoiceHash = generateHMACSHA256(hashData, secret)

		saveErr = s.billRepo.Create(ctx, bill)
		if saveErr == nil {
			break
		}
		if strings.Contains(saveErr.Error(), "SQLSTATE 23505") {
			// Retry on duplicate key (race condition)
			continue
		}
		break // Break on other errors
	}

	if saveErr != nil {
		return nil, fmt.Errorf("failed to save bill: %w", saveErr)
	}

	if req.ConvertFromID != "" {
		if oldID, err := uuid.Parse(req.ConvertFromID); err == nil {
			_ = s.billRepo.Delete(ctx, orgID, oldID)
		}
	}

	// === CRM Integration ===
	// Auto-create or update customer based on the bill
	if bill.CustomerName != "" && bill.Type != "estimate" {
		var cust *domain.Customer
		var err error
		
		normalizedPhone := normalizePhone(bill.CustomerPhone)
		
		if normalizedPhone != "" {
			cust, err = s.customerRepo.GetByPhone(ctx, orgID, normalizedPhone)
		}

		if err == nil && cust != nil {
			// Update existing customer
			if bill.CustomerName != "" {
				cust.Name = bill.CustomerName
			}
			cust.TotalPurchases += bill.GrandTotal
			errUpdate := s.customerRepo.Update(ctx, cust)
			if errUpdate != nil {
				fmt.Printf("Error updating customer in CRM integration: %v\n", errUpdate)
			}
		} else {
			// Create new customer
			newCust := &domain.Customer{
				OrganizationID: orgID,
				Name:           bill.CustomerName,
				Phone:          normalizedPhone, // Save normalized phone
				TotalPurchases: bill.GrandTotal,
			}
			errCreate := s.customerRepo.Create(ctx, newCust)
			if errCreate != nil {
				fmt.Printf("Error creating customer in CRM integration: %v\n", errCreate)
			}
		}
	}
	// ========================

	return bill, nil
}

// ── Read Operations ────────────────────────────────────────────────────

func (s *BillService) GetByID(ctx context.Context, orgID, id uuid.UUID) (*domain.Bill, error) {
	return s.billRepo.GetByID(ctx, orgID, id)
}

func (s *BillService) GetAll(ctx context.Context, orgID uuid.UUID, filter domain.BillFilter) ([]domain.Bill, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PerPage <= 0 {
		filter.PerPage = 20
	}
	return s.billRepo.GetAll(ctx, orgID, filter)
}

func (s *BillService) Delete(ctx context.Context, orgID, id uuid.UUID) error {
	return s.billRepo.Delete(ctx, orgID, id)
}

// ── Payment Methods ────────────────────────────────────────────────────

func (s *BillService) AddPayment(ctx context.Context, orgID, billID uuid.UUID, paymentAmount float64) error {
	bill, err := s.billRepo.GetByID(ctx, orgID, billID)
	if err != nil {
		return fmt.Errorf("failed to get bill: %w", err)
	}

	if bill.BalanceDue <= 0 {
		return fmt.Errorf("bill is already fully paid")
	}

	if paymentAmount <= 0 {
		return fmt.Errorf("payment amount must be greater than zero")
	}

	if paymentAmount > bill.BalanceDue {
		return fmt.Errorf("payment amount (%.2f) exceeds balance due (%.2f)", paymentAmount, bill.BalanceDue)
	}

	// Calculate new amounts
	newAdvanceAmount := math.Round((bill.AdvanceAmount + paymentAmount) * 100) / 100
	newBalanceDue := math.Round((bill.BalanceDue - paymentAmount) * 100) / 100

	status := bill.Status
	if newBalanceDue <= 0 && status != "cancelled" && bill.Type != "estimate" {
		status = "completed"
	} else if newBalanceDue > 0 && bill.Type != "estimate" {
		status = "pending"
	}

	return s.billRepo.UpdatePayment(ctx, orgID, billID, paymentAmount, newAdvanceAmount, newBalanceDue, status)
}

// ── Verification ───────────────────────────────────────────────────────

func generateHMACSHA256(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *BillService) VerifyInvoice(ctx context.Context, token string) (*domain.Bill, domain.PublicVerificationResponse, error) {
	bill, err := s.billRepo.GetByVerificationToken(ctx, token)
	if err != nil {
		return nil, domain.PublicVerificationResponse{}, err
	}

	secret := os.Getenv("VERIFICATION_SECRET")
	expectedData := fmt.Sprintf("%s|%.2f|%s", bill.InvoiceNumber, bill.GrandTotal, bill.InvoiceDate)
	expectedHash := generateHMACSHA256(expectedData, secret)

	status := "VERIFIED"
	if bill.InvoiceHash != expectedHash {
		status = "TAMPERED"
	} else if bill.VerificationStatus != "ACTIVE" {
		status = bill.VerificationStatus // e.g., VOID or CANCELLED
	}

	// We need the shop name. We can fetch it via setting repo, but we don't have setting repo injected here.
	// Wait, PublicVerificationResponse needs ShopName.
	// I should inject SettingRepository into BillService or fetch it inside the handler.
	// Let's return just what we have and let the handler add the ShopName, or we can just return a placeholder, or we can inject it.
	// We didn't add settingRepo to BillService yet. Let's just return what we have and the handler will do the rest.
	
	return bill, domain.PublicVerificationResponse{
		ShopName:           "", // Handled in handler
		InvoiceNumber:      bill.InvoiceNumber,
		InvoiceDate:        bill.InvoiceDate,
		GrandTotal:         bill.GrandTotal,
		BalanceDue:         bill.BalanceDue,
		VerificationStatus: status,
	}, nil
}
