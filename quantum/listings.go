package quantum

import (
	"context"
	"net/http"
)

// RegistrationBookLine is a single income or expense line of the registration
// book.
type RegistrationBookLine struct {
	Total         float64 `json:"total" xml:"total"`
	Concept       string  `json:"concept" xml:"concept"`
	Description   string  `json:"description" xml:"description"`
	DeductionIVA  float64 `json:"deductionIVA" xml:"deductionIVA"`
	DeductionIRPF float64 `json:"deductionIRPF" xml:"deductionIRPF"`
	// Type is "INCOME" or "EXPENSES".
	Type string `json:"type" xml:"type"`
}

// RegistrationBookResponse is the envelope returned by GET
// /listing/registrationBook.
type RegistrationBookResponse struct {
	apiResponse
	Income          []RegistrationBookLine `json:"registrationBookLineIncome" xml:"registrationBookLineIncome"`
	Expenses        []RegistrationBookLine `json:"registrationBookLineExpenses" xml:"registrationBookLineExpenses"`
	TotalIncome     float64                `json:"totalIncome" xml:"totalIncome"`
	TotalExpenses   float64                `json:"totalExpenses" xml:"totalExpenses"`
	TotalDifference float64                `json:"totalDifference" xml:"totalDifference"`
	TotalPercentage float64                `json:"totalPercentage" xml:"totalPercentage"`
	NetReturn       float64                `json:"netReturn" xml:"netReturn"`
}

// AccountStatementParams are the parameters for the account statement (mayor).
// All fields except the amount bounds are required.
type AccountStatementParams struct {
	Year          int
	Type          AccountType
	StartDate     string
	EndDate       string
	StartAccount  string
	EndAccount    string
	Language      string
	LowestAmount  *float64
	HighestAmount *float64
	CompanyID     int64
}

// PeriodStatementParams are the shared parameters for the balance sheet, profit
// & loss and trial balance listings. All fields are required.
type PeriodStatementParams struct {
	Year int
	// PeriodType is "P2", "P" or "M".
	PeriodType string
	Starting   string
	Ending     string
	Language   string
	CompanyID  int64
}

// RegistrationBookParams are the parameters for the registration book. All
// fields are required.
type RegistrationBookParams struct {
	Year int
	// Period is one of the period codes (e.g. "1T", "01", "0A").
	Period string
	// Type is "IRPF" or "IVA".
	Type      string
	CompanyID int64
}

// ListingsService groups the accounting listing/report operations. The balance
// reports are returned as FileBase64Response (a Base64 PDF plus, when
// available, a structured JSON breakdown).
type ListingsService struct {
	client *Client
}

// AccountStatement returns the account statement, "mayor"
// (GET /listing/accountStatement).
func (s *ListingsService) AccountStatement(ctx context.Context, params AccountStatementParams) (*FileBase64Response, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setString("type", params.Type).
		setString("startDate", params.StartDate).
		setString("endDate", params.EndDate).
		setString("startAccount", params.StartAccount).
		setString("endAccount", params.EndAccount).
		setString("language", params.Language).
		setFloatOpt("lowestAmount", params.LowestAmount).
		setFloatOpt("highestAmount", params.HighestAmount).
		values()
	out := &FileBase64Response{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/listing/accountStatement", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// BalanceSheet returns the balance sheet (GET /listing/balanceSheet).
func (s *ListingsService) BalanceSheet(ctx context.Context, params PeriodStatementParams) (*FileBase64Response, error) {
	return s.periodStatement(ctx, "/listing/balanceSheet", params)
}

// ProfitAndLoss returns the profit & loss account (GET /listing/profitAndLoss).
func (s *ListingsService) ProfitAndLoss(ctx context.Context, params PeriodStatementParams) (*FileBase64Response, error) {
	return s.periodStatement(ctx, "/listing/profitAndLoss", params)
}

// TrialBalance returns the trial balance, "sumas y saldos"
// (GET /listing/trialBalance).
func (s *ListingsService) TrialBalance(ctx context.Context, params PeriodStatementParams) (*FileBase64Response, error) {
	return s.periodStatement(ctx, "/listing/trialBalance", params)
}

// RegistrationBook returns the registration book (GET /listing/registrationBook).
func (s *ListingsService) RegistrationBook(ctx context.Context, params RegistrationBookParams) (*RegistrationBookResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setString("period", params.Period).
		setString("type", params.Type).
		values()
	out := &RegistrationBookResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/listing/registrationBook", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *ListingsService) periodStatement(ctx context.Context, path string, params PeriodStatementParams) (*FileBase64Response, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setString("periodType", params.PeriodType).
		setString("starting", params.Starting).
		setString("ending", params.Ending).
		setString("language", params.Language).
		values()
	out := &FileBase64Response{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: path, query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}
