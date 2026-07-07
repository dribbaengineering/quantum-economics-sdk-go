package quantum

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// testServer spins up an httptest.Server whose handler is provided by the test,
// and returns a client wired to it with a fixed API key and company id.
func testServer(t *testing.T, handler http.HandlerFunc, opts ...Option) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	base := append([]Option{
		WithBaseURL(srv.URL + "/contabilidad/ws/"),
		WithAPIKey("test-key"),
		WithCompanyID(28218),
		WithHTTPClient(srv.Client()),
	}, opts...)

	client, err := NewClient(base...)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client, srv
}

func TestNewClientValidation(t *testing.T) {
	t.Run("missing api key", func(t *testing.T) {
		_, err := NewClient(WithCompanyID(1))
		if !errors.Is(err, ErrMissingAPIKey) {
			t.Fatalf("want ErrMissingAPIKey, got %v", err)
		}
	})

	t.Run("invalid base url", func(t *testing.T) {
		_, err := NewClient(WithAPIKey("k"), WithBaseURL("://nope"))
		if !errors.Is(err, ErrInvalidBaseURL) {
			t.Fatalf("want ErrInvalidBaseURL, got %v", err)
		}
	})

	t.Run("custom authenticator replaces api key requirement", func(t *testing.T) {
		_, err := NewClient(WithAuthenticator(APIKeyAuthenticator{Key: "k"}))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestAuthorizationAndCompanyID(t *testing.T) {
	var gotAuth, gotCompany, gotAccept string
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotCompany = r.URL.Query().Get("companyId")
		gotAccept = r.Header.Get("Accept")
		_, _ = io.WriteString(w, `{"error":{"message":"OK","errorCode":0},"customers":[]}`)
	})

	if _, err := client.Customers.List(context.Background()); err != nil {
		t.Fatalf("List: %v", err)
	}
	if gotAuth != "API-KEY test-key" {
		t.Errorf("Authorization = %q", gotAuth)
	}
	if gotCompany != "28218" {
		t.Errorf("companyId = %q", gotCompany)
	}
	if gotAccept != "application/json" {
		t.Errorf("Accept = %q", gotAccept)
	}
}

func TestAPIErrorFromEnvelope(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		// HTTP 200 but a business-level error inside the envelope.
		_, _ = io.WriteString(w, `{"error":{"message":"customer not found","errorCode":42}}`)
	})

	_, err := client.Customers.List(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("want *APIError, got %v", err)
	}
	if apiErr.Code != 42 || apiErr.Message != "customer not found" {
		t.Errorf("unexpected APIError: %+v", apiErr)
	}
	if apiErr.HTTPStatus != http.StatusOK {
		t.Errorf("HTTPStatus = %d", apiErr.HTTPStatus)
	}
}

func TestHTTPErrorAndSentinels(t *testing.T) {
	cases := []struct {
		status   int
		sentinel error
	}{
		{http.StatusUnauthorized, ErrUnauthorized},
		{http.StatusForbidden, ErrUnauthorized},
		{http.StatusNotFound, ErrNotFound},
		{http.StatusBadRequest, ErrBadRequest},
		{http.StatusInternalServerError, ErrServer},
	}
	for _, tc := range cases {
		client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.status)
			_, _ = io.WriteString(w, "plain text failure")
		})
		_, err := client.Customers.List(context.Background())
		if !errors.Is(err, tc.sentinel) {
			t.Errorf("status %d: want sentinel %v, got %v", tc.status, tc.sentinel, err)
		}
		var httpErr *HTTPError
		if !errors.As(err, &httpErr) {
			t.Errorf("status %d: want *HTTPError, got %v", tc.status, err)
		}
	}
}

func TestMissingCompanyID(t *testing.T) {
	client, err := NewClient(WithAPIKey("k"))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	_, err = client.Customers.List(context.Background())
	if !errors.Is(err, ErrMissingCompanyID) {
		t.Fatalf("want ErrMissingCompanyID, got %v", err)
	}
}

func TestTextEndpointReturnsRawBody(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "  https://files.example/tmp/invoice.pdf  \n")
	})
	url, err := client.Invoices.PDFURL(context.Background(), DocumentParams{ID: "10"})
	if err != nil {
		t.Fatalf("PDFURL: %v", err)
	}
	if url != "https://files.example/tmp/invoice.pdf" {
		t.Errorf("url = %q (expected trimmed body)", url)
	}
}

func TestContextCancellation(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{}`)
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.Customers.List(ctx)
	if err == nil || !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("want context canceled error, got %v", err)
	}
}
