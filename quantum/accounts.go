package quantum

import (
	"context"
	"net/http"
)

// Account represents a ledger account with its balances.
type Account struct {
	// ID is the account code.
	ID             string         `json:"id" xml:"id"`
	Name           string         `json:"name" xml:"name"`
	CurrentBalance float64        `json:"currentBalance" xml:"currentBalance"`
	Debit          float64        `json:"debit" xml:"debit"`
	Credit         float64        `json:"credit" xml:"credit"`
	MonthlyBalance []MonthBalance `json:"monthlyBalance" xml:"monthlyBalance"`
}

// MonthBalance is the balance of an account for a single month.
type MonthBalance struct {
	Month   int     `json:"month" xml:"month"`
	Year    int     `json:"year" xml:"year"`
	Balance float64 `json:"balance" xml:"balance"`
}

// AccountListResponse is the envelope returned by GET /account.
type AccountListResponse struct {
	apiResponse
	Accounts       []Account `json:"getaccounts" xml:"getaccounts"`
	FirstAccountID string    `json:"firstAccountId" xml:"firstAccountId"`
	LastAccountID  string    `json:"lastAccountId" xml:"lastAccountId"`
}

// AccountingPlanEntry is a row of the accounting plan listing.
type AccountingPlanEntry struct {
	AccountNumber string  `json:"account_number" xml:"account_number"`
	AccountTitle  string  `json:"account_title" xml:"account_title"`
	Month         int     `json:"month" xml:"month"`
	Debits        float64 `json:"debits" xml:"debits"`
	Credits       float64 `json:"credits" xml:"credits"`
	NetMovement   float64 `json:"net_movement" xml:"net_movement"`
}

// AccountingPlanResponse is the envelope returned by GET /account/accountingPlan.
type AccountingPlanResponse struct {
	apiResponse
	Data    []AccountingPlanEntry `json:"data" xml:"data"`
	Company string                `json:"company" xml:"company"`
}

// AccountSearchParams are the parameters for searching accounts.
type AccountSearchParams struct {
	// Year is required.
	Year int
	// AccountType is required: AccountTypeGeneral/Customer/Provider.
	AccountType AccountType
	// Text optionally filters by name/code.
	Text string
	// CompanyID overrides the client's default company id (optional).
	CompanyID int64
}

// AccountingPlanParams are the parameters for the accounting plan listing.
type AccountingPlanParams struct {
	// Year is required.
	Year int
	// AdmitEntries optionally filters accounts admitting entries ("S"/"N"/"T").
	AdmitEntries    string
	ShowEmptyValues *bool
	DetailByMonth   *bool
	Formated        *bool
	StartAccount    string
	EndAccount      string
	// AccountType optionally filters ("G"/"C"/"P"/"A").
	AccountType string
	// CompanyID overrides the client's default company id (optional).
	CompanyID int64
}

// AccountDefinitionParams are the parameters for the account definitions
// endpoint.
type AccountDefinitionParams struct {
	// TpPGC, Country and Paquete are all required.
	TpPGC   string
	Country string
	Paquete string
	// CompanyID overrides the client's default company id (optional).
	CompanyID int64
}

// AccountsService groups the ledger account operations.
type AccountsService struct {
	client *Client
}

// Search returns the accounts matching the given filters (GET /account).
func (s *AccountsService) Search(ctx context.Context, params AccountSearchParams) (*AccountListResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setString("accountType", params.AccountType).
		setStringOpt("text", params.Text).
		values()
	out := &AccountListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/account", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// AccountingPlan returns the accounting plan (GET /account/accountingPlan).
func (s *AccountsService) AccountingPlan(ctx context.Context, params AccountingPlanParams) (*AccountingPlanResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setStringOpt("admitEntries", params.AdmitEntries).
		setBoolOpt("showEmptyValues", params.ShowEmptyValues).
		setBoolOpt("detailByMonth", params.DetailByMonth).
		setBoolOpt("formated", params.Formated).
		setStringOpt("startAccount", params.StartAccount).
		setStringOpt("endAccount", params.EndAccount).
		setStringOpt("accountType", params.AccountType).
		values()
	out := &AccountingPlanResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/account/accountingPlan", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// SearchFull returns the full, legacy-shaped account records
// (GET /account/full). The payload is a wide legacy table and is returned raw.
func (s *AccountsService) SearchFull(ctx context.Context, params AccountSearchParams) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setString("accountType", params.AccountType).
		values()
	out := &RawResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/account/full", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Definitions returns the full account definitions
// (GET /account/definition/full). The payload is a wide legacy table and is
// returned raw.
func (s *AccountsService) Definitions(ctx context.Context, params AccountDefinitionParams) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setString("tpPGC", params.TpPGC).
		setString("country", params.Country).
		setString("paquete", params.Paquete).
		values()
	out := &RawResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/account/definition/full", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}
