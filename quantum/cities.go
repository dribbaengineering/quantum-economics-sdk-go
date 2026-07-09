package quantum

import (
	"context"
	"net/http"
	"strings"
)

// City is a locality record returned by the city lookup endpoints. CityCode is
// the value Quantum expects in Customer.CityCode / Provider.CityCode. It is
// Quantum's own locality code: the GeoNames geonameId for Andorran localities,
// and a Quantum-internal code (not the raw INE municipality code) for Spanish
// ones.
type City struct {
	CityName     string `json:"cityName" xml:"cityName"`
	CityCode     string `json:"cityCode" xml:"cityCode"`
	ProvinceName string `json:"provinceName" xml:"provinceName"`
	ProvinceCode string `json:"provinceCode" xml:"provinceCode"`
}

// CityListResponse is the envelope returned by the city lookup endpoints.
type CityListResponse struct {
	apiResponse
	Cities []City `json:"cities" xml:"cities"`
}

// CitiesService resolves localities to the cityCode required when creating a
// domestic customer or provider.
//
// NOTE: the /city endpoints are NOT part of Quantum's published API contract —
// they are absent from the official Swagger/Postman specification and may change
// without notice. They are included because resolving a cityCode is a
// prerequisite for creating domestic customers/providers and no documented
// endpoint offers it.
type CitiesService struct {
	client *Client
}

// SearchByZipCode returns the localities matching a postal code in the given
// country (ISO code). This is the most reliable lookup: a postal code maps to a
// single locality, so the first result is the one to use.
func (s *CitiesService) SearchByZipCode(ctx context.Context, countryISO, postalCode string) (*CityListResponse, error) {
	return s.get(ctx, "/city/searchByZipCode", newQuery().
		setString("country", strings.ToUpper(countryISO)).
		setString("code", postalCode))
}

// SearchByCode returns the locality with the given city code — the reverse of
// SearchByZipCode / Search.
func (s *CitiesService) SearchByCode(ctx context.Context, countryISO, cityCode string) (*CityListResponse, error) {
	return s.get(ctx, "/city/searchByCode", newQuery().
		setString("country", strings.ToUpper(countryISO)).
		setString("code", cityCode))
}

// Search returns the localities whose name matches query, in the given country.
// The server-side match is case- and accent-sensitive.
func (s *CitiesService) Search(ctx context.Context, countryISO, query string) (*CityListResponse, error) {
	return s.get(ctx, "/city/search", newQuery().
		setString("country", strings.ToUpper(countryISO)).
		setString("query", query))
}

func (s *CitiesService) get(ctx context.Context, path string, q *queryBuilder) (*CityListResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &CityListResponse{}
	query := q.setInt("companyId", companyID).values()
	if err := s.client.do(ctx, request{method: http.MethodGet, path: path, query: query}, out); err != nil {
		return nil, err
	}
	return out, nil
}
