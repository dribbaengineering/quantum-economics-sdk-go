package quantum

import (
	"context"
	"net/http"
	"strconv"
)

// Provider represents a supplier ("proveedor") record. Like Customer it is used
// both as a request body and a response item, so all fields are optional; the
// minimum to create a provider is NIF, Name and CityCode (CityCode is not
// required for foreign providers).
type Provider struct {
	// Regid is the internal registry id (read-only, present in responses).
	Regid int64 `json:"regid,omitempty" xml:"regid,omitempty"`
	// CustomerID is the provider code to reference in invoices. Despite the
	// name, this is the value to pass as Invoice.CustomerProviderID.
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
	// CityCode is Quantum's own locality code (the GeoNames geonameId for
	// Andorra; a Quantum-internal code, not the raw INE code, for Spain).
	// Required for domestic providers; resolve it with the Cities service.
	CityCode string `json:"cityCode,omitempty" xml:"cityCode,omitempty"`

	IBAN          string `json:"iban,omitempty" xml:"iban,omitempty"`
	Swift         string `json:"swift,omitempty" xml:"swift,omitempty"`
	PaymentMethod int    `json:"paymentMethod,omitempty" xml:"paymentMethod,omitempty"`
	Family        int    `json:"family,omitempty" xml:"family,omitempty"`
}

// ProviderListResponse is the envelope returned by the provider read endpoints.
type ProviderListResponse struct {
	apiResponse
	Providers []Provider `json:"providers" xml:"providers"`
}

// ProvidersService groups the operations on suppliers ("proveedores").
type ProvidersService struct {
	client *Client
}

// List returns every provider of the company.
func (s *ProvidersService) List(ctx context.Context) (*ProviderListResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &ProviderListResponse{}
	q := newQuery().setInt("companyId", companyID).values()
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/provider", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetByID fetches a single provider by its internal id. Returns ErrNotFound
// when no provider matches.
func (s *ProvidersService) GetByID(ctx context.Context, id int64) (*Provider, error) {
	return s.getOne(ctx, "/provider/"+strconv.FormatInt(id, 10))
}

// GetByNIF fetches a single provider by its NIF. Returns ErrNotFound when no
// provider matches.
func (s *ProvidersService) GetByNIF(ctx context.Context, nif string) (*Provider, error) {
	return s.getOne(ctx, "/provider/nif/"+nif)
}

// Create registers a new provider and returns the assigned id.
func (s *ProvidersService) Create(ctx context.Context, p Provider) (*IDResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &IDResponse{}
	q := newQuery().setInt("companyId", companyID).values()
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/provider", query: q, body: p}, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *ProvidersService) getOne(ctx context.Context, path string) (*Provider, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &ProviderListResponse{}
	q := newQuery().setInt("companyId", companyID).values()
	if err := s.client.do(ctx, request{method: http.MethodGet, path: path, query: q}, out); err != nil {
		return nil, err
	}
	if len(out.Providers) == 0 {
		return nil, ErrNotFound
	}
	return &out.Providers[0], nil
}
