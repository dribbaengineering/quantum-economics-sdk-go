package quantum

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors returned (usually wrapped) by the client. Match them with
// errors.Is. Configuration errors are reported eagerly by NewClient and the
// request builder; the HTTP-status sentinels are attached to *HTTPError and
// *APIError so callers can branch on the failure category without inspecting
// numeric codes.
var (
	// ErrMissingAPIKey means no API key was configured.
	ErrMissingAPIKey = errors.New("quantum: missing API key")
	// ErrMissingCompanyID means no company id was available for a request that
	// requires one (neither on the client nor on the call parameters).
	ErrMissingCompanyID = errors.New("quantum: missing company id")
	// ErrInvalidBaseURL means the configured base URL could not be parsed.
	ErrInvalidBaseURL = errors.New("quantum: invalid base URL")

	// ErrUnauthorized maps to HTTP 401/403 (bad or missing API key).
	ErrUnauthorized = errors.New("quantum: unauthorized")
	// ErrNotFound maps to HTTP 404 (resource does not exist).
	ErrNotFound = errors.New("quantum: not found")
	// ErrBadRequest maps to HTTP 400 (malformed request).
	ErrBadRequest = errors.New("quantum: bad request")
	// ErrServer maps to HTTP 5xx (server-side failure).
	ErrServer = errors.New("quantum: server error")
)

// APIError represents a business-level failure reported inside the Quantum
// response envelope, i.e. an "error" object whose errorCode is non-zero. The
// HTTP status may still be 200 in this case, which is why it is inspected
// separately from transport errors.
type APIError struct {
	// Code is the Quantum error code (errorCode). Zero means success and is
	// never wrapped in an APIError.
	Code int
	// Message is the human-readable message returned by Quantum.
	Message string
	// HTTPStatus is the HTTP status code the envelope arrived with.
	HTTPStatus int
	// Method and Endpoint identify the call that failed, for diagnostics.
	Method   string
	Endpoint string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("quantum: API error %d on %s %s: %s", e.Code, e.Method, e.Endpoint, e.Message)
}

// Is lets callers match an APIError against the HTTP-status sentinels, e.g.
// errors.Is(err, quantum.ErrNotFound).
func (e *APIError) Is(target error) bool {
	return matchStatusSentinel(e.HTTPStatus, target)
}

// HTTPError represents a transport-level failure: a non-2xx status whose body
// could not be interpreted as a Quantum error envelope.
type HTTPError struct {
	// StatusCode is the HTTP status code returned by the server.
	StatusCode int
	// Status is the textual HTTP status (e.g. "404 Not Found").
	Status string
	// Body is the raw (possibly truncated) response body, useful for debugging.
	Body []byte
	// Method and Endpoint identify the call that failed.
	Method   string
	Endpoint string
}

func (e *HTTPError) Error() string {
	msg := fmt.Sprintf("quantum: HTTP %d on %s %s", e.StatusCode, e.Method, e.Endpoint)
	if len(e.Body) > 0 {
		snippet := e.Body
		const max = 512
		if len(snippet) > max {
			snippet = snippet[:max]
		}
		msg += ": " + string(snippet)
	}
	return msg
}

// Is lets callers match an HTTPError against the HTTP-status sentinels.
func (e *HTTPError) Is(target error) bool {
	return matchStatusSentinel(e.StatusCode, target)
}

// DecodeError wraps a failure to deserialize a response body. It keeps the raw
// body so callers can inspect what actually came back.
type DecodeError struct {
	Err  error
	Body []byte
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("quantum: decode response: %v", e.Err)
}

func (e *DecodeError) Unwrap() error { return e.Err }

// matchStatusSentinel reports whether an HTTP status corresponds to one of the
// exported category sentinels.
func matchStatusSentinel(status int, target error) bool {
	switch target {
	case ErrUnauthorized:
		return status == http.StatusUnauthorized || status == http.StatusForbidden
	case ErrNotFound:
		return status == http.StatusNotFound
	case ErrBadRequest:
		return status == http.StatusBadRequest
	case ErrServer:
		return status >= 500
	default:
		return false
	}
}
