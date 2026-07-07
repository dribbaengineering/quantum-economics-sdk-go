// Command xml_format shows how to send and receive XML instead of JSON. Quantum
// accepts both formats; the WithXML option flips the Content-Type and Accept
// headers and the request/response codecs.
//
// Run with:
//
//	QUANTUM_API_KEY=... QUANTUM_COMPANY_ID=28218 go run ./examples/xml_format
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
		quantum.WithXML(), // send and receive application/xml
	)
	if err != nil {
		log.Fatalf("client: %v", err)
	}

	invoice := quantum.Invoice{
		Type:                    quantum.InvoiceTypeIssued,
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
		Installment:    []quantum.InvoiceInstallment{{Date: "22/09/2025", Amount: 121.00}},
		DescriptionSII: "example invoice sent as XML",
	}

	res, err := client.Invoices.Create(context.Background(), invoice)
	if err != nil {
		log.Fatalf("issue invoice: %v", err)
	}
	log.Printf("issued invoice via XML, id=%d", res.ID)
}
