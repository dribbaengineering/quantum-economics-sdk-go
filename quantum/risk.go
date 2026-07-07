package quantum

import (
	"context"
	"net/http"
)

// RiskResponse is the envelope returned by GET /risk. Risk is the outstanding
// risk amount of a customer/provider.
type RiskResponse struct {
	apiResponse
	Risk float64 `json:"risk" xml:"risk"`
}

// RiskParams identify the counterparty whose risk is requested. All fields are
// required.
type RiskParams struct {
	// ClientProviderNIF is the NIF of the customer/provider.
	ClientProviderNIF string
	// Type identifies whether the NIF is a customer or a provider.
	Type      string
	CompanyID int64
}

// RiskService groups the risk operations.
type RiskService struct {
	client *Client
}

// Get returns the outstanding risk of a customer/provider (GET /risk).
func (s *RiskService) Get(ctx context.Context, params RiskParams) (*RiskResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setString("clientProviderNif", params.ClientProviderNIF).
		setString("type", params.Type).
		values()
	out := &RiskResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/risk", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}
