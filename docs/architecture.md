# Architecture & design

The SDK is built as a thin, well-isolated core plus one service per API domain.
The goal is a library that is easy to use, easy to test, and open to extension
without modification.

## Layout

```
quantum-economics-sdk-go/
├── quantum/            # the SDK package (import this)
│   ├── client.go           # Client, NewClient, service wiring
│   ├── options.go          # functional options (config)
│   ├── auth.go             # Authenticator + APIKeyAuthenticator
│   ├── transport.go        # Doer interface + request pipeline (do)
│   ├── codec.go            # JSON/XML codecs + content types
│   ├── query.go            # query builder + date helpers
│   ├── errors.go           # APIError, HTTPError, DecodeError, sentinels
│   ├── response.go         # shared envelope (apiResponse)
│   ├── raw.go              # RawResponse for under-specified endpoints
│   ├── enums.go            # closed-set code constants
│   ├── logger.go           # Logger interface + no-op
│   ├── invoices.go         # one file per domain: types + service
│   ├── customers.go / providers.go / proforma.go / ...
│   └── *_test.go           # httptest-based tests
├── examples/           # runnable programs, one per scenario
├── docs/               # this documentation
└── spec/               # upstream API reference (Swagger, PDF, Postman)
```

Because Go groups a package by directory, all `quantum` source files live
together under `quantum/`; the per-domain split is by filename. Types and their
service are colocated in one file per domain for cohesion.

## SOLID

**Single responsibility.** Each concern is its own unit: `auth.go` authorizes,
`codec.go` (de)serializes, `transport.go` runs the request/response pipeline,
`errors.go` maps failures, `query.go` assembles parameters. Each service handles
exactly one API domain.

**Open/closed.** The client is configured entirely through functional `Option`s.
New knobs are added as new options without changing existing call sites, and the
config struct stays unexported so construction always goes through `NewClient`.

**Liskov substitution.** The client depends on the `Doer` interface, which the
standard `*http.Client` satisfies — as does any retrying, tracing or mock
implementation. They are interchangeable.

**Interface segregation.** The abstractions are deliberately tiny:

```go
type Doer interface          { Do(*http.Request) (*http.Response, error) }
type Authenticator interface { Authorize(*http.Request) error }
type Logger interface        { Logf(string, ...any) }
```

Implementers are never forced to provide methods they do not need.

**Dependency inversion.** The HTTP client, authenticator and logger are injected
into the client; the core depends on the interfaces above, not on concrete
types. This is what makes the SDK trivially testable with `httptest`.

## Request pipeline

`Client.do` is the single choke point for every call:

1. Resolve the full URL from the base URL, path and query.
2. Marshal the body (JSON or XML) if present.
3. Set headers (`Authorization`, `Content-Type`, `Accept`, `User-Agent`).
4. Execute via the injected `Doer`.
5. On a non-2xx status, produce an `*APIError` (if the body is a Quantum
   envelope) or an `*HTTPError`.
6. On success, decode into the typed response and check the envelope for a
   business error, returning an `*APIError` if present.

Services never touch HTTP directly — they build a `request` value and hand it to
`do`, which keeps behaviour uniform and the surface area small.

## Response typing policy

Responses are modelled from the upstream specification. A handful of endpoints
reuse an unrelated schema in the spec (some `Worker`, `Labour`, `Diary` and
`DUA` listing endpoints); rather than expose a type that may be inaccurate,
those return a [`RawResponse`](../quantum/raw.go): it still provides uniform
error handling via the envelope and the complete payload, which you can decode
into your own type with `RawResponse.Decode`. Every such method documents this
in its doc comment.

## Dates

Quantum uses two date formats: `dd-mm-yyyy` in query parameters and `dd/mm/yyyy`
in bodies. Date fields are plain strings so no format is imposed on you, but
helpers are provided: `FormatQueryDate` / `ParseQueryDate` and `FormatBodyDate` /
`ParseBodyDate`.
