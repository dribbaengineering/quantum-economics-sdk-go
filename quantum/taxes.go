package quantum

import (
	"context"
	"net/http"
	"net/url"
)

// Tax is a filed or computed tax form (modelo). Depending on Model, exactly one
// of the Tax303/Tax111/Tax115/Tax202/Tax130 detail blocks is populated.
type Tax struct {
	Model        string  `json:"model" xml:"model"`
	Year         int     `json:"year" xml:"year"`
	Period       string  `json:"period" xml:"period"`
	Amount       float64 `json:"amount" xml:"amount"`
	Justificante string  `json:"justificante" xml:"justificante"`
	CSV          string  `json:"csv" xml:"csv"`
	Filed        bool    `json:"filed" xml:"filed"`
	FilingMonth  int     `json:"filingMonth" xml:"filingMonth"`
	FilingYear   int     `json:"filingYear" xml:"filingYear"`
	ID           int64   `json:"id" xml:"id"`

	Tax303 *Tax303 `json:"tax303,omitempty" xml:"tax303,omitempty"`
	Tax111 *Tax111 `json:"tax111,omitempty" xml:"tax111,omitempty"`
	Tax115 *Tax115 `json:"tax115,omitempty" xml:"tax115,omitempty"`
	Tax202 *Tax202 `json:"tax202,omitempty" xml:"tax202,omitempty"`
	Tax130 *Tax130 `json:"tax130,omitempty" xml:"tax130,omitempty"`
}

// Tax303 is the VAT self-assessment (modelo 303) detail.
type Tax303 struct {
	Difference            float64      `json:"difference" xml:"difference"`
	EarlierCompensation   float64      `json:"earlierCompensation" xml:"earlierCompensation"`
	Amount                float64      `json:"amount" xml:"amount"`
	Result                string       `json:"result" xml:"result"`
	ProRataRegularization float64      `json:"proRataRegularization" xml:"proRataRegularization"`
	PendingCustomsDutyTax float64      `json:"pendingCustomsDutyTax" xml:"pendingCustomsDutyTax"`
	StatePercentage       float64      `json:"statePercentage" xml:"statePercentage"`
	Detail                []Tax303Line `json:"detail" xml:"detail"`
}

// Tax303Line is a single line of a modelo 303 breakdown.
type Tax303Line struct {
	Block       string  `json:"block" xml:"block"`
	Type        string  `json:"type" xml:"type"`
	Base        float64 `json:"base" xml:"base"`
	Amount      float64 `json:"amount" xml:"amount"`
	BaseCC      float64 `json:"baseCC" xml:"baseCC"`
	AmountCC    float64 `json:"amountCC" xml:"amountCC"`
	Percentage  float64 `json:"percentage" xml:"percentage"`
	Class390    int     `json:"class390" xml:"class390"`
	Description string  `json:"description" xml:"description"`
	Invoices    int     `json:"invoices" xml:"invoices"`
	Origin      string  `json:"origin" xml:"origin"`
}

// Tax111 is the withholdings self-assessment (modelo 111) detail.
type Tax111 struct {
	Amount float64      `json:"amount" xml:"amount"`
	Detail []Tax111Line `json:"detail" xml:"detail"`
}

// Tax111Line is a single line of a modelo 111 breakdown.
type Tax111Line struct {
	Type            string  `json:"type" xml:"type"`
	Recipients      int     `json:"recipients" xml:"recipients"`
	ReceptionAmount float64 `json:"receptionAmount" xml:"receptionAmount"`
	WithheldAmount  float64 `json:"withheldAmount" xml:"withheldAmount"`
}

// Tax115 is the rental withholdings self-assessment (modelo 115) detail.
type Tax115 struct {
	Amount float64      `json:"amount" xml:"amount"`
	Detail []Tax115Line `json:"detail" xml:"detail"`
}

// Tax115Line is a single line of a modelo 115 breakdown.
type Tax115Line struct {
	Invoices int     `json:"invoices" xml:"invoices"`
	Base     float64 `json:"base" xml:"base"`
	Rate     float64 `json:"rate" xml:"rate"`
	Amount   float64 `json:"amount" xml:"amount"`
}

// Tax130 is the fractioned income-tax payment (modelo 130) detail.
type Tax130 struct {
	TotalBases         float64 `json:"totalBases" xml:"totalBases"`
	TotalBasesWithheld float64 `json:"totalBasesWithheld" xml:"totalBasesWithheld"`
}

// Tax202 is the fractioned corporate-tax payment (modelo 202) detail.
type Tax202 struct {
	FractionedPaymentBase float64 `json:"fractionedPaymentBase" xml:"fractionedPaymentBase"`
	Amount                float64 `json:"amount" xml:"amount"`
	IBAN                  string  `json:"iban" xml:"iban"`
}

// TaxListResponse is the envelope returned by the tax read endpoints.
type TaxListResponse struct {
	apiResponse
	Taxes []Tax `json:"taxes" xml:"taxes"`
}

// TaxListParams are the parameters for listing taxes (GET /tax). All are
// optional.
type TaxListParams struct {
	Year      int
	Period    string
	Model     string
	Presentar *bool
	CompanyID int64
}

// TaxDocParams identify a specific tax form for the detail and document
// endpoints. Model, Year and Period are required.
type TaxDocParams struct {
	Model     string
	Year      int
	Period    string
	CompanyID int64
}

// TaxesService groups the tax operations.
type TaxesService struct {
	client *Client
}

// List returns the taxes matching the given filters (GET /tax).
func (s *TaxesService) List(ctx context.Context, params TaxListParams) (*TaxListResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setIntOpt("year", int64(params.Year)).
		setStringOpt("period", params.Period).
		setStringOpt("model", params.Model).
		setBoolOpt("presentar", params.Presentar).
		values()
	out := &TaxListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/tax", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Detail returns the detail of a specific tax form (GET /tax/detail).
func (s *TaxesService) Detail(ctx context.Context, params TaxDocParams) (*TaxListResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := s.docQuery(companyID, params)
	out := &TaxListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/tax/detail", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Full returns every tax type configured for a year (GET /tax/full).
func (s *TaxesService) Full(ctx context.Context, year int) (*TaxTypeListResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).setInt("year", int64(year)).values()
	out := &TaxTypeListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/tax/full", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PDF returns the tax form PDF as a temporary URL (GET /tax/pdf).
func (s *TaxesService) PDF(ctx context.Context, params TaxDocParams) (string, error) {
	return s.docString(ctx, "/tax/pdf", params)
}

// PDFBase64 returns the tax form PDF Base64-encoded (GET /tax/pdfB64).
func (s *TaxesService) PDFBase64(ctx context.Context, params TaxDocParams) (string, error) {
	return s.docString(ctx, "/tax/pdfB64", params)
}

// PDFOld returns a temporary URL to the tax form PDF via the legacy endpoint
// (GET /tax/pdfOld).
func (s *TaxesService) PDFOld(ctx context.Context, params TaxDocParams) (*URLResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := s.docQuery(companyID, params)
	out := &URLResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/tax/pdfOld", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *TaxesService) docQuery(companyID int64, params TaxDocParams) url.Values {
	return newQuery().
		setInt("companyId", companyID).
		setString("model", params.Model).
		setInt("year", int64(params.Year)).
		setString("period", params.Period).
		values()
}

func (s *TaxesService) docString(ctx context.Context, path string, params TaxDocParams) (string, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return "", err
	}
	var out string
	if err := s.client.do(ctx, request{method: http.MethodGet, path: path, query: s.docQuery(companyID, params)}, &out); err != nil {
		return "", err
	}
	return out, nil
}
