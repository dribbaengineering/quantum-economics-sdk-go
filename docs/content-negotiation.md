# Content negotiation (JSON / XML)

Quantum accepts and returns both JSON and XML. The SDK defaults to JSON (the
recommended format) but can switch either direction independently.

## Options

| Option | Effect |
| --- | --- |
| `WithRequestFormat(ct)` | Encode request bodies as `ct` and set `Content-Type`. |
| `WithResponseFormat(ct)` | Set `Accept` to `ct`, asking Quantum to answer in that format. |
| `WithXML()` | Shortcut for XML in both directions. |

`ct` is a `ContentType`: `ContentTypeJSON` or `ContentTypeXML`.

```go
// Send XML, receive XML.
client, _ := quantum.NewClient(
    quantum.WithAPIKey(key),
    quantum.WithCompanyID(id),
    quantum.WithXML(),
)

// Or mix: send JSON, ask for XML back.
client, _ := quantum.NewClient(
    quantum.WithAPIKey(key),
    quantum.WithCompanyID(id),
    quantum.WithResponseFormat(quantum.ContentTypeXML),
)
```

## How it maps to the wire

- The request body of `POST`/`PUT` calls is marshalled with the request codec.
  For XML, an `<?xml ...?>` header is prepended and the documented root element
  names are used (`<Invoice>`, `<ProformaInvoice>`, `<DeliveryNote>`, `<Customer>`, ...).
- The `Accept` header is always set to the response format; if omitted by the
  server, Quantum defaults to JSON.

Every model carries both `json` and `xml` struct tags, so the same types work in
either mode.

## Recommendation

Prefer JSON unless you have a specific reason to use XML (for example, an
existing XML pipeline). JSON responses are the API's default and the most
thoroughly exercised path.

For importing electronic invoices in the **Facturae** format, you do **not** need
XML content negotiation: send the Facturae document Base64-encoded through
`Invoices.CreateWithFacturae` while keeping the client in JSON mode. See
[`examples/invoice_with_facturae`](../examples/invoice_with_facturae).
