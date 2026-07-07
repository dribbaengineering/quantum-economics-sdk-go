package quantum

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Doer is the HTTP transport abstraction the client depends on. The standard
// *http.Client satisfies it, and so does any wrapper adding retries, tracing or
// canned responses in tests. Depending on this narrow interface instead of the
// concrete *http.Client is the dependency-inversion pillar of the design.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// maxErrorBody caps how much of an unexpected response body is retained for
// error messages, so a huge HTML error page cannot blow up memory or logs.
const maxErrorBody = 4 << 10 // 4 KiB

// request describes a single API call in transport-neutral terms. Service
// methods populate it and hand it to Client.do, keeping the wiring (URL,
// headers, encoding, error mapping) in one place.
type request struct {
	method string
	// path is joined onto the client base URL; it must start with "/".
	path string
	// query holds the URL query parameters, already including companyId.
	query url.Values
	// body, when non-nil, is marshalled using the client's request format.
	body any
}

// do executes an API call and, when out is non-nil, decodes the response into
// it. out is expected to embed apiResponse so business-level errors can be
// detected uniformly. Passing a *[]byte or *string for out yields the raw body
// (base64/URL/text endpoints use this).
func (c *Client) do(ctx context.Context, r request, out any) error {
	endpoint, err := c.buildURL(r.path, r.query)
	if err != nil {
		return err
	}

	var bodyReader io.Reader
	reqCodec := codecFor(c.requestFormat)
	if r.body != nil {
		raw, err := reqCodec.Marshal(r.body)
		if err != nil {
			return fmt.Errorf("quantum: encode request body: %w", err)
		}
		bodyReader = bytes.NewReader(raw)
	}

	httpReq, err := http.NewRequestWithContext(ctx, r.method, endpoint, bodyReader)
	if err != nil {
		return fmt.Errorf("quantum: build request: %w", err)
	}
	c.applyHeaders(httpReq, r.body != nil)
	if err := c.auth.Authorize(httpReq); err != nil {
		return err
	}

	c.logger.Logf("quantum: %s %s", r.method, endpoint)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("quantum: %s %s: %w", r.method, r.path, err)
	}
	defer drainAndClose(resp.Body)

	return c.handleResponse(resp, r, out)
}

// handleResponse maps an HTTP response onto either a decoded value or an error.
// The ordering matters: an HTTP error status short-circuits to an *HTTPError
// (or *APIError when the body is a Quantum envelope), otherwise the body is
// decoded and finally inspected for an envelope-level error.
func (c *Client) handleResponse(resp *http.Response, r request, out any) error {
	respCodec := codecFor(c.responseFormat)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.httpError(resp, r, respCodec)
	}

	// Raw body passthrough for text endpoints (PDF URLs, base64 strings...).
	switch dst := out.(type) {
	case nil:
		return nil
	case *[]byte:
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("quantum: read response body: %w", err)
		}
		*dst = data
		return nil
	case *string:
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("quantum: read response body: %w", err)
		}
		*dst = string(bytes.TrimSpace(data))
		return nil
	}

	if err := decodeInto(respCodec, resp.Body, out); err != nil {
		return err
	}

	if carrier, ok := out.(errorCarrier); ok {
		if apiErr := carrier.apiError(); apiErr != nil {
			return &APIError{
				Code:       apiErr.ErrorCode,
				Message:    apiErr.Message,
				HTTPStatus: resp.StatusCode,
				Method:     r.method,
				Endpoint:   r.path,
			}
		}
	}
	return nil
}

// httpError builds the richest error it can for a non-2xx response: an
// *APIError when the body decodes into a Quantum envelope, otherwise a plain
// *HTTPError carrying a snippet of the raw body.
func (c *Client) httpError(resp *http.Response, r request, respCodec codec) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))

	var env apiResponse
	if len(bytes.TrimSpace(body)) > 0 {
		if err := respCodec.Unmarshal(body, &env); err == nil && env.Err != nil && env.Err.ErrorCode != 0 {
			return &APIError{
				Code:       env.Err.ErrorCode,
				Message:    env.Err.Message,
				HTTPStatus: resp.StatusCode,
				Method:     r.method,
				Endpoint:   r.path,
			}
		}
	}
	return &HTTPError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       body,
		Method:     r.method,
		Endpoint:   r.path,
	}
}

// buildURL resolves a service path against the base URL and attaches the query.
func (c *Client) buildURL(path string, query url.Values) (string, error) {
	rel := &url.URL{Path: strings.TrimPrefix(path, "/")}
	u := c.baseURL.ResolveReference(rel)
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}
	return u.String(), nil
}

// applyHeaders sets the standard headers on a request.
func (c *Client) applyHeaders(req *http.Request, hasBody bool) {
	for k, v := range c.defaultHeaders {
		req.Header.Set(k, v)
	}
	if hasBody {
		req.Header.Set("Content-Type", c.requestFormat.mediaType())
	}
	req.Header.Set("Accept", c.responseFormat.mediaType())
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
}

// drainAndClose fully drains and closes a response body so the underlying
// connection can be reused by the HTTP transport.
func drainAndClose(body io.ReadCloser) {
	if body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, body)
	_ = body.Close()
}
