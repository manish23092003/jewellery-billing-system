package service

import (
	"fmt"
	"time"

	"jewellery-billing/internal/domain"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type PDFService struct {
}

func NewPDFService() *PDFService {
	return &PDFService{}
}

// GenerateInvoicePDF creates a professional invoice PDF natively in Go using Maroto.
func (s *PDFService) GenerateInvoicePDF(bill *domain.Bill, settings *domain.ShopSettings) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithTopMargin(15).
		WithLeftMargin(15).
		WithRightMargin(15).
		Build()

	m := maroto.New(cfg)

	// Header: Shop Name & Details
	m.AddRows(
		row.New(20).Add(
			col.New(12).Add(
				text.New(settings.ShopName, props.Text{
					Top:   5,
					Style: fontstyle.Bold,
					Align: align.Center,
					Size:  20,
				}),
			),
		),
		row.New(6).Add(
			col.New(12).Add(
				text.New(settings.Address, props.Text{
					Align: align.Center,
					Size:  10,
				}),
			),
		),
		row.New(6).Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("Phone: %s | GSTIN: %s", settings.Phone, settings.GSTIN), props.Text{
					Align: align.Center,
					Size:  10,
				}),
			),
		),
	)

	// Divider
	m.AddRows(row.New(10).Add(col.New(12).Add(text.New("____________________________________________________________________________________", props.Text{Align: align.Center, Size: 10}))))

	// Invoice Info
	dateStr := bill.InvoiceDate
	if t, err := time.Parse(time.RFC3339, bill.InvoiceDate); err == nil {
		dateStr = t.Format("02-Jan-2006")
	}

	m.AddRows(
		row.New(15).Add(
			col.New(6).Add(
				text.New(fmt.Sprintf("Invoice No: %s", bill.InvoiceNumber), props.Text{Style: fontstyle.Bold, Size: 11}),
				text.New(fmt.Sprintf("Date: %s", dateStr), props.Text{Top: 5, Size: 10}),
			),
			col.New(6).Add(
				text.New(fmt.Sprintf("Customer: %s", bill.CustomerName), props.Text{Align: align.Right, Style: fontstyle.Bold, Size: 11}),
				text.New(fmt.Sprintf("Phone: %s", bill.CustomerPhone), props.Text{Align: align.Right, Top: 5, Size: 10}),
			),
		),
	)

	m.AddRows(row.New(5).Add(col.New(12).Add(text.New(""))))

	// Table Header
	m.AddRows(
		row.New(8).Add(
			col.New(4).Add(text.New("Item", props.Text{Style: fontstyle.Bold, Size: 10})),
			col.New(2).Add(text.New("Metal", props.Text{Style: fontstyle.Bold, Size: 10})),
			col.New(2).Add(text.New("Weight(g)", props.Text{Style: fontstyle.Bold, Size: 10})),
			col.New(2).Add(text.New("Rate(Rs)", props.Text{Style: fontstyle.Bold, Size: 10})),
			col.New(2).Add(text.New("Total", props.Text{Style: fontstyle.Bold, Align: align.Right, Size: 10})),
		),
	)

	// Table Rows
	for _, item := range bill.Items {
		m.AddRows(
			row.New(8).Add(
				col.New(4).Add(text.New(item.ItemName, props.Text{Size: 9})),
				col.New(2).Add(text.New(fmt.Sprintf("%s %s", item.MetalType, item.Purity), props.Text{Size: 9})),
				col.New(2).Add(text.New(fmt.Sprintf("%.3f", item.Weight), props.Text{Size: 9})),
				col.New(2).Add(text.New(fmt.Sprintf("%.2f", item.RatePerGram), props.Text{Size: 9})),
				col.New(2).Add(text.New(fmt.Sprintf("%.2f", item.LineTotal), props.Text{Align: align.Right, Size: 9})),
			),
		)
	}

	m.AddRows(row.New(5).Add(col.New(12).Add(text.New("____________________________________________________________________________________", props.Text{Align: align.Center, Size: 10}))))

	// Totals
	m.AddRows(
		row.New(8).Add(
			col.New(9).Add(text.New("Subtotal:", props.Text{Align: align.Right, Size: 10})),
			col.New(3).Add(text.New(fmt.Sprintf("Rs %.2f", bill.Subtotal), props.Text{Align: align.Right, Size: 10})),
		),
		row.New(8).Add(
			col.New(9).Add(text.New("GST Amount:", props.Text{Align: align.Right, Size: 10})),
			col.New(3).Add(text.New(fmt.Sprintf("Rs %.2f", bill.GSTAmount), props.Text{Align: align.Right, Size: 10})),
		),
		row.New(10).Add(
			col.New(9).Add(text.New("Grand Total:", props.Text{Align: align.Right, Style: fontstyle.Bold, Size: 12})),
			col.New(3).Add(text.New(fmt.Sprintf("Rs %.2f", bill.GrandTotal), props.Text{Align: align.Right, Style: fontstyle.Bold, Size: 12})),
		),
	)

	// Footer
	m.AddRows(
		row.New(30).Add(
			col.New(12).Add(
				text.New("Thank you for your business!", props.Text{Top: 15, Align: align.Center, Style: fontstyle.Italic, Size: 10}),
			),
		),
	)

	// Generate
	doc, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("could not generate PDF: %w", err)
	}

	return doc.GetBytes(), nil
}
