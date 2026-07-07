# Quantum Economics SDK for Go

A Go client library for the [Quantum Economics](https://app.quantumeconomics.es)
accounting and billing API. It provides full coverage of the public API — invoicing,
pro forma invoices, customers and suppliers, accounting, taxes, banking, labour,
time tracking and more — behind a small, dependency-free, idiomatic Go surface.

- **Zero dependencies** — standard library only.
- **Full API coverage** — every documented endpoint, grouped into 18 services.
- **Solid error handling** — typed `*APIError` / `*HTTPError` plus sentinel errors that work with `errors.Is` / `errors.As`.
- **JSON and XML** — send and receive either format.
- **SOLID design** — pluggable HTTP client, authenticator and logger via small interfaces.

> The full API reference material this SDK is built from lives in [`spec/`](spec/)
> (Swagger definition, the official PDF and a Postman collection).

## Installation

```sh
go get github.com/dribbaengineering/quantum-economics-sdk-go/quantum
```

Requires Go 1.22+.

## Authentication

Every request needs an **API key** and a **company id** (`companyId`). Generate the
API key from the Quantum web app: *Mi Configuración → Conectores → Q-API*. The
company id is the number shown in front of the company name. See
[docs/authentication.md](docs/authentication.md).

```go
client, err := quantum.NewClient(
    quantum.WithAPIKey("RGowc3lHV2NFSTMxVzZ6cm1wTnc2OFFJVjR6UjBiOTM="),
    quantum.WithCompanyID(28218),
)
if err != nil {
    log.Fatal(err)
}
```

## Quick start

Issue a customer invoice with the minimum required data (Quantum generates the
number and, for VERIFACTU companies, reports it to the AEAT):

```go
res, err := client.Invoices.Create(ctx, quantum.Invoice{
    Type:                    quantum.InvoiceTypeIssued,
    CustomerProviderID:      2,
    TotalAmountWithoutTaxes: 100.00,
    TotalAmount:             121.00,
    Line: []quantum.InvoiceLine{{
        Description: "Order 1", Quantity: 1, Amount: 121, Base: 100, Percentage: 21, TaxCode: 21,
    }},
    Installment:    []quantum.InvoiceInstallment{{Date: "22/09/2025", Amount: 121.00}},
    DescriptionSII: "example invoice",
})
if err != nil {
    log.Fatal(err)
}
log.Printf("issued invoice id=%d", res.ID)
```

### Issued vs. registered invoices

The same `Invoices.Create` call covers both flows, controlled by
`SeriesAndNumber`:

| `SeriesAndNumber` | `Type` | Result |
| --- | --- | --- |
| empty | `C` (issued) | Quantum **issues** the invoice, generating the number (uses `Series`, or the company's first series) and creating a billing record. |
| set | any | Invoice is **registered** (assumed issued elsewhere); no billing record is generated. |

## Error handling

```go
_, err := client.Customers.GetByID(ctx, 999)
switch {
case errors.Is(err, quantum.ErrNotFound):
    // resource does not exist
case errors.Is(err, quantum.ErrUnauthorized):
    // bad API key
default:
    var apiErr *quantum.APIError
    if errors.As(err, &apiErr) {
        log.Printf("quantum error %d: %s", apiErr.Code, apiErr.Message)
    }
}
```

A business error reported inside a `200 OK` envelope becomes an `*APIError`; a
non-2xx transport failure becomes an `*HTTPError`. Both match the category
sentinels (`ErrNotFound`, `ErrUnauthorized`, `ErrBadRequest`, `ErrServer`). See
[docs/error-handling.md](docs/error-handling.md).

## Services

The client exposes one service per API domain:

| Field | Domain |
| --- | --- |
| `client.Invoices` | Issue/register invoices, list, PDF, payment state |
| `client.Proforma` | Pro forma invoices |
| `client.Customers` | Customers (clientes) |
| `client.Providers` | Suppliers (proveedores) |
| `client.Companies` | Companies accessible to the key |
| `client.Banks` | Bank accounts, cash registers and movements |
| `client.Accounts` | Ledger accounts and accounting plan |
| `client.Taxes` | Tax forms (modelos) and their PDFs |
| `client.TaxTypes` | Configured tax rates |
| `client.Listings` | Balance sheet, P&L, trial balance, ledger, registration book |
| `client.Labour` | Labour costs, payrolls and documents |
| `client.Workers` | Time tracking and absences |
| `client.Tickets` | Expense tickets |
| `client.Diaries` | Accounting journal entries |
| `client.DUA` | Customs single administrative documents |
| `client.Risk` | Customer/supplier outstanding risk |
| `client.Portfolio` | Portfolio (collection/payment) movements |
| `client.DeliveryNotes` | Delivery notes (albaranes) |
| `client.QuantumBI` | Advisor honoraria / BI reports |

A full endpoint-to-method map is in [docs/endpoints.md](docs/endpoints.md).

## Examples

Runnable examples live under [`examples/`](examples/). Each reads
`QUANTUM_API_KEY` and `QUANTUM_COMPANY_ID` from the environment:

```sh
QUANTUM_API_KEY=... QUANTUM_COMPANY_ID=28218 go run ./examples/issue_invoice
```

| Example | What it shows |
| --- | --- |
| `issue_invoice` | Issue a customer invoice with minimal data |
| `register_supplier_invoice` | Register a supplier invoice, creating the supplier inline |
| `create_customer` | Look up a customer by NIF, create if missing |
| `list_invoices` | List and paginate invoices |
| `invoice_with_facturae` | Import an invoice from a Facturae file |
| `proforma` | Create a pro forma invoice and download its PDF |
| `error_handling` | Categorize failures with sentinels and typed errors |
| `xml_format` | Send and receive XML |

## Documentation

- [Authentication](docs/authentication.md)
- [Error handling](docs/error-handling.md)
- [Content negotiation (JSON / XML)](docs/content-negotiation.md)
- [Architecture & design](docs/architecture.md)
- [Endpoint coverage](docs/endpoints.md)

## Design

The library is organized around SOLID principles: a thin core (authentication,
request pipeline, content negotiation, error mapping) plus one service per
domain. Behaviour is extended through functional options and small interfaces
(`Doer`, `Authenticator`, `Logger`) — never by editing the core. See
[docs/architecture.md](docs/architecture.md).

## License

See [LICENSE](LICENSE).
