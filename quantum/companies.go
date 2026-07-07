package quantum

import (
	"context"
	"net/http"
)

// Company describes a company the API key has access to. The Code field is the
// companyId used throughout the API.
type Company struct {
	ID   int64  `json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
	// Code is the company code (companyId).
	Code int64  `json:"code" xml:"code"`
	NIF  string `json:"nif" xml:"nif"`
	// PersonType is "SELF_EMPLOYED" or "COMPANY".
	PersonType            string   `json:"personType" xml:"personType"`
	ControlTreasury       bool     `json:"controlTreasury" xml:"controlTreasury"`
	CompanyManagement     int      `json:"companyManagement" xml:"companyManagement"`
	ObligationsPeriodType string   `json:"obligationsPeriodType" xml:"obligationsPeriodType"`
	ComplaintKey          string   `json:"complaintKey" xml:"complaintKey"`
	Plan                  string   `json:"plan" xml:"plan"`
	Series                []string `json:"series" xml:"series"`
	Addons                []string `json:"addons" xml:"addons"`
	WorkTimeGPS           bool     `json:"workTimeGPS" xml:"workTimeGPS"`

	// Module version markers, as returned by the API.
	VFiscal         string `json:"vfiscal" xml:"vfiscal"`
	VLaboral        string `json:"vlaboral" xml:"vlaboral"`
	VOFiscales      string `json:"vofiscales" xml:"vofiscales"`
	VRegMercantil   string `json:"vregmercantil" xml:"vregmercantil"`
	VDocument       string `json:"vdocument" xml:"vdocument"`
	VIncidencias    string `json:"vincidencias" xml:"vincidencias"`
	VActividades    string `json:"vactividades" xml:"vactividades"`
	VFirmar         string `json:"vfirmar" xml:"vfirmar"`
	VContabilidad   string `json:"vcontabilidad" xml:"vcontabilidad"`
	VControlHorario string `json:"vcontrolhorario" xml:"vcontrolhorario"`
	VAusencias      string `json:"vausencias" xml:"vausencias"`
	VCanalDenuncias string `json:"vcanaldenuncias" xml:"vcanaldenuncias"`
}

// CompanyListResponse is the envelope returned by GET /company.
type CompanyListResponse struct {
	apiResponse
	Companies []Company `json:"companies" xml:"companies"`
}

// CompaniesService groups the company operations.
type CompaniesService struct {
	client *Client
}

// List returns the companies the API key can access for a given year
// (GET /company). It does not require a company id.
func (s *CompaniesService) List(ctx context.Context, year int) (*CompanyListResponse, error) {
	q := newQuery().setInt("year", int64(year)).values()
	out := &CompanyListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/company", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFull returns the full, legacy-shaped company record (GET /company/full).
// The payload is extensive and only loosely specified upstream, so it is
// returned raw for the caller to decode as needed.
func (s *CompaniesService) GetFull(ctx context.Context, empresa int64, year int) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("empresa", empresa).
		setInt("year", int64(year)).
		values()
	out := &RawResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/company/full", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// State returns the company state token (GET /company/state).
func (s *CompaniesService) State(ctx context.Context) (*RawResponse, error) {
	out := &RawResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/company/state"}, out); err != nil {
		return nil, err
	}
	return out, nil
}
