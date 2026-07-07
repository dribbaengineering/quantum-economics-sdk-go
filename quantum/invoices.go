package quantum

import (
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
)

// Invoice is the payload for creating an invoice and the item shape returned by
// the read endpoints. It models both issued (customer, Type "C") and received
// (supplier, Type "P") invoices.
//
// Issued vs. registered:
//
//   - If SeriesAndNumber is set, the invoice is REGISTERED (Quantum assumes it
//     was issued by another billing system) and no billing record is generated.
//   - If SeriesAndNumber is empty and Type is "C", Quantum ISSUES the invoice:
//     it generates the number (using Series, or the company's first series when
//     Series is empty) and, for VERIFACTU companies, reports it to the AEAT.
//
// Customer/Provider identity: pass CustomerProviderID for an existing
// customer/supplier, or embed Customer / Provider to create (or match by NIF)
// one on the fly — in which case leave CustomerProviderID unset.
//
// All fields are optional on the wire; the minimum for an issued invoice is
// Type, CustomerProviderID, TotalAmountWithoutTaxes, TotalAmount and either a
// line or a tax breakdown.
type Invoice struct {
	XMLName xml.Name `json:"-" xml:"Invoice"`

	ID              int64  `json:"id,omitempty" xml:"id,omitempty"`
	Type            string `json:"type,omitempty" xml:"type,omitempty"`
	SeriesAndNumber string `json:"seriesAndNumber,omitempty" xml:"seriesAndNumber,omitempty"`
	Series          string `json:"series,omitempty" xml:"series,omitempty"`
	DateForVAT      string `json:"dateForVAT,omitempty" xml:"dateForVAT,omitempty"`

	// CustomerProviderID is the customer (issued) or provider (received) code.
	CustomerProviderID      int64   `json:"customerProviderId,omitempty" xml:"customerProviderId,omitempty"`
	Name                    string  `json:"name,omitempty" xml:"name,omitempty"`
	TotalAmountWithoutTaxes float64 `json:"totalAmountWithoutTaxes,omitempty" xml:"totalAmountWithoutTaxes,omitempty"`
	TotalAmount             float64 `json:"totalAmount,omitempty" xml:"totalAmount,omitempty"`

	Rental                    bool   `json:"rental,omitempty" xml:"rental,omitempty"`
	CreditMemoSeriesAndNumber string `json:"creditMemoSeriesAndNumber,omitempty" xml:"creditMemoSeriesAndNumber,omitempty"`

	PaymentMethod                     int     `json:"paymentMethod,omitempty" xml:"paymentMethod,omitempty"`
	PaymentType                       int     `json:"paymentType,omitempty" xml:"paymentType,omitempty"`
	BankID                            int     `json:"bankId,omitempty" xml:"bankId,omitempty"`
	IBAN                              string  `json:"iban,omitempty" xml:"iban,omitempty"`
	CustomerProviderBankAccountNumber string  `json:"customerProviderBankAccountNumber,omitempty" xml:"customerProviderBankAccountNumber,omitempty"`
	ActivityID                        int64   `json:"activityId,omitempty" xml:"activityId,omitempty"`
	Customs                           bool    `json:"customs,omitempty" xml:"customs,omitempty"`
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

	Tobacco        bool    `json:"tobacco,omitempty" xml:"tobacco,omitempty"`
	Packs          int64   `json:"packs,omitempty" xml:"packs,omitempty"`
	InvoiceDate    string  `json:"invoiceDate,omitempty" xml:"invoiceDate,omitempty"`
	DescriptionSII string  `json:"descriptionSII,omitempty" xml:"descriptionSII,omitempty"`
	TravelAgency   bool    `json:"travelAgency,omitempty" xml:"travelAgency,omitempty"`
	Simplified     bool    `json:"simplified,omitempty" xml:"simplified,omitempty"`
	OperationDate  string  `json:"operationDate,omitempty" xml:"operationDate,omitempty"`
	Paid           float64 `json:"paid,omitempty" xml:"paid,omitempty"`
	Pending        float64 `json:"pending,omitempty" xml:"pending,omitempty"`

	EndOfSeriesNumber string `json:"endOfSeriesNumber,omitempty" xml:"endOfSeriesNumber,omitempty"`
	CountryISO        string `json:"countryISO,omitempty" xml:"countryISO,omitempty"`
	Currency          string `json:"currency,omitempty" xml:"currency,omitempty"`
	Agrarian          bool   `json:"agrarian,omitempty" xml:"agrarian,omitempty"`
	InstallmentDate   string `json:"installmentDate,omitempty" xml:"installmentDate,omitempty"`
	SWCFC             int    `json:"swcfc,omitempty" xml:"swcfc,omitempty"`

	// Customer / Provider let you create (or match by NIF) the counterparty
	// inline instead of referencing CustomerProviderID.
	Customer *Customer `json:"customer,omitempty" xml:"customer,omitempty"`
	Provider *Provider `json:"provider,omitempty" xml:"provider,omitempty"`
}

// InvoiceTax is a single tax entry of an invoice's tax breakdown.
type InvoiceTax struct {
	// TaxType is the tax-type code (e.g. "G" for general VAT).
	TaxType     string  `json:"taxType,omitempty" xml:"taxType,omitempty"`
	TaxableBase float64 `json:"taxableBase,omitempty" xml:"taxableBase,omitempty"`
	Percentage  float64 `json:"percentage,omitempty" xml:"percentage,omitempty"`
	Amount      float64 `json:"amount,omitempty" xml:"amount,omitempty"`
	// TaxSubType applies only to IRPF.
	TaxSubType int64 `json:"taxSubType,omitempty" xml:"taxSubType,omitempty"`
	// Taxable ("S"/"N") marks the base as subject to tax.
	Taxable string `json:"taxable,omitempty" xml:"taxable,omitempty"`
	// Exempt ("S"/"N") marks the base as exempt.
	Exempt string `json:"exempt,omitempty" xml:"exempt,omitempty"`
	// GoodServices ("B"/"S") is the operation kind; issued invoices only.
	GoodServices string `json:"goodServices,omitempty" xml:"goodServices,omitempty"`
	// TaxInversion ("S"/"N") marks reverse charge; received invoices only.
	TaxInversion string `json:"taxInversion,omitempty" xml:"taxInversion,omitempty"`
	// Reason is the exemption reason code (E1..E6, M1, M2).
	Reason string `json:"reason,omitempty" xml:"reason,omitempty"`
	// Calification is the operation qualification code (S1, S2, S3).
	Calification string `json:"calification,omitempty" xml:"calification,omitempty"`
	Account      string `json:"account,omitempty" xml:"account,omitempty"`
	Clase390     int64  `json:"clase390,omitempty" xml:"clase390,omitempty"`
	IncludeIS    bool   `json:"includeIS,omitempty" xml:"includeIS,omitempty"`
	Codigo       int    `json:"codigo,omitempty" xml:"codigo,omitempty"`
}

// InvoiceInstallment is a payment/collection due date of an invoice.
type InvoiceInstallment struct {
	// Date is the due date, formatted dd/mm/yyyy.
	Date   string  `json:"date,omitempty" xml:"date,omitempty"`
	Amount float64 `json:"amount,omitempty" xml:"amount,omitempty"`
	// CollectedPaid is the amount already collected/paid for this installment.
	CollectedPaid float64 `json:"collectedPaid,omitempty" xml:"collectedPaid,omitempty"`
}

// InvoiceLine is a single line of an invoice.
type InvoiceLine struct {
	Description string  `json:"description,omitempty" xml:"description,omitempty"`
	Quantity    float64 `json:"quantity,omitempty" xml:"quantity,omitempty"`
	// Amount is the total line amount, taxes included.
	Amount float64 `json:"amount,omitempty" xml:"amount,omitempty"`
	// Base is the taxable base per item.
	Base float64 `json:"base,omitempty" xml:"base,omitempty"`

	// Reference is the rate/tariff code; ReferenceType classifies it ("S"/"M").
	Reference     string  `json:"reference,omitempty" xml:"reference,omitempty"`
	ReferenceType string  `json:"referenceType,omitempty" xml:"referenceType,omitempty"`
	Discount      float64 `json:"discount,omitempty" xml:"discount,omitempty"`
	// Percentage is the VAT percentage.
	Percentage float64 `json:"percentage,omitempty" xml:"percentage,omitempty"`
	Account    string  `json:"account,omitempty" xml:"account,omitempty"`
	Class390   int     `json:"class390,omitempty" xml:"class390,omitempty"`
	TaxCode    int     `json:"taxCode,omitempty" xml:"taxCode,omitempty"`

	Free1 string `json:"free1,omitempty" xml:"free1,omitempty"`
	Free2 string `json:"free2,omitempty" xml:"free2,omitempty"`
	Free3 string `json:"free3,omitempty" xml:"free3,omitempty"`
	Free4 string `json:"free4,omitempty" xml:"free4,omitempty"`
	Free5 string `json:"free5,omitempty" xml:"free5,omitempty"`
	Free6 string `json:"free6,omitempty" xml:"free6,omitempty"`
}

// InvoiceListResponse is the envelope returned by the invoice listing endpoints.
// It also carries aggregate figures and, for the paged List call, pagination
// metadata.
type InvoiceListResponse struct {
	apiResponse
	Invoices []Invoice `json:"invoices" xml:"invoices"`
	RadChart []Invoice `json:"radChart" xml:"radChart"`
	Income   float64   `json:"income" xml:"income"`
	Expenses float64   `json:"expenses" xml:"expenses"`
	Balance  float64   `json:"balance" xml:"balance"`

	ActualPage       int `json:"actualPage" xml:"actualPage"`
	TotalPages       int `json:"totalPages" xml:"totalPages"`
	InvoicesQuantity int `json:"invoicesQuantity" xml:"invoicesQuantity"`
}

// InvoiceWithDocument bundles an invoice with an attached document (PDF or
// similar) for the invoiceWithFile endpoint.
type InvoiceWithDocument struct {
	// The xml tag matches Invoice.XMLName ("Invoice") to avoid a name clash; the
	// invoiceWithFile endpoint is documented for JSON, where the "invoice" key
	// is used.
	Invoice  Invoice  `json:"invoice" xml:"Invoice"`
	Document Document `json:"document" xml:"document"`
}

// Document is an attachment sent alongside an invoice.
type Document struct {
	// Document is the file contents, Base64-encoded.
	Document string `json:"document" xml:"document"`
	// Filename is the filename including extension (e.g. "invoice.pdf").
	Filename string `json:"filename,omitempty" xml:"filename,omitempty"`
	Comments string `json:"comments,omitempty" xml:"comments,omitempty"`
	// DocumentType optionally classifies the attachment. See the DocumentType*
	// constants.
	DocumentType string `json:"document_type,omitempty" xml:"document_type,omitempty"`
}

// Document type classifications accepted by the document_type field.
const (
	DocumentTypeIssuedInvoice             = "ISSUED_INVOICE"
	DocumentTypeReceivedInvoice           = "RECEIVED_INVOICE"
	DocumentTypeIssuedSimplifiedInvoice   = "ISSUED_SIMPLIFIED_INVOICE"
	DocumentTypeReceivedSimplifiedInvoice = "RECEIVED_SIMPLIFIED_INVOICE"
	DocumentTypeTax                       = "TAX"
	DocumentTypeExpenses                  = "EXPENSES"
	DocumentTypeDUA                       = "DUA"
	DocumentTypeMonthlyAccountingSummary  = "MONTHLY_ACCOUNTING_SUMMARY"
	DocumentTypeUsedGoodsSummary          = "USED_GOODS_SUMMARY"
	DocumentTypePayStub                   = "PAY_STUB"
	DocumentTypePayroll                   = "PAYROLL"
	DocumentTypeNotDefined                = "NOT_DEFINED"
)

// InvoiceWithFacturae wraps a Facturae document for the invoiceWithFacturae
// endpoint.
type InvoiceWithFacturae struct {
	Facturae Facturae `json:"facturae" xml:"facturae"`
}

// Facturae carries an electronic invoice (Facturae format) to be imported.
// Quantum derives the customer/provider from the document, creating it if
// necessary. The document must include an InvoiceSeriesCode.
type Facturae struct {
	// FiscalYear is the exercise the invoice belongs to.
	FiscalYear int `json:"fiscalYear,omitempty" xml:"fiscalYear,omitempty"`
	// Base64 is the Facturae XML, Base64-encoded.
	Base64 string `json:"base64,omitempty" xml:"base64,omitempty"`
	// DescriptionSII describes the object of the invoice.
	DescriptionSII string `json:"descriptionSII,omitempty" xml:"descriptionSII,omitempty"`
	// Extra carries optional accounting hints applied per line.
	Extra *FacturaeExtra `json:"extra,omitempty" xml:"extra,omitempty"`
}

// FacturaeExtra carries optional accounting information for a Facturae import.
type FacturaeExtra struct {
	// DocumentType ("F"/"P") distinguishes issued from received when ambiguous.
	DocumentType string `json:"documentType,omitempty" xml:"documentType,omitempty"`
	ProjectID    int64  `json:"projectID,omitempty" xml:"projectID,omitempty"`
	ProjectName  string `json:"projectName,omitempty" xml:"projectName,omitempty"`
	DateIVA      string `json:"dateIVA,omitempty" xml:"dateIVA,omitempty"`
	// Accounts provides per-line accounting overrides.
	Accounts      *FacturaeAccounts `json:"accounts,omitempty" xml:"accounts,omitempty"`
	Form349Key    string            `json:"form349Key,omitempty" xml:"form349Key,omitempty"`
	ActivityCode  string            `json:"activityCode,omitempty" xml:"activityCode,omitempty"`
	ActivityOrder string            `json:"activityOrder,omitempty" xml:"activityOrder,omitempty"`
}

// FacturaeAccounts groups the per-line accounting overrides of a Facturae
// import.
type FacturaeAccounts struct {
	InvoiceLine []FacturaeLine `json:"invoiceLine,omitempty" xml:"invoiceLine,omitempty"`
}

// FacturaeLine overrides the accounting details of a single Facturae line.
type FacturaeLine struct {
	LineNumber int     `json:"lineNumber,omitempty" xml:"lineNumber,omitempty"`
	Account    string  `json:"account,omitempty" xml:"account,omitempty"`
	ClassIVA   int     `json:"classIVA,omitempty" xml:"classIVA,omitempty"`
	ISP        float64 `json:"isp,omitempty" xml:"isp,omitempty"`
}

// InvoiceCollectionPayment is a single collection/payment event of an invoice.
type InvoiceCollectionPayment struct {
	Date   string  `json:"date" xml:"date"`
	Amount float64 `json:"amount" xml:"amount"`
}

// InvoicePaymentStateResponse describes the collection/payment state of an
// invoice.
type InvoicePaymentStateResponse struct {
	apiResponse
	ID                  int64                      `json:"id" xml:"id"`
	CollectionsPayments []InvoiceCollectionPayment `json:"collectionsPayments" xml:"collectionsPayments"`
	State               string                     `json:"state" xml:"state"`
	Total               float64                    `json:"total" xml:"total"`
	CollectedPaid       float64                    `json:"collectedPaid" xml:"collectedPaid"`
	Pending             float64                    `json:"pending" xml:"pending"`
}

// ListInvoicesParams are the parameters for listing invoices.
type ListInvoicesParams struct {
	// Type is required: InvoiceTypeIssued ("C") or InvoiceTypeReceived ("P").
	Type InvoiceType
	// StartDate / EndDate bound the search, formatted dd-mm-yyyy (optional).
	StartDate string
	EndDate   string
	// Page selects a page for the paged List call (optional, ignored by
	// ListFull).
	Page int
	// CompanyID overrides the client's default company id (optional).
	CompanyID int64
}

// DocumentParams identify a document to fetch (PDF URL, Base64 or inline).
type DocumentParams struct {
	// ID is the invoice/proforma id (required).
	ID string
	// Language optionally selects the document language ("es"/"ca").
	Language Language
	// CompanyID overrides the client's default company id (optional).
	CompanyID int64
}

// InvoicesService groups the invoice operations.
type InvoicesService struct {
	client *Client
}

// Create sends an invoice (POST /invoice). Depending on SeriesAndNumber the
// invoice is issued or registered — see the Invoice documentation. It returns
// the id of the created invoice.
func (s *InvoicesService) Create(ctx context.Context, inv Invoice) (*IDResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &IDResponse{}
	q := newQuery().setInt("companyId", companyID).values()
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/invoice", query: q, body: inv}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateWithDocument sends an invoice together with an attached document
// (POST /invoice/invoiceWithFile).
func (s *InvoicesService) CreateWithDocument(ctx context.Context, payload InvoiceWithDocument) (*IDResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &IDResponse{}
	q := newQuery().setInt("companyId", companyID).values()
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/invoice/invoiceWithFile", query: q, body: payload}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateWithFacturae imports an invoice from a Facturae document
// (POST /invoice/invoiceWithFacturae).
func (s *InvoicesService) CreateWithFacturae(ctx context.Context, payload InvoiceWithFacturae) (*IDResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &IDResponse{}
	q := newQuery().setInt("companyId", companyID).values()
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/invoice/invoiceWithFacturae", query: q, body: payload}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// List returns a paged list of invoices (GET /invoice).
func (s *InvoicesService) List(ctx context.Context, params ListInvoicesParams) (*InvoiceListResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setString("type", params.Type).
		setStringOpt("startDate", params.StartDate).
		setStringOpt("endDate", params.EndDate).
		setIntOpt("page", int64(params.Page)).
		values()
	out := &InvoiceListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/invoice", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListFull returns the full (non-paged) list of invoices (GET /invoice/full).
func (s *InvoicesService) ListFull(ctx context.Context, params ListInvoicesParams) (*InvoiceListResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setString("type", params.Type).
		setStringOpt("startDate", params.StartDate).
		setStringOpt("endDate", params.EndDate).
		values()
	out := &InvoiceListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/invoice/full", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single invoice by id (GET /invoice/{id}). It returns
// ErrNotFound when the invoice does not exist.
func (s *InvoicesService) Get(ctx context.Context, id int64) (*Invoice, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).values()
	out := &InvoiceListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/invoice/" + strconv.FormatInt(id, 10), query: q}, out); err != nil {
		return nil, err
	}
	if len(out.Invoices) == 0 {
		return nil, ErrNotFound
	}
	return &out.Invoices[0], nil
}

// Dashboard returns the invoice dashboard list (GET /invoice/dashboard).
func (s *InvoicesService) Dashboard(ctx context.Context) (*InvoiceListResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).values()
	out := &InvoiceListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/invoice/dashboard", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentState returns the collection/payment state of an invoice
// (GET /invoice/state/{id}). companyNIF is optional.
func (s *InvoicesService) PaymentState(ctx context.Context, id int64, companyNIF string) (*InvoicePaymentStateResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).setStringOpt("companyNIF", companyNIF).values()
	out := &InvoicePaymentStateResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/invoice/state/" + strconv.FormatInt(id, 10), query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// PDFURL returns a temporary URL to the invoice PDF (GET /invoice/pdf).
func (s *InvoicesService) PDFURL(ctx context.Context, params DocumentParams) (string, error) {
	return s.docString(ctx, "/invoice/pdf", params)
}

// PDFBase64 returns the invoice PDF Base64-encoded (GET /invoice/pdfB64).
func (s *InvoicesService) PDFBase64(ctx context.Context, params DocumentParams) (string, error) {
	return s.docString(ctx, "/invoice/pdfB64", params)
}

// Document returns the invoice document embedded in a JSON envelope
// (GET /invoice/document).
func (s *InvoicesService) Document(ctx context.Context, params DocumentParams) (*DocumentResponse, error) {
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
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/invoice/document", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// docString is the shared implementation for the text-returning document
// endpoints (PDF URL, Base64).
func (s *InvoicesService) docString(ctx context.Context, path string, params DocumentParams) (string, error) {
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
