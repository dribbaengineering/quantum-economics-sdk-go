package quantum

import (
	"context"
	"net/http"
	"strconv"
)

// Customer represents a customer ("cliente") record. The same structure is used
// both as a create/update request body and as a list item in responses, so all
// fields are optional on the wire (omitempty): the minimum needed to create a
// customer is NIF, Name and CityCode (CityCode is not required for foreign
// customers, identified by CountryISO).
type Customer struct {
	// Regid is the internal registry id (read-only, present in responses).
	Regid int64 `json:"regid,omitempty" xml:"regid,omitempty"`
	// CustomerID is the customer code to reference in invoices. This is the
	// value to pass as Invoice.CustomerProviderID.
	CustomerID string `json:"customerId,omitempty" xml:"customerId,omitempty"`

	NIF        string `json:"nif,omitempty" xml:"nif,omitempty"`
	Name       string `json:"name,omitempty" xml:"name,omitempty"`
	CountryISO string `json:"countryISO,omitempty" xml:"countryISO,omitempty"`
	Email      string `json:"email,omitempty" xml:"email,omitempty"`
	Phone      string `json:"phone,omitempty" xml:"phone,omitempty"`

	StreetType   string `json:"streetType,omitempty" xml:"streetType,omitempty"`
	StreetName   string `json:"streetName,omitempty" xml:"streetName,omitempty"`
	StreetNumber string `json:"streetNumber,omitempty" xml:"streetNumber,omitempty"`
	Staircase    string `json:"staircase,omitempty" xml:"staircase,omitempty"`
	Floor        string `json:"floor,omitempty" xml:"floor,omitempty"`
	Room         string `json:"room,omitempty" xml:"room,omitempty"`
	PostCode     string `json:"postCode,omitempty" xml:"postCode,omitempty"`
	// CityCode is the INE municipality code. Required for domestic customers.
	CityCode string `json:"cityCode,omitempty" xml:"cityCode,omitempty"`

	IBAN          string `json:"iban,omitempty" xml:"iban,omitempty"`
	Swift         string `json:"swift,omitempty" xml:"swift,omitempty"`
	PaymentMethod int    `json:"paymentMethod,omitempty" xml:"paymentMethod,omitempty"`
	Family        int    `json:"family,omitempty" xml:"family,omitempty"`

	// Fields specific to customers (not present on providers).
	MandateReference string `json:"mandateReference,omitempty" xml:"mandateReference,omitempty"`
	MandateDate      string `json:"mandateDate,omitempty" xml:"mandateDate,omitempty"`
	AdvisorOrCompany string `json:"advisorOrCompany,omitempty" xml:"advisorOrCompany,omitempty"`
	IsCashCustomer   bool   `json:"isCashCustomer,omitempty" xml:"isCashCustomer,omitempty"`
}

// CustomerListResponse is the envelope returned by the customer read endpoints.
type CustomerListResponse struct {
	apiResponse
	Customers []Customer `json:"customers" xml:"customers"`
}

// CustomersService groups the operations on customers ("clientes").
type CustomersService struct {
	client *Client
}

// List returns every customer of the company.
func (s *CustomersService) List(ctx context.Context) (*CustomerListResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &CustomerListResponse{}
	q := newQuery().setInt("companyId", companyID).values()
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/customer", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetByID fetches a single customer by its internal id. It returns ErrNotFound
// when no customer matches.
func (s *CustomersService) GetByID(ctx context.Context, id int64) (*Customer, error) {
	return s.getOne(ctx, "/customer/"+strconv.FormatInt(id, 10))
}

// GetByNIF fetches a single customer by its NIF. Handy to confirm a customer
// exists (and obtain its id) before sending an invoice. Returns ErrNotFound
// when no customer matches.
func (s *CustomersService) GetByNIF(ctx context.Context, nif string) (*Customer, error) {
	return s.getOne(ctx, "/customer/nif/"+nif)
}

// Create registers a new customer and returns the assigned id.
func (s *CustomersService) Create(ctx context.Context, c Customer) (*IDResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &IDResponse{}
	q := newQuery().setInt("companyId", companyID).values()
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/customer", query: q, body: c}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update modifies an existing customer.
func (s *CustomersService) Update(ctx context.Context, c Customer) (*IDResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &IDResponse{}
	q := newQuery().setInt("companyId", companyID).values()
	if err := s.client.do(ctx, request{method: http.MethodPut, path: "/customer", query: q, body: c}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// getOne performs a customer lookup that returns a list and unwraps the first
// element.
func (s *CustomersService) getOne(ctx context.Context, path string) (*Customer, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &CustomerListResponse{}
	q := newQuery().setInt("companyId", companyID).values()
	if err := s.client.do(ctx, request{method: http.MethodGet, path: path, query: q}, out); err != nil {
		return nil, err
	}
	if len(out.Customers) == 0 {
		return nil, ErrNotFound
	}
	return &out.Customers[0], nil
}
