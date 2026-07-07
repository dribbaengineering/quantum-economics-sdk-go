# Error handling

The Quantum API reports failures in two different ways, and the SDK maps each to
a distinct, inspectable error type.

## The response envelope

Every Quantum response is wrapped in an envelope containing an `error` object:

```json
{ "error": { "message": "OK", "errorCode": 0 }, "apiVersion": 1.0, "...": "..." }
```

An `errorCode` of `0` means success. Crucially, Quantum can return a **non-zero
error code with an HTTP 200 status** — the business call failed even though the
transport succeeded. The SDK inspects the envelope on every response, so these
are never silently ignored.

## Error types

### `*APIError`

Returned when the envelope carries a non-zero `errorCode` (a business-level
failure). Carries the Quantum code and message plus the HTTP status and the
call that failed.

```go
type APIError struct {
    Code       int
    Message    string
    HTTPStatus int
    Method     string
    Endpoint   string
}
```

### `*HTTPError`

Returned for a non-2xx HTTP status whose body is **not** a Quantum envelope
(gateway errors, HTML error pages, etc.). Carries the status and a snippet of
the raw body.

```go
type HTTPError struct {
    StatusCode int
    Status     string
    Body       []byte
    Method     string
    Endpoint   string
}
```

> If a non-2xx response *does* contain a Quantum envelope with an error code,
> you get an `*APIError` instead — the most specific error wins.

### `*DecodeError`

Returned when a successful response body cannot be deserialized. It wraps the
underlying error (available via `errors.Unwrap`) and keeps the raw body.

## Sentinel errors

Both `*APIError` and `*HTTPError` implement `Is`, so you can match them against
category sentinels regardless of the concrete type:

| Sentinel | Meaning |
| --- | --- |
| `ErrUnauthorized` | HTTP 401 / 403 — bad or missing API key |
| `ErrNotFound` | HTTP 404 — resource does not exist |
| `ErrBadRequest` | HTTP 400 — malformed request |
| `ErrServer` | HTTP 5xx — server-side failure |

Configuration-time sentinels are returned directly by the constructor and the
request builder:

| Sentinel | Meaning |
| --- | --- |
| `ErrMissingAPIKey` | No API key (and no custom authenticator) configured |
| `ErrMissingCompanyID` | No company id available for a call that needs one |
| `ErrInvalidBaseURL` | The configured base URL could not be parsed |

The single-resource getters (`Invoices.Get`, `Customers.GetByID`,
`Customers.GetByNIF`, `Providers.GetByID`, `Providers.GetByNIF`) return
`ErrNotFound` when the API responds successfully but with an empty list.

## Recommended pattern

```go
_, err := client.Invoices.Get(ctx, id)
switch {
case err == nil:
    // ok
case errors.Is(err, quantum.ErrNotFound):
    // handle missing invoice
case errors.Is(err, quantum.ErrUnauthorized):
    // re-check credentials
default:
    var apiErr *quantum.APIError
    var httpErr *quantum.HTTPError
    switch {
    case errors.As(err, &apiErr):
        log.Printf("quantum error %d: %s", apiErr.Code, apiErr.Message)
    case errors.As(err, &httpErr):
        log.Printf("HTTP %d: %s", httpErr.StatusCode, httpErr.Body)
    default:
        log.Printf("unexpected: %v", err)
    }
}
```

See the runnable [`examples/error_handling`](../examples/error_handling).
