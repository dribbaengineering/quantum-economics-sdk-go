package quantum

import (
	"context"
	"encoding/xml"
	"net/http"
)

// DeliveryNote is the payload for creating a delivery note ("albarán").
type DeliveryNote struct {
	XMLName xml.Name `json:"-" xml:"DeliveryNote"`

	ID              int64  `json:"id,omitempty" xml:"id,omitempty"`
	SeriesAndNumber string `json:"seriesAndNumber,omitempty" xml:"seriesAndNumber,omitempty"`
	Series          string `json:"series,omitempty" xml:"series,omitempty"`
	DateForVAT      string `json:"dateForVAT,omitempty" xml:"dateForVAT,omitempty"`
	// CustomerID is the customer code the delivery note is issued to.
	CustomerID              int64   `json:"customerId,omitempty" xml:"customerId,omitempty"`
	Name                    string  `json:"name,omitempty" xml:"name,omitempty"`
	TotalAmountWithoutTaxes float64 `json:"totalAmountWithoutTaxes,omitempty" xml:"totalAmountWithoutTaxes,omitempty"`
	TotalAmount             float64 `json:"totalAmount,omitempty" xml:"totalAmount,omitempty"`

	CreditMemoSeriesAndNumber string `json:"creditMemoSeriesAndNumber,omitempty" xml:"creditMemoSeriesAndNumber,omitempty"`
	PaymentMethod             int    `json:"paymentMethod,omitempty" xml:"paymentMethod,omitempty"`
	PaymentType               int    `json:"paymentType,omitempty" xml:"paymentType,omitempty"`
	BankID                    int    `json:"bankId,omitempty" xml:"bankId,omitempty"`
	ActivityID                int64  `json:"activityId,omitempty" xml:"activityId,omitempty"`
	Comments                  string `json:"comments,omitempty" xml:"comments,omitempty"`

	Tax  []DeliveryNoteTax  `json:"tax,omitempty" xml:"tax,omitempty"`
	Line []DeliveryNoteLine `json:"line,omitempty" xml:"line,omitempty"`

	IRPFTaxCode       int     `json:"irpfTaxCode,omitempty" xml:"irpfTaxCode,omitempty"`
	IRPFTaxPercentage float64 `json:"irpfTaxPercentage,omitempty" xml:"irpfTaxPercentage,omitempty"`
	DeliveryNoteDate  string  `json:"deliveryNoteDate,omitempty" xml:"deliveryNoteDate,omitempty"`
	CountryISO        string  `json:"countryISO,omitempty" xml:"countryISO,omitempty"`
	Currency          string  `json:"currency,omitempty" xml:"currency,omitempty"`
	TipoIRPF          int     `json:"tipoIRPF,omitempty" xml:"tipoIRPF,omitempty"`
}

// DeliveryNoteLine is a single line of a delivery note.
type DeliveryNoteLine struct {
	Description   string  `json:"description,omitempty" xml:"description,omitempty"`
	Quantity      float64 `json:"quantity,omitempty" xml:"quantity,omitempty"`
	Amount        float64 `json:"amount,omitempty" xml:"amount,omitempty"`
	Base          float64 `json:"base,omitempty" xml:"base,omitempty"`
	Reference     string  `json:"reference,omitempty" xml:"reference,omitempty"`
	ReferenceType string  `json:"referenceType,omitempty" xml:"referenceType,omitempty"`
	Discount      float64 `json:"discount,omitempty" xml:"discount,omitempty"`
	Percentage    float64 `json:"percentage,omitempty" xml:"percentage,omitempty"`
	TaxCode       int     `json:"taxCode,omitempty" xml:"taxCode,omitempty"`
	Cuenta        string  `json:"cuenta,omitempty" xml:"cuenta,omitempty"`
	Clase390      int64   `json:"clase390,omitempty" xml:"clase390,omitempty"`
	Codigo        int     `json:"codigo,omitempty" xml:"codigo,omitempty"`
}

// DeliveryNoteTax is a single tax entry of a delivery note.
type DeliveryNoteTax struct {
	TaxType     string  `json:"taxType,omitempty" xml:"taxType,omitempty"`
	TaxableBase float64 `json:"taxableBase,omitempty" xml:"taxableBase,omitempty"`
	Percentage  float64 `json:"percentage,omitempty" xml:"percentage,omitempty"`
	Amount      float64 `json:"amount,omitempty" xml:"amount,omitempty"`
	TaxSubType  int64   `json:"taxSubType,omitempty" xml:"taxSubType,omitempty"`
	Codigo      int     `json:"codigo,omitempty" xml:"codigo,omitempty"`
}

// DeliveryNotesService groups the delivery note ("albarán") operations.
type DeliveryNotesService struct {
	client *Client
}

// Create registers a new delivery note (POST /deliverynote) and returns its id.
func (s *DeliveryNotesService) Create(ctx context.Context, dn DeliveryNote) (*IDResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).values()
	out := &IDResponse{}
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/deliverynote", query: q, body: dn}, out); err != nil {
		return nil, err
	}
	return out, nil
}
