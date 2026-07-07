# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-07-07

Initial release: a dependency-free Go client for the Quantum Economics
accounting and billing API, covering all documented endpoints.

### Added

- **Core client** (`quantum.NewClient`) with functional options:
  `WithAPIKey`, `WithCompanyID`, `WithAuthenticator`, `WithBaseURL`,
  `WithHTTPClient`, `WithHTTPTimeout`, `WithUserAgent`, `WithLogger`,
  `WithRequestFormat`, `WithResponseFormat`, `WithXML`, `WithDefaultHeader`.
- **Pluggable abstractions**: `Doer` (HTTP transport), `Authenticator`
  (with `APIKeyAuthenticator`) and `Logger` (with `LoggerFunc`).
- **JSON and XML** content negotiation in both directions.
- **Typed error model**: `*APIError` (business errors inside a 200 envelope),
  `*HTTPError` (transport failures) and `*DecodeError`, plus sentinels
  (`ErrUnauthorized`, `ErrNotFound`, `ErrBadRequest`, `ErrServer`,
  `ErrMissingAPIKey`, `ErrMissingCompanyID`, `ErrInvalidBaseURL`) that work
  with `errors.Is` / `errors.As`.
- **18 services covering every documented endpoint**: `Invoices`, `Proforma`,
  `Customers`, `Providers`, `Companies`, `Banks`, `Accounts`, `Taxes`,
  `TaxTypes`, `Listings`, `Labour`, `Workers`, `Tickets`, `Diaries`, `DUA`,
  `Risk`, `Portfolio`, `DeliveryNotes` and `QuantumBI`.
- **Date helpers** for the API's two formats: `FormatQueryDate` /
  `ParseQueryDate` (dd-mm-yyyy) and `FormatBodyDate` / `ParseBodyDate`
  (dd/mm/yyyy).
- **`RawResponse`** for the few endpoints whose upstream schema is
  under-specified, giving uniform error handling plus a decodable payload.
- **Tests** (`httptest`-based), **8 runnable examples** under `examples/`, and
  documentation under `docs/` (authentication, error handling, content
  negotiation, architecture and full endpoint coverage).

[Unreleased]: https://github.com/dribbaengineering/quantum-economics-sdk-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/dribbaengineering/quantum-economics-sdk-go/releases/tag/v0.1.0
