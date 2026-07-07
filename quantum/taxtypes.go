package quantum

import (
	"context"
	"net/http"
)

// TaxType is a tax rate configured for the company (VAT/IGIC, IRPF or recargo).
type TaxType struct {
	// Type is "G" (VAT/IVA general), "I" (IRPF) or "R" (recargo).
	Type          string  `json:"type" xml:"type"`
	Code          int     `json:"code" xml:"code"`
	Description   string  `json:"description" xml:"description"`
	Percentage    float64 `json:"percentage" xml:"percentage"`
	ReverseCharge bool    `json:"reverseCharge" xml:"reverseCharge"`
	Country       string  `json:"country" xml:"country"`
	Valid         bool    `json:"valid" xml:"valid"`
	IGIC          bool    `json:"igic" xml:"igic"`
}

// TaxTypeListResponse is the envelope returned by the tax-type endpoints.
type TaxTypeListResponse struct {
	apiResponse
	Types []TaxType `json:"types" xml:"types"`
}

// TaxTypesParams are the parameters for listing tax types (GET /taxesTypes).
type TaxTypesParams struct {
	// Type optionally filters by kind ("G"/"I"/"R").
	Type string
	// Valid optionally filters to only valid (or invalid) tax types.
	Valid     *bool
	CompanyID int64
}

// TaxTypesService groups the tax-type operations.
type TaxTypesService struct {
	client *Client
}

// List returns the defined tax types (GET /taxesTypes).
func (s *TaxTypesService) List(ctx context.Context, params TaxTypesParams) (*TaxTypeListResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setStringOpt("type", params.Type).
		setBoolOpt("valid", params.Valid).
		values()
	out := &TaxTypeListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/taxesTypes", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Full returns the defined tax types via the full endpoint
// (GET /taxesTypes/full). Valid optionally filters the result.
func (s *TaxTypesService) Full(ctx context.Context, valid *bool) (*TaxTypeListResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).setBoolOpt("valid", valid).values()
	out := &TaxTypeListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/taxesTypes/full", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}
