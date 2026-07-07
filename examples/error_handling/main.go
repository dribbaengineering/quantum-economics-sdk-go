// Command error_handling demonstrates the SDK's error model: sentinel errors,
// *APIError (business errors inside a 200 envelope) and *HTTPError (transport
// failures).
//
// Run with:
//
//	QUANTUM_API_KEY=... QUANTUM_COMPANY_ID=28218 go run ./examples/error_handling
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

	_, err = client.Customers.GetByID(context.Background(), 999999999)
	if err == nil {
		log.Println("customer found")
		return
	}

	// Categorize the failure.
	switch {
	case errors.Is(err, quantum.ErrNotFound):
		log.Println("no such customer")
	case errors.Is(err, quantum.ErrUnauthorized):
		log.Println("check your API key")
	default:
		var apiErr *quantum.APIError
		var httpErr *quantum.HTTPError
		switch {
		case errors.As(err, &apiErr):
			log.Printf("quantum business error %d: %s", apiErr.Code, apiErr.Message)
		case errors.As(err, &httpErr):
			log.Printf("transport error: HTTP %d", httpErr.StatusCode)
		default:
			log.Printf("unexpected error: %v", err)
		}
	}
}
