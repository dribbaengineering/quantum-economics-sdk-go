// Command invoice_with_facturae imports an invoice from a Facturae XML file.
//
// The Facturae document is read from disk, Base64-encoded and sent to the
// invoiceWithFacturae endpoint. Quantum derives the customer/provider from the
// document, creating it if necessary.
//
// Run with:
//
//	QUANTUM_API_KEY=... QUANTUM_COMPANY_ID=28218 go run ./examples/invoice_with_facturae invoice.xsig
package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"github.com/dribbaengineering/quantum-economics-sdk-go/examples/internal/exampleutil"
	"github.com/dribbaengineering/quantum-economics-sdk-go/quantum"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: invoice_with_facturae <facturae-file>")
	}

	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("read facturae: %v", err)
	}

	client, err := quantum.NewClient(
		quantum.WithAPIKey(exampleutil.MustEnv("QUANTUM_API_KEY")),
		quantum.WithCompanyID(exampleutil.MustEnvInt("QUANTUM_COMPANY_ID")),
	)
	if err != nil {
		log.Fatalf("client: %v", err)
	}

	payload := quantum.InvoiceWithFacturae{
		Facturae: quantum.Facturae{
			FiscalYear:     2025,
			Base64:         base64.StdEncoding.EncodeToString(raw),
			DescriptionSII: "invoice imported from facturae",
			Extra: &quantum.FacturaeExtra{
				ProjectName: "Example project",
				Accounts: &quantum.FacturaeAccounts{
					InvoiceLine: []quantum.FacturaeLine{
						{LineNumber: 1, Account: "70500000", ClassIVA: 80},
					},
				},
			},
		},
	}

	res, err := client.Invoices.CreateWithFacturae(context.Background(), payload)
	if err != nil {
		log.Fatalf("import facturae: %v", err)
	}
	log.Printf("imported invoice from facturae, id=%d", res.ID)
}
