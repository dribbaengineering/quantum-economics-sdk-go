// Command register_supplier_invoice registers a received (supplier) invoice.
//
// Because SeriesAndNumber is provided, the invoice is REGISTERED (it is assumed
// to have been issued by another system) and no billing record is generated.
// The supplier is supplied inline: if no supplier with that NIF exists, Quantum
// creates it from the provided data.
//
// Run with:
//
//	QUANTUM_API_KEY=... QUANTUM_COMPANY_ID=28218 go run ./examples/register_supplier_invoice
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

	invoice := quantum.Invoice{
		Type:                    quantum.InvoiceTypeReceived,
		SeriesAndNumber:         "2025/001125",
		InvoiceDate:             "06/11/2025",
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
		Installment:    []quantum.InvoiceInstallment{{Date: "06/11/2025", Amount: 121.00}},
		DescriptionSII: "example supplier invoice",
		// Supplier created on the fly when the NIF is unknown.
		Provider: &quantum.Provider{
			NIF:      "B00000009",
			Name:     "Example Supplier",
			CityCode: "08100",
		},
	}

	res, err := client.Invoices.Create(context.Background(), invoice)
	if err != nil {
		log.Fatalf("register invoice: %v", err)
	}
	log.Printf("registered supplier invoice, id=%d", res.ID)
}
