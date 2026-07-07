package quantum

import (
	"context"
	"net/http"
	"net/url"
)

// SectionHonorary is an aggregated honorarium figure broken down by section,
// client or user, depending on the endpoint that returned it.
type SectionHonorary struct {
	SectionID   int64   `json:"sectionId" xml:"sectionId"`
	SectionName string  `json:"sectionName" xml:"sectionName"`
	Honorary    float64 `json:"honorary" xml:"honorary"`
	TotalCost   float64 `json:"totalCost" xml:"totalCost"`
	Margin      float64 `json:"margin" xml:"margin"`
}

// SectionHonoraryResponse is the envelope returned by most QuantumBI endpoints.
type SectionHonoraryResponse struct {
	apiResponse
	Data []SectionHonorary `json:"data" xml:"data"`
}

// RateHonorary is an aggregated honorarium figure broken down by rate.
type RateHonorary struct {
	RateName string  `json:"rateName" xml:"rateName"`
	Honorary float64 `json:"honorary" xml:"honorary"`
}

// RateHonoraryResponse is the envelope returned by the rate-honorary endpoint.
type RateHonoraryResponse struct {
	apiResponse
	Data []RateHonorary `json:"data" xml:"data"`
}

// HonoraryParams are the shared parameters for the QuantumBI advisor reports.
// StartDate and EndDate are required (dd-mm-yyyy); the id filters are optional.
// These endpoints operate at advisor level and do not take a company id.
type HonoraryParams struct {
	StartDate string
	EndDate   string
	FeeID     int64
	SectionID int64
	ClientID  int64
}

// QuantumBIService groups the QuantumBI advisor reporting operations.
type QuantumBIService struct {
	client *Client
}

// query assembles the shared honorary query parameters.
func (p HonoraryParams) query() url.Values {
	return newQuery().
		setString("startDate", p.StartDate).
		setString("endDate", p.EndDate).
		setIntOpt("feeId", p.FeeID).
		setIntOpt("sectionId", p.SectionID).
		setIntOpt("clientId", p.ClientID).
		values()
}

// section runs a section-honorary style call.
func (s *QuantumBIService) section(ctx context.Context, path string, params HonoraryParams) (*SectionHonoraryResponse, error) {
	out := &SectionHonoraryResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: path, query: params.query()}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// TotalHonorarium returns the total honorarium of the agency
// (GET /advisorQuantumBi/totalHonorarium).
func (s *QuantumBIService) TotalHonorarium(ctx context.Context, params HonoraryParams) (*SectionHonoraryResponse, error) {
	return s.section(ctx, "/advisorQuantumBi/totalHonorarium", params)
}

// SectionHonorary returns the honorarium broken down by section
// (GET /advisorQuantumBi/sectionHonorary).
func (s *QuantumBIService) SectionHonorary(ctx context.Context, params HonoraryParams) (*SectionHonoraryResponse, error) {
	return s.section(ctx, "/advisorQuantumBi/sectionHonorary", params)
}

// ClientOGHonorary returns the honorarium broken down by client organizational
// group (GET /advisorQuantumBi/clientOGHonorary).
func (s *QuantumBIService) ClientOGHonorary(ctx context.Context, params HonoraryParams) (*SectionHonoraryResponse, error) {
	return s.section(ctx, "/advisorQuantumBi/clientOGHonorary", params)
}

// UserHonorariums returns the honorarium broken down by user
// (GET /advisorQuantumBi/userHonorariums).
func (s *QuantumBIService) UserHonorariums(ctx context.Context, params HonoraryParams) (*SectionHonoraryResponse, error) {
	return s.section(ctx, "/advisorQuantumBi/userHonorariums", params)
}

// UserClientHonorariums returns the honorarium broken down by user and client
// (GET /advisorQuantumBi/userClientHonorariums).
func (s *QuantumBIService) UserClientHonorariums(ctx context.Context, params HonoraryParams) (*SectionHonoraryResponse, error) {
	return s.section(ctx, "/advisorQuantumBi/userClientHonorariums", params)
}

// UserSectionHonorariums returns the honorarium broken down by user and section
// (GET /advisorQuantumBi/userSectionHonorariums).
func (s *QuantumBIService) UserSectionHonorariums(ctx context.Context, params HonoraryParams) (*SectionHonoraryResponse, error) {
	return s.section(ctx, "/advisorQuantumBi/userSectionHonorariums", params)
}

// UserTimeCost returns the time cost broken down by user
// (GET /advisorQuantumBi/userTimeCost).
func (s *QuantumBIService) UserTimeCost(ctx context.Context, params HonoraryParams) (*SectionHonoraryResponse, error) {
	return s.section(ctx, "/advisorQuantumBi/userTimeCost", params)
}

// RateHonorary returns the honorarium broken down by rate
// (GET /advisorQuantumBi/rateHonorary).
func (s *QuantumBIService) RateHonorary(ctx context.Context, params HonoraryParams) (*RateHonoraryResponse, error) {
	out := &RateHonoraryResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/advisorQuantumBi/rateHonorary", query: params.query()}, out); err != nil {
		return nil, err
	}
	return out, nil
}
