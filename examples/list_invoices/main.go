// Command list_invoices lists issued invoices for a date range, paging through
// all results, and prints the aggregate totals.
//
// Run with:
//
//	QUANTUM_API_KEY=... QUANTUM_COMPANY_ID=28218 go run ./examples/list_invoices
package main

import (
	"context"
	"log"

	"github.com/dribbaengineering/quantum-economics-sdk-go/examples/internal/exampleutil"
	"github.com/dribbaengineering/quantum-economics-sdk-go/quantum"
)

func main() {
	client, err := quantum.NewClient(
		quantum.WithAPIKey(exampleutil.MustEnv("QUANTUM_API_KEY")),
		quantum.WithCompanyID(exampleutil.MustEnvInt("QUANTUM_COMPANY_ID")),
	)
	if err != nil {
		log.Fatalf("client: %v", err)
	}

	ctx := context.Background()
	params := quantum.ListInvoicesParams{
		Type:      quantum.InvoiceTypeIssued,
		StartDate: "01-01-2025",
		EndDate:   "31-12-2025",
		Page:      1,
	}

	var total int
	for {
		page, err := client.Invoices.List(ctx, params)
		if err != nil {
			log.Fatalf("list invoices: %v", err)
		}
		for _, inv := range page.Invoices {
			total++
			log.Printf("#%s %s %.2f", inv.SeriesAndNumber, inv.InvoiceDate, inv.TotalAmount)
		}
		if params.Page >= page.TotalPages || page.TotalPages == 0 {
			log.Printf("done: %d invoices, income=%.2f expenses=%.2f balance=%.2f",
				total, page.Income, page.Expenses, page.Balance)
			break
		}
		params.Page++
	}
}
