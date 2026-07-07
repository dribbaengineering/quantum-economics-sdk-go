package quantum

import (
	"context"
	"net/http"
)

// DUAParams are the filters for listing DUAs (single administrative documents
// for customs). Year is required; the date range is optional.
type DUAParams struct {
	Year      int
	StartDate string
	EndDate   string
	CompanyID int64
}

// DUAService groups the DUA ("Documento Único Administrativo") operations.
type DUAService struct {
	client *Client
}

// List returns the company DUAs (GET /dua). The response schema is not reliably
// specified upstream, so it is returned raw.
func (s *DUAService) List(ctx context.Context, params DUAParams) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setStringOpt("startDate", params.StartDate).
		setStringOpt("endDate", params.EndDate).
		values()
	out := &RawResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/dua", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}
