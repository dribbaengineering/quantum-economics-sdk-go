package quantum

import (
	"context"
	"net/http"
	"strconv"
)

// Bank represents a bank account or cash register.
type Bank struct {
	Code         int64   `json:"code" xml:"code"`
	Name         string  `json:"name" xml:"name"`
	BankCode     string  `json:"bankCode" xml:"bankCode"`
	Office       string  `json:"office" xml:"office"`
	ControlDigit string  `json:"controlDigit" xml:"controlDigit"`
	Account      string  `json:"account" xml:"account"`
	Currency     string  `json:"currency" xml:"currency"`
	Balance      float64 `json:"balance" xml:"balance"`
}

// BankListResponse is the envelope returned by the bank read endpoints.
type BankListResponse struct {
	apiResponse
	Banks []Bank `json:"banks" xml:"banks"`
}

// BankMovement is a single movement of a bank account or cash register.
type BankMovement struct {
	OperationDate string  `json:"operationDate" xml:"operationDate"`
	Text          string  `json:"text" xml:"text"`
	Amount        float64 `json:"amount" xml:"amount"`
}

// BankMovementListResponse is the paged envelope returned by the movements
// endpoint.
type BankMovementListResponse struct {
	apiResponse
	Movements      []BankMovement `json:"movements" xml:"movements"`
	ActualPage     int            `json:"actualPage" xml:"actualPage"`
	TotalPages     int            `json:"totalPages" xml:"totalPages"`
	TotalMovements int            `json:"totalMovements" xml:"totalMovements"`
}

// BankMovementsParams are the parameters for listing bank movements.
type BankMovementsParams struct {
	// BankType is required: BankTypeBank ("B") or BankTypeCash ("C").
	BankType BankType
	// StartDate / EndDate bound the search (dd-mm-yyyy), both required.
	StartDate string
	EndDate   string
	// SortBy is the required sort field.
	SortBy string
	// Page selects a page (optional).
	Page int
	// StartAmount / EndAmount optionally bound the movement amount.
	StartAmount *float64
	EndAmount   *float64
	// CompanyID overrides the client's default company id (optional).
	CompanyID int64
}

// BanksService groups the banking and cash register operations.
type BanksService struct {
	client *Client
}

// List returns every bank account or cash register of the given type
// (GET /bank).
func (s *BanksService) List(ctx context.Context, bankType BankType) (*BankListResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).setString("bankType", bankType).values()
	out := &BankListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/bank", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single bank account or cash register by id (GET /bank/{id}).
func (s *BanksService) Get(ctx context.Context, id int64, bankType BankType) (*BankListResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).setString("bankType", bankType).values()
	out := &BankListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/bank/" + strconv.FormatInt(id, 10), query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Movements returns the paged list of movements of a bank account or cash
// register (GET /bank/{id}/movements).
func (s *BanksService) Movements(ctx context.Context, id int64, params BankMovementsParams) (*BankMovementListResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setString("bankType", params.BankType).
		setString("startDate", params.StartDate).
		setString("endDate", params.EndDate).
		setString("sortBy", params.SortBy).
		setIntOpt("page", int64(params.Page)).
		setFloatOpt("startAmount", params.StartAmount).
		setFloatOpt("endAmount", params.EndAmount).
		values()
	out := &BankMovementListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/bank/" + strconv.FormatInt(id, 10) + "/movements", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}
