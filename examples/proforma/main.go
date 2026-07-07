// Command proforma creates a pro forma invoice and downloads its PDF.
//
// Run with:
//
//	QUANTUM_API_KEY=... QUANTUM_COMPANY_ID=28218 go run ./examples/proforma
package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"strconv"

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

	res, err := client.Proforma.Create(ctx, quantum.ProformaInvoice{
		CustomerProviderID:      2,
		TotalAmountWithoutTaxes: 100.00,
		TotalAmount:             121.00,
		Line: []quantum.InvoiceLine{{
			Description: "Order 1",
			Quantity:    1,
			Amount:      121,
			Base:        100,
			Percentage:  21,
			TaxCode:     21,
		}},
		Installment: []quantum.InvoiceInstallment{{Date: "22/09/2025", Amount: 121.00}},
	})
	if err != nil {
		log.Fatalf("create proforma: %v", err)
	}
	log.Printf("created pro forma invoice, id=%d", res.ID)

	// Download the PDF (Base64) and save it to disk.
	doc, err := client.Proforma.Document(ctx, quantum.DocumentParams{
		ID:       strconv.FormatInt(res.ID, 10),
		Language: quantum.LanguageSpanish,
	})
	if err != nil {
		log.Fatalf("fetch proforma document: %v", err)
	}
	pdf, err := base64.StdEncoding.DecodeString(doc.Document)
	if err != nil {
		log.Fatalf("decode base64: %v", err)
	}
	name := doc.Filename
	if name == "" {
		name = "proforma.pdf"
	}
	if err := os.WriteFile(name, pdf, 0o644); err != nil {
		log.Fatalf("write file: %v", err)
	}
	log.Printf("saved %s (%d bytes)", name, len(pdf))
}
