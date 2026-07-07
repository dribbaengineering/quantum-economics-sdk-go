package quantum

import (
	"context"
	"net/http"
)

// Movement is a portfolio (collection/payment) movement to register.
type Movement struct {
	// Type is "C" (collection) or "B" (payment).
	Type          string  `json:"type,omitempty" xml:"type,omitempty"`
	Date          string  `json:"date,omitempty" xml:"date,omitempty"`
	NIF           string  `json:"nif,omitempty" xml:"nif,omitempty"`
	InvoiceNumber string  `json:"invoiceNumber,omitempty" xml:"invoiceNumber,omitempty"`
	Total         float64 `json:"total,omitempty" xml:"total,omitempty"`
	Account       string  `json:"account,omitempty" xml:"account,omitempty"`
}

// PortfolioService groups the portfolio ("cartera") operations.
type PortfolioService struct {
	client *Client
}

// NewMovement registers a new portfolio movement (POST /portfolio/movement).
func (s *PortfolioService) NewMovement(ctx context.Context, m Movement) (*IDResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).values()
	out := &IDResponse{}
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/portfolio/movement", query: q, body: m}, out); err != nil {
		return nil, err
	}
	return out, nil
}
