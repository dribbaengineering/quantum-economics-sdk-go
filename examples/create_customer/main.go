// Command create_customer creates a customer (idempotently by NIF) and shows
// how to look one up before invoicing.
//
// Run with:
//
//	QUANTUM_API_KEY=... QUANTUM_COMPANY_ID=28218 go run ./examples/create_customer
package main

import (
	"context"
	"errors"
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
	const nif = "00000011s"

	// Reuse the customer if it already exists, otherwise create it.
	existing, err := client.Customers.GetByNIF(ctx, nif)
	switch {
	case err == nil:
		log.Printf("customer already exists: id=%s name=%s", existing.CustomerID, existing.Name)
		return
	case errors.Is(err, quantum.ErrNotFound):
		// fall through to create
	default:
		log.Fatalf("lookup customer: %v", err)
	}

	res, err := client.Customers.Create(ctx, quantum.Customer{
		NIF:      nif,
		Name:     "Example Customer",
		CityCode: "08100", // INE municipality code (required for domestic customers)
	})
	if err != nil {
		log.Fatalf("create customer: %v", err)
	}
	log.Printf("created customer, id=%d", res.ID)
}
