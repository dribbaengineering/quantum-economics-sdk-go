package quantum

import (
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
)

// ProformaInvoice is the payload for creating a pro forma invoice and the item
// shape returned by the read endpoints. Its structure mirrors Invoice but omits
// the invoice type and the series/number fields, since a pro forma is validated
// and issued later from the web application.
//
// As with invoices, pass CustomerProviderID for an existing customer or embed
// Customer to create/match one inline.
type ProformaInvoice struct {
	XMLName xml.Name `json:"-" xml:"ProformaInvoice"`

	ID              int64  `json:"id,omitempty" xml:"id,omitempty"`
	SeriesAndNumber string `json:"seriesAndNumber,omitempty" xml:"seriesAndNumber,omitempty"`
	Series          string `json:"series,omitempty" xml:"series,omitempty"`
	DateForVAT      string `json:"dateForVAT,omitempty" xml:"dateForVAT,omitempty"`

	CustomerProviderID      int64   `json:"customerProviderId,omitempty" xml:"customerProviderId,omitempty"`
	Name                    string  `json:"name,omitempty" xml:"name,omitempty"`
	TotalAmountWithoutTaxes float64 `json:"totalAmountWithoutTaxes,omitempty" xml:"totalAmountWithoutTaxes,omitempty"`
	TotalAmount             float64 `json:"totalAmount,omitempty" xml:"totalAmount,omitempty"`

	Rental                            bool    `json:"rental,omitempty" xml:"rental,omitempty"`
	PaymentMethod                     int     `json:"paymentMethod,omitempty" xml:"paymentMethod,omitempty"`
	PaymentType                       int     `json:"paymentType,omitempty" xml:"paymentType,omitempty"`
	BankID                            int     `json:"bankId,omitempty" xml:"bankId,omitempty"`
	CustomerProviderBankAccountNumber string  `json:"customerProviderBankAccountNumber,omitempty" xml:"customerProviderBankAccountNumber,omitempty"`
	ActivityID                        int64   `json:"activityId,omitempty" xml:"activityId,omitempty"`
	PropertyID                        int64   `json:"propertyId,omitempty" xml:"propertyId,omitempty"`
	VAT                               float64 `json:"vat,omitempty" xml:"vat,omitempty"`
	RecargoEquivalencia               float64 `json:"recargoEquivalencia,omitempty" xml:"recargoEquivalencia,omitempty"`
	Comments                          string  `json:"comments,omitempty" xml:"comments,omitempty"`

	Tax         []InvoiceTax         `json:"tax,omitempty" xml:"tax,omitempty"`
	Installment []InvoiceInstallment `json:"installment,omitempty" xml:"installment,omitempty"`
	Line        []InvoiceLine        `json:"line,omitempty" xml:"line,omitempty"`

	IRPFTaxCode       int     `json:"irpfTaxCode,omitempty" xml:"irpfTaxCode,omitempty"`
	IRPFTaxPercentage float64 `json:"irpfTaxPercentage,omitempty" xml:"irpfTaxPercentage,omitempty"`
	CashCriterion     bool    `json:"cashCriterion,omitempty" xml:"cashCriterion,omitempty"`

	Free1 string `json:"free1,omitempty" xml:"free1,omitempty"`
	Free2 string `json:"free2,omitempty" xml:"free2,omitempty"`
	Free3 string `json:"free3,omitempty" xml:"free3,omitempty"`
	Free4 string `json:"free4,omitempty" xml:"free4,omitempty"`
	Free5 string `json:"free5,omitempty" xml:"free5,omitempty"`
	Free6 string `json:"free6,omitempty" xml:"free6,omitempty"`

	Tobacco        bool   `json:"tobacco,omitempty" xml:"tobacco,omitempty"`
	Packs          int64  `json:"packs,omitempty" xml:"packs,omitempty"`
	InvoiceDate    string `json:"invoiceDate,omitempty" xml:"invoiceDate,omitempty"`
	DescriptionSII string `json:"descriptionSII,omitempty" xml:"descriptionSII,omitempty"`
	TravelAgency   bool   `json:"travelAgency,omitempty" xml:"travelAgency,omitempty"`
	OperationDate  string `json:"operationDate,omitempty" xml:"operationDate,omitempty"`

	CountryISO      string `json:"countryISO,omitempty" xml:"countryISO,omitempty"`
	Currency        string `json:"currency,omitempty" xml:"currency,omitempty"`
	Agrarian        bool   `json:"agrarian,omitempty" xml:"agrarian,omitempty"`
	InstallmentDate string `json:"installmentDate,omitempty" xml:"installmentDate,omitempty"`

	Customer     *Customer `json:"customer,omitempty" xml:"customer,omitempty"`
	ShopifyOrder string    `json:"shopifyOrder,omitempty" xml:"shopifyOrder,omitempty"`
}

// ProformaListResponse is the envelope returned by the pro forma listing
// endpoints, including aggregate figures and pagination metadata.
type ProformaListResponse struct {
	apiResponse
	Proformas []ProformaInvoice `json:"proformas" xml:"proformas"`
	RadChart  []ProformaInvoice `json:"radChart" xml:"radChart"`
	Income    float64           `json:"income" xml:"income"`
	Expenses  float64           `json:"expenses" xml:"expenses"`
	Balance   float64           `json:"balance" xml:"balance"`

	ActualPage       int `json:"actualPage" xml:"actualPage"`
	TotalPages       int `json:"totalPages" xml:"totalPages"`
	InvoicesQuantity int `json:"invoicesQuantity" xml:"invoicesQuantity"`
}

// ListProformaParams are the parameters for listing pro forma invoices.
type ListProformaParams struct {
	// StartDate / EndDate bound the search, formatted dd-mm-yyyy (optional).
	StartDate string
	EndDate   string
	// Page selects a page for the paged List call (optional, ignored by
	// ListFull).
	Page int
	// CompanyID overrides the client's default company id (optional).
	CompanyID int64
}

// ProformaService groups the pro forma invoice operations.
type ProformaService struct {
	client *Client
}

// Create sends a pro forma invoice (POST /proforma) and returns its id.
func (s *ProformaService) Create(ctx context.Context, p ProformaInvoice) (*IDResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).values()
	out := &IDResponse{}
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/proforma", query: q, body: p}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// List returns a paged list of pro forma invoices (GET /proforma).
func (s *ProformaService) List(ctx context.Context, params ListProformaParams) (*ProformaListResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setStringOpt("startDate", params.StartDate).
		setStringOpt("endDate", params.EndDate).
		setIntOpt("page", int64(params.Page)).
		values()
	out := &ProformaListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/proforma", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListFull returns the full (non-paged) list of pro forma invoices
// (GET /proforma/full).
func (s *ProformaService) ListFull(ctx context.Context, params ListProformaParams) (*ProformaListResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setStringOpt("startDate", params.StartDate).
		setStringOpt("endDate", params.EndDate).
		values()
	out := &ProformaListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/proforma/full", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single pro forma invoice by id (GET /proforma/{id}). It returns
// ErrNotFound when the pro forma does not exist.
func (s *ProformaService) Get(ctx context.Context, id int64) (*ProformaInvoice, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).values()
	out := &ProformaListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/proforma/" + strconv.FormatInt(id, 10), query: q}, out); err != nil {
		return nil, err
	}
	if len(out.Proformas) == 0 {
		return nil, ErrNotFound
	}
	return &out.Proformas[0], nil
}

// Dashboard returns the pro forma dashboard list (GET /proforma/dashboard).
func (s *ProformaService) Dashboard(ctx context.Context) (*ProformaListResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).values()
	out := &ProformaListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/proforma/dashboard", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PDFURL returns a temporary URL to the pro forma PDF (GET /proforma/pdf).
func (s *ProformaService) PDFURL(ctx context.Context, params DocumentParams) (string, error) {
	return s.docString(ctx, "/proforma/pdf", params)
}

// PDFBase64 returns the pro forma PDF Base64-encoded (GET /proforma/pdfB64).
func (s *ProformaService) PDFBase64(ctx context.Context, params DocumentParams) (string, error) {
	return s.docString(ctx, "/proforma/pdfB64", params)
}

// Document returns the pro forma document embedded in a JSON envelope
// (GET /proforma/document).
func (s *ProformaService) Document(ctx context.Context, params DocumentParams) (*DocumentResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setString("id", params.ID).
		setStringOpt("language", params.Language).
		values()
	out := &DocumentResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/proforma/document", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *ProformaService) docString(ctx context.Context, path string, params DocumentParams) (string, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return "", err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setString("id", params.ID).
		setStringOpt("language", params.Language).
		values()
	var out string
	if err := s.client.do(ctx, request{method: http.MethodGet, path: path, query: q}, &out); err != nil {
		return "", err
	}
	return out, nil
}
