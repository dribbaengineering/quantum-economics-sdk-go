package quantum

import (
	"context"
	"net/http"
)

// JournalEntry is a single accounting journal entry ("apunte").
type JournalEntry struct {
	MovementDate  string  `json:"movementDate" xml:"movementDate"`
	Entry         int     `json:"entry" xml:"entry"`
	Diary         int     `json:"diary" xml:"diary"`
	Concept       string  `json:"concept" xml:"concept"`
	AccountCode   int64   `json:"accountCode" xml:"accountCode"`
	InvoiceNumber string  `json:"invoiceNumber" xml:"invoiceNumber"`
	Debits        float64 `json:"debits" xml:"debits"`
	Credits       float64 `json:"credits" xml:"credits"`
	// AccountType is "G", "C" or "P".
	AccountType string `json:"accountType" xml:"accountType"`
}

// JournalEntriesResponse is the envelope returned by GET /diaries/obtainDiaries.
type JournalEntriesResponse struct {
	apiResponse
	Entries []JournalEntry `json:"entries" xml:"entries"`
}

// ObtainDiariesParams are the filters for the journal entries query. Only the
// company id is required.
type ObtainDiariesParams struct {
	StartMovementDate string
	EndMovementDate   string
	Entry             int
	Diary             int
	AccountCode       int64
	// AccountType optionally filters by "G", "C" or "P".
	AccountType   string
	Concept       string
	InvoiceNumber string
	CompanyID     int64
}

// DiariesService groups the journal ("diario") operations.
type DiariesService struct {
	client *Client
}

// ObtainDiaries returns the journal entries matching the filters
// (GET /diaries/obtainDiaries).
func (s *DiariesService) ObtainDiaries(ctx context.Context, params ObtainDiariesParams) (*JournalEntriesResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setStringOpt("startMovementDate", params.StartMovementDate).
		setStringOpt("endMovementDate", params.EndMovementDate).
		setIntOpt("entry", int64(params.Entry)).
		setIntOpt("diary", int64(params.Diary)).
		setIntOpt("accountCode", params.AccountCode).
		setStringOpt("accountType", params.AccountType).
		setStringOpt("concept", params.Concept).
		setStringOpt("invoiceNumber", params.InvoiceNumber).
		values()
	out := &JournalEntriesResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/diaries/obtainDiaries", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DiaryDefinition returns the diary definitions for a year
// (GET /diaries/diaryDefinition). The response schema is not reliably specified
// upstream, so it is returned raw.
func (s *DiariesService) DiaryDefinition(ctx context.Context, year int) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).setInt("year", int64(year)).values()
	out := &RawResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/diaries/diaryDefinition", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}
