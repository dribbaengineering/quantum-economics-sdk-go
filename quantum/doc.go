// Package quantum is a Go client library for the Quantum Economics accounting
// and billing API (https://app.quantumeconomics.es/contabilidad/apidoc/).
//
// The library provides a small, dependency-free core (authentication, request
// pipeline, content negotiation and error handling) and a set of focused
// services grouped by business domain. It is designed around SOLID principles:
//
//   - Single responsibility: each service handles one API domain, while the
//     transport, codec, auth and error concerns live in their own types.
//   - Open/closed: behaviour is extended through functional options and small
//     interfaces (Doer, Authenticator, Logger) rather than by editing the core.
//   - Liskov / interface segregation: the client depends on narrow interfaces
//     that any implementation can satisfy (for example a custom *http.Client
//     wrapper or a mock in tests).
//   - Dependency inversion: the concrete HTTP client, credentials and logger
//     are injected into the client, never hard-wired.
//
// # Authentication
//
// Every request must carry an "Authorization: API-KEY <key>" header and a
// "companyId" query parameter. The API key is generated per company from the
// Quantum web application (Mi Configuración → Conectores → Q-API). Provide both
// when constructing the client:
//
//	client, err := quantum.NewClient(
//		quantum.WithAPIKey("RGowc3lHV2NFSTMxVzZ6cm1wTnc2OFFJVjR6UjBiOTM="),
//		quantum.WithCompanyID(28218),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// # Making calls
//
// Services are reached through fields on the client and every method takes a
// context.Context as its first argument:
//
//	invoices, err := client.Invoices.List(ctx, quantum.ListInvoicesParams{
//		Type:      quantum.InvoiceTypeIssued,
//		StartDate: "01-01-2025",
//		EndDate:   "31-12-2025",
//	})
//
// # Error handling
//
// The API always answers with an envelope carrying an "error" object. A
// non-zero error code is surfaced as an *APIError, while transport and HTTP
// level problems become *HTTPError. Both integrate with errors.Is / errors.As
// and a set of sentinel errors (ErrUnauthorized, ErrNotFound, ...):
//
//	_, err := client.Customers.GetByID(ctx, 999)
//	var apiErr *quantum.APIError
//	if errors.As(err, &apiErr) {
//		log.Printf("quantum error %d: %s", apiErr.Code, apiErr.Message)
//	}
package quantum
