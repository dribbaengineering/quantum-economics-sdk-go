package quantum

import (
	"context"
	"net/http"
)

// LabourLine is a monthly labour cost total.
type LabourLine struct {
	Month int     `json:"month" xml:"month"`
	Year  int     `json:"year" xml:"year"`
	Total float64 `json:"total" xml:"total"`
}

// LabourResponse is the envelope returned by the labour summary endpoints.
type LabourResponse struct {
	apiResponse
	LabourLine []LabourLine `json:"labourLine" xml:"labourLine"`
}

// Labour document type codes accepted by the labour document endpoint.
const (
	LabourDocPayroll = "PAYROLL"
	LabourDocRMC     = "RMC"
	LabourDocRNT     = "RNT"
	LabourDocRLC     = "RLC"
	LabourDocFolder  = "FOLDER"
	LabourDocOther   = "OTHER"
)

// LabourDocumentParams identify a labour document to fetch. All fields are
// required.
type LabourDocumentParams struct {
	// Regid is the document id.
	Regid string
	// Type is one of the LabourDoc* constants.
	Type      string
	Year      int
	CompanyID int64
}

// PayrollDocumentParams identify a payroll document to fetch. All fields are
// required.
type PayrollDocumentParams struct {
	CompanyCIF string
	Year       int
	// Regid is the payroll document id.
	Regid     string
	CompanyID int64
}

// LabourService groups the labour ("laboral") operations. The document
// endpoints return typed DocumentResponse values; the document/folder listing
// endpoints are only loosely specified upstream and return RawResponse.
type LabourService struct {
	client *Client
}

// Summary returns the labour cost detail of the last 12 months (GET /labour).
func (s *LabourService) Summary(ctx context.Context) (*LabourResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).values()
	out := &LabourResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/labour", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Detail returns the labour cost detail of a specific month
// (GET /labour/detail).
func (s *LabourService) Detail(ctx context.Context, year, month int) (*LabourResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(year)).
		setInt("month", int64(month)).
		values()
	out := &LabourResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/labour/detail", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Document returns a labour document, Base64-encoded (GET /labour/document).
func (s *LabourService) Document(ctx context.Context, params LabourDocumentParams) (*DocumentResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setString("regid", params.Regid).
		setString("type", params.Type).
		setInt("year", int64(params.Year)).
		values()
	out := &DocumentResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/labour/document", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Payroll returns the payroll list for a year (GET /labour/payroll). The
// response schema is only loosely specified upstream, so it is returned raw.
func (s *LabourService) Payroll(ctx context.Context, year int) (*RawResponse, error) {
	return s.rawByYear(ctx, "/labour/payroll", year)
}

// PayrollByCompany returns the payroll list for a month
// (GET /labour/payroll/company). Returned raw.
func (s *LabourService) PayrollByCompany(ctx context.Context, year, month int) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(year)).
		setInt("month", int64(month)).
		values()
	return s.raw(ctx, "/labour/payroll/company", q)
}

// PayrollDocument returns a payroll document, Base64-encoded
// (GET /labour/payroll/document).
func (s *LabourService) PayrollDocument(ctx context.Context, params PayrollDocumentParams) (*DocumentResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setString("companyCIF", params.CompanyCIF).
		setInt("year", int64(params.Year)).
		setString("regid", params.Regid).
		values()
	out := &DocumentResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/labour/payroll/document", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ValidatePayroll validates a list of payrolls by id
// (POST /labour/payroll/validate). ids is a comma-separated list of payroll ids;
// validation is the validation state to apply. Returned raw.
func (s *LabourService) ValidatePayroll(ctx context.Context, validation, year int, ids string) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("validation", int64(validation)).
		setInt("year", int64(year)).
		setString("id", ids).
		values()
	return s.rawPost(ctx, "/labour/payroll/validate", q)
}

// ManageBasicDocuments lists basic labour documents
// (GET /labour/manageBasicDocuments). Returned raw.
func (s *LabourService) ManageBasicDocuments(ctx context.Context, year int) (*RawResponse, error) {
	return s.rawByYear(ctx, "/labour/manageBasicDocuments", year)
}

// ManageCertDocuments lists certificate documents
// (GET /labour/manageCertDocuments). Returned raw.
func (s *LabourService) ManageCertDocuments(ctx context.Context, year int) (*RawResponse, error) {
	return s.rawByYear(ctx, "/labour/manageCertDocuments", year)
}

// ManageDocuments lists the folders and documents of a folder
// (GET /labour/manageDocuments). Returned raw.
func (s *LabourService) ManageDocuments(ctx context.Context, folderID int64, year int) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("folderId", folderID).
		setInt("year", int64(year)).
		values()
	return s.raw(ctx, "/labour/manageDocuments", q)
}

// OtherDocuments lists documents that are not payrolls
// (GET /labour/otherDocuments). Returned raw.
func (s *LabourService) OtherDocuments(ctx context.Context, year int) (*RawResponse, error) {
	return s.rawByYear(ctx, "/labour/otherDocuments", year)
}

func (s *LabourService) rawByYear(ctx context.Context, path string, year int) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).setInt("year", int64(year)).values()
	return s.raw(ctx, path, q)
}

func (s *LabourService) raw(ctx context.Context, path string, q map[string][]string) (*RawResponse, error) {
	out := &RawResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: path, query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *LabourService) rawPost(ctx context.Context, path string, q map[string][]string) (*RawResponse, error) {
	out := &RawResponse{}
	if err := s.client.do(ctx, request{method: http.MethodPost, path: path, query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}
