package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

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

	"jewellery-billing/internal/domain"
)

func main() {
	bill := &domain.Bill{
		ID:                 uuid.New(),
		InvoiceNumber:      "INV-2024-001",
		InvoiceDate:        time.Now().Format(time.RFC3339),
		Type:               "invoice",
		Status:             "completed",
		CustomerName:       "Rajesh Kumar",
		CustomerPhone:      "9876543210",
		Subtotal:           100000,
		GSTAmount:          3000,
		GrandTotal:         103000,
		AdvanceAmount:      50000,
		BalanceDue:         53000,
		VerificationToken:  uuid.New(),
		Items: []domain.BillItem{
			{
				ItemName:     "Gold Chain",
				MetalType:    "Gold",
				Purity:       "22K",
				Weight:       15.500,
				RatePerGram:  6500,
				MakingCharge: 1500,
				LineTotal:    102250,
				Charges: []domain.BillItemCharge{
					{ChargeName: "Hallmark", Amount: 200},
				},
			},
		},
	}
	settings := &domain.ShopSettings{
		ShopName: "SHREE JEWELLERS",
		Address:  "123 Main Market, MG Road, City - 400001",
		Phone:    "+91 9876543210",
		GSTIN:    "27AAAAA0000A1Z5",
	}

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

	goldColor := &props.Color{Red: 198, Green: 169, Blue: 98}
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
	dateStr := time.Now().Format("02-Jan-2006")
	m.AddRows(
		row.New(30).WithStyle(&props.Cell{BackgroundColor: lightGray}).Add(
			col.New(6).Add(
				text.New(" Billed To:", props.Text{Style: fontstyle.Bold, Size: 10, Top: 3}),
				text.New(" "+bill.CustomerName, props.Text{Top: 8, Size: 11, Style: fontstyle.Bold}),
				text.New(fmt.Sprintf(" Phone: %s", bill.CustomerPhone), props.Text{Top: 16, Size: 10}),
			),
			col.New(6).Add(
				text.New(fmt.Sprintf("Invoice No: %s ", bill.InvoiceNumber), props.Text{Align: align.Right, Style: fontstyle.Bold, Size: 11, Top: 3}),
				text.New(fmt.Sprintf("Date: %s ", dateStr), props.Text{Align: align.Right, Top: 9, Size: 10}),
			),
		),
		row.New(5), // Spacing
	)

	// TABLE HEADER
	m.AddRows(
		row.New(10).WithStyle(&props.Cell{BackgroundColor: goldColor}).Add(
			col.New(1).Add(text.New(" S.No", props.Text{Style: fontstyle.Bold, Size: 10, Color: whiteColor, Top: 2})),
			col.New(4).Add(text.New("Description", props.Text{Style: fontstyle.Bold, Size: 10, Color: whiteColor, Top: 2})),
			col.New(2).Add(text.New("Metal", props.Text{Style: fontstyle.Bold, Size: 10, Color: whiteColor, Top: 2})),
			col.New(2).Add(text.New("Wt(g)", props.Text{Style: fontstyle.Bold, Align: align.Right, Size: 10, Color: whiteColor, Top: 2})),
			col.New(1).Add(text.New("Rate", props.Text{Style: fontstyle.Bold, Align: align.Right, Size: 10, Color: whiteColor, Top: 2})),
			col.New(2).Add(text.New("Amount ", props.Text{Style: fontstyle.Bold, Align: align.Right, Size: 10, Color: whiteColor, Top: 2})),
		),
	)

	// TABLE ITEMS
	for i, item := range bill.Items {
		m.AddRows(
			row.New(8).Add(
				col.New(1).Add(text.New(fmt.Sprintf(" %d", i+1), props.Text{Size: 9, Top: 2})),
				col.New(4).Add(text.New(item.ItemName, props.Text{Size: 9, Style: fontstyle.Bold, Top: 2})),
				col.New(2).Add(text.New(fmt.Sprintf("%s %s", item.MetalType, item.Purity), props.Text{Size: 9, Top: 2})),
				col.New(2).Add(text.New(fmt.Sprintf("%.3f", item.Weight), props.Text{Align: align.Right, Size: 9, Top: 2})),
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

	// TOTALS & TERMS
	m.AddRows(
		row.New(45).Add(
			// Terms
			col.New(7).Add(
				text.New("Terms & Conditions:", props.Text{Style: fontstyle.Bold, Size: 9, Color: goldColor}),
				text.New("1. Goods once sold will not be taken back.", props.Text{Top: 6, Size: 8}),
				text.New("2. Subject to local jurisdiction only.", props.Text{Top: 12, Size: 8}),
				text.New("3. Making charges and taxes apply as per government norms.", props.Text{Top: 18, Size: 8}),
			),
			// Totals
			col.New(5).Add(
				text.New(fmt.Sprintf("Subtotal: Rs %.2f", bill.Subtotal), props.Text{Align: align.Right, Size: 10}),
				text.New(fmt.Sprintf("GST Amount: Rs %.2f", bill.GSTAmount), props.Text{Top: 7, Align: align.Right, Size: 10}),
				text.New(fmt.Sprintf("Grand Total: Rs %.2f", bill.GrandTotal), props.Text{Top: 16, Align: align.Right, Style: fontstyle.Bold, Size: 12, Color: goldColor}),
				text.New(fmt.Sprintf("Advance/Paid: Rs %.2f", bill.AdvanceAmount), props.Text{Top: 24, Align: align.Right, Size: 10, Color: &props.Color{Red: 0, Green: 128, Blue: 0}}),
				text.New(fmt.Sprintf("Balance Due: Rs %.2f", bill.BalanceDue), props.Text{Top: 32, Align: align.Right, Style: fontstyle.Bold, Size: 11, Color: &props.Color{Red: 200, Green: 0, Blue: 0}}),
			),
		),
	)

	m.AddRows(row.New(10).Add(col.New(12).Add(line.New(props.Line{Thickness: 0.5, Color: goldColor}))))

	// FOOTER & QR CODE
	verificationURL := fmt.Sprintf("http://localhost:5173/verify/%s", bill.VerificationToken.String())

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
				text.New(verificationURL, props.Text{Top: 24, Align: align.Left, Size: 7}),
			),
			// Right side: Signature
			col.New(4).Add(
				text.New(fmt.Sprintf("For %s", settings.ShopName), props.Text{Top: 15, Align: align.Right, Style: fontstyle.Bold, Size: 10}),
				text.New("Authorized Signatory", props.Text{Top: 28, Align: align.Right, Size: 9, Style: fontstyle.Italic}),
			),
		),
	)

	doc, err := m.Generate()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	os.WriteFile("test_invoice.pdf", doc.GetBytes(), 0644)
	fmt.Println("Saved test_invoice.pdf")
}
