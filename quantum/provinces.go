package quantum

import (
	"context"
	"net/http"
	"strings"
)

// Province is a province record returned by the province lookup endpoint.
// ProvinceCode is the value CitiesService.SearchInProvince expects: the
// two-digit Spanish province code (e.g. "08" for Barcelona, "28" for Madrid).
type Province struct {
	ProvinceName string `json:"provinceName" xml:"provinceName"`
	ProvinceCode string `json:"provinceCode" xml:"provinceCode"`
}

// ProvinceListResponse is the envelope returned by the province lookup endpoint.
type ProvinceListResponse struct {
	apiResponse
	Provinces []Province `json:"provinces" xml:"provinces"`
}

// ProvincesService lists the provinces of a country. Its only purpose is to
// supply the province code that a Spanish city-name search requires.
//
// NOTE: like the /city endpoints, /provinces is NOT part of Quantum's published
// API contract — it is absent from the official Swagger/Postman specification
// and may change without notice. It is included because a Spanish city-name
// search is impossible without a province code and no documented endpoint
// offers one.
type ProvincesService struct {
	client *Client
}

// List returns every province of the given country (ISO code).
func (s *ProvincesService) List(ctx context.Context, countryISO string) (*ProvinceListResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	out := &ProvinceListResponse{}
	q := newQuery().
		setString("country", strings.ToUpper(countryISO)).
		setInt("companyId", companyID).
		values()
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/provinces", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}
