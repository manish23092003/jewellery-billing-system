package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"jewellery-billing/internal/domain"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type PDFService struct {
}

func NewPDFService() *PDFService {
	return &PDFService{}
}

// GenerateInvoicePDF creates a professional invoice PDF natively in Go using Maroto.
func (s *PDFService) GenerateInvoicePDF(bill *domain.Bill, settings *domain.ShopSettings) ([]byte, error) {
	builder := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithTopMargin(15).
		WithLeftMargin(15).
		WithRightMargin(15)

	m := maroto.New(builder.Build())

	docTitle := "TAX INVOICE"
	if bill.Type == "estimate" {
		docTitle = "ESTIMATE"
	}

	// Brand Colors
	goldColor := &props.Color{Red: 198, Green: 169, Blue: 98} // #c6a962
	whiteColor := &props.Color{Red: 255, Green: 255, Blue: 255}
	lightGray := &props.Color{Red: 245, Green: 245, Blue: 245}

	// HEADER ROW
	m.AddRows(
		row.New(45).Add(
			col.New(8).Add(
				text.New(settings.ShopName, props.Text{
					Style: fontstyle.Bold,
					Size:  22,
					Align: align.Left,
					Color: goldColor,
				}),
				text.New(settings.Address, props.Text{
					Top:   25,
					Size:  10,
					Align: align.Left,
				}),
				text.New(fmt.Sprintf("Phone: %s | GSTIN: %s", settings.Phone, settings.GSTIN), props.Text{
					Top:   32,
					Size:  10,
					Align: align.Left,
				}),
			),
			col.New(4).Add(
				text.New(docTitle, props.Text{
					Style: fontstyle.Bold,
					Size:  20,
					Align: align.Right,
					Color: goldColor,
				}),
			),
		),
		row.New(5).Add(col.New(12).Add(line.New(props.Line{Thickness: 0.5, Color: goldColor}))),
	)

	// INVOICE DETAILS
	dateStr := bill.InvoiceDate
	if t, err := time.Parse(time.RFC3339, bill.InvoiceDate); err == nil {
		dateStr = t.Format("02-Jan-2006")
	}

	m.AddRows(
		row.New(30).WithStyle(&props.Cell{BackgroundColor: lightGray}).Add(
			col.New(6).Add(
				text.New(" Billed To:", props.Text{Style: fontstyle.Bold, Size: 10, Top: 3}),
				text.New(" "+bill.CustomerName, props.Text{Top: 8, Size: 11, Style: fontstyle.Bold}),
				text.New(fmt.Sprintf(" Phone: %s", bill.CustomerPhone), props.Text{Top: 16, Size: 10}),
			),
			col.New(6).Add(
				text.New(fmt.Sprintf("Invoice No: %s ", bill.InvoiceNumber), props.Text{Align: align.Right, Style: fontstyle.Bold, Size: 11, Top: 3}),
				text.New(fmt.Sprintf("Date: %s ", dateStr), props.Text{Align: align.Right, Top: 10, Size: 10}),
				text.New(fmt.Sprintf("Payment Mode: %s ", strings.ToUpper(bill.PaymentMethod)), props.Text{Align: align.Right, Top: 17, Size: 10}),
			),
		),
		row.New(5), // Spacing
	)

	// TABLE HEADER
	m.AddRows(
		row.New(10).WithStyle(&props.Cell{BackgroundColor: goldColor}).Add(
			col.New(1).Add(text.New(" S.No", props.Text{Style: fontstyle.Bold, Size: 10, Color: whiteColor, Top: 2})),
			col.New(2).Add(text.New("Description", props.Text{Style: fontstyle.Bold, Size: 10, Color: whiteColor, Top: 2})),
			col.New(2).Add(text.New("HSN", props.Text{Style: fontstyle.Bold, Size: 10, Color: whiteColor, Top: 2})),
			col.New(2).Add(text.New("Metal", props.Text{Style: fontstyle.Bold, Size: 10, Color: whiteColor, Top: 2})),
			col.New(1).Add(text.New("Wt", props.Text{Style: fontstyle.Bold, Align: align.Right, Size: 10, Color: whiteColor, Top: 2})),
			col.New(1).Add(text.New("Qty", props.Text{Style: fontstyle.Bold, Align: align.Right, Size: 10, Color: whiteColor, Top: 2})),
			col.New(1).Add(text.New("Rate", props.Text{Style: fontstyle.Bold, Align: align.Right, Size: 10, Color: whiteColor, Top: 2})),
			col.New(2).Add(text.New("Amount ", props.Text{Style: fontstyle.Bold, Align: align.Right, Size: 10, Color: whiteColor, Top: 2})),
		),
	)

	// TABLE ITEMS
	for i, item := range bill.Items {
		m.AddRows(
			row.New(8).Add(
				col.New(1).Add(text.New(fmt.Sprintf(" %d", i+1), props.Text{Size: 9, Top: 2})),
				col.New(2).Add(text.New(item.ItemName, props.Text{Size: 9, Style: fontstyle.Bold, Top: 2})),
				col.New(2).Add(text.New(item.HSNCode, props.Text{Size: 9, Top: 2})),
				col.New(2).Add(text.New(fmt.Sprintf("%s %s", item.MetalType, item.Purity), props.Text{Size: 9, Top: 2})),
				col.New(1).Add(text.New(fmt.Sprintf("%.3f", item.Weight), props.Text{Align: align.Right, Size: 9, Top: 2})),
				col.New(1).Add(text.New(fmt.Sprintf("%d", item.Quantity), props.Text{Align: align.Right, Size: 9, Top: 2})),
				col.New(1).Add(text.New(fmt.Sprintf("%.0f", item.RatePerGram), props.Text{Align: align.Right, Size: 9, Top: 2})),
				col.New(2).Add(text.New(fmt.Sprintf("%.2f ", item.LineTotal), props.Text{Align: align.Right, Size: 9, Top: 2})),
			),
		)
		
		if item.MakingCharge > 0 {
			m.AddRows(
				row.New(6).Add(
					col.New(1),
					col.New(11).Add(text.New(fmt.Sprintf("+ Making Charge: Rs %.2f", item.MakingCharge), props.Text{Size: 8, Style: fontstyle.Italic, Color: &props.Color{Red: 100, Green: 100, Blue: 100}})),
				),
			)
		}
		for _, c := range item.Charges {
			m.AddRows(
				row.New(6).Add(
					col.New(1),
					col.New(11).Add(text.New(fmt.Sprintf("+ %s: Rs %.2f", c.ChargeName, c.Amount), props.Text{Size: 8, Style: fontstyle.Italic, Color: &props.Color{Red: 100, Green: 100, Blue: 100}})),
				),
			)
		}
		m.AddRows(row.New(3).Add(col.New(12).Add(line.New(props.Line{Thickness: 0.1, Color: lightGray}))))
	}

	m.AddRows(row.New(5))

	// OLD GOLD DEDUCTIONS
	if len(bill.OldGoldItems) > 0 {
		m.AddRows(row.New(8).Add(col.New(12).Add(text.New("Old Gold Deductions", props.Text{Style: fontstyle.Bold, Size: 10, Color: goldColor}))))
		
		for _, og := range bill.OldGoldItems {
			m.AddRows(
				row.New(6).Add(
					col.New(6).Add(text.New(fmt.Sprintf("%s (%s) - %.3fg", og.Name, og.Purity, og.Weight), props.Text{Size: 9})),
					col.New(3).Add(text.New(fmt.Sprintf("@ Rs %.2f/g", og.RatePerGram), props.Text{Size: 9})),
					col.New(3).Add(text.New(fmt.Sprintf("- Rs %.2f", og.TotalValue), props.Text{Align: align.Right, Size: 9})),
				),
			)
		}
		m.AddRows(row.New(5).Add(col.New(12).Add(line.New(props.Line{Thickness: 0.5, Color: goldColor}))))
	}

	// TOTALS & TERMS
	m.AddRows(
		row.New(6).Add(
			col.New(7), col.New(5).Add(text.New(fmt.Sprintf("Subtotal: Rs %.2f", bill.Subtotal), props.Text{Align: align.Right, Size: 10})),
		),
		row.New(6).Add(
			col.New(7), col.New(5).Add(text.New(fmt.Sprintf("CGST: Rs %.2f", bill.GSTAmount/2), props.Text{Align: align.Right, Size: 10})),
		),
		row.New(6).Add(
			col.New(7), col.New(5).Add(text.New(fmt.Sprintf("SGST: Rs %.2f", bill.GSTAmount/2), props.Text{Align: align.Right, Size: 10})),
		),
		row.New(4).Add(
			col.New(7), col.New(5).Add(line.New(props.Line{Thickness: 0.5, Color: goldColor})),
		),
		row.New(8).Add(
			col.New(7), col.New(5).Add(text.New(fmt.Sprintf("Grand Total: Rs %.2f", bill.GrandTotal), props.Text{Top: 2, Align: align.Right, Style: fontstyle.Bold, Size: 12, Color: goldColor})),
		),
	)

	if len(bill.Payments) > 0 {
		for _, pay := range bill.Payments {
			pDate := pay.PaymentDate.Format("02-Jan-2006")
			m.AddRows(
				row.New(6).Add(
					col.New(6), col.New(6).Add(text.New(fmt.Sprintf("Paid (%s): Rs %.2f", pDate, pay.Amount), props.Text{Top: 2, Align: align.Right, Size: 10, Color: &props.Color{Red: 0, Green: 128, Blue: 0}})),
				),
			)
		}
	} else {
		m.AddRows(
			row.New(6).Add(
				col.New(7), col.New(5).Add(text.New(fmt.Sprintf("Advance/Paid: Rs %.2f", bill.AdvanceAmount), props.Text{Top: 2, Align: align.Right, Size: 10, Color: &props.Color{Red: 0, Green: 128, Blue: 0}})),
			),
		)
	}

	m.AddRows(
		row.New(8).Add(
			col.New(7), col.New(5).Add(text.New(fmt.Sprintf("Balance Due: Rs %.2f", bill.BalanceDue), props.Text{Top: 2, Align: align.Right, Style: fontstyle.Bold, Size: 11, Color: &props.Color{Red: 200, Green: 0, Blue: 0}})),
		),
	)

	m.AddRows(row.New(10).Add(col.New(12).Add(line.New(props.Line{Thickness: 0.5, Color: goldColor}))))

	// FOOTER & QR CODE
	appURL := os.Getenv("FRONTEND_URL")
	if appURL == "" || strings.Contains(appURL, "localhost") || strings.Contains(appURL, "onrender.com") {
		appURL = "https://jewellery-billing-system-psi.vercel.app"
	}
	verificationURL := fmt.Sprintf("%s/#/verify/%s", appURL, bill.VerificationToken.String())

	m.AddRows(
		row.New(40).Add(
			// Left side: QR Code
			col.New(3).Add(
				code.NewQr(verificationURL, props.Rect{Center: true, Percent: 100}),
			),
			// Middle side: Verification text
			col.New(5).Add(
				text.New("Thank you for your business!", props.Text{Top: 5, Align: align.Left, Style: fontstyle.Italic, Size: 10, Color: goldColor}),
				text.New("Scan QR code to verify this authentic invoice.", props.Text{Top: 12, Align: align.Left, Size: 9}),
				text.New("Secured & Verified", props.Text{Top: 18, Align: align.Left, Style: fontstyle.Bold, Size: 9, Color: &props.Color{Red: 34, Green: 139, Blue: 34}}),
			),
			// Right side: Signature (removed)
			col.New(4),
		),
	)

	// Gap
	m.AddRows(row.New(15))

	// Terms & Conditions
	m.AddRows(
		row.New(4).Add(col.New(12).Add(text.New("Terms & Conditions:", props.Text{Style: fontstyle.Bold, Size: 8}))),
		row.New(4).Add(col.New(12).Add(text.New("1. Goods once sold will not be taken back or exchanged.", props.Text{Size: 7}))),
		row.New(4).Add(col.New(12).Add(text.New("2. All disputes are subject to local jurisdiction.", props.Text{Size: 7}))),
		row.New(4).Add(col.New(12).Add(text.New("3. Taxes and making charges are applicable as per government norms.", props.Text{Size: 7}))),
		row.New(4).Add(col.New(12).Add(text.New("4. Please retain this invoice for future reference.", props.Text{Size: 7}))),
		row.New(8).Add(col.New(12).Add(text.New("5. E.& O.E. (Errors and Omissions Excepted).", props.Text{Size: 7}))),
	)

	// Generate
	doc, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("could not generate PDF: %w", err)
	}

	return doc.GetBytes(), nil
}
