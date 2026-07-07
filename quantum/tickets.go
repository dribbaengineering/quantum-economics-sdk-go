package quantum

import (
	"context"
	"net/http"
)

// Ticket is an expense ticket record.
type Ticket struct {
	Regid                 int64   `json:"regid" xml:"regid"`
	Codem                 int64   `json:"codem" xml:"codem"`
	Date                  string  `json:"date" xml:"date"`
	Description           string  `json:"description" xml:"description"`
	Amount                float64 `json:"amount" xml:"amount"`
	File                  int     `json:"file" xml:"file"`
	SWCFC                 int     `json:"swcfc" xml:"swcfc"`
	DateTransfer          string  `json:"dateTransfer" xml:"dateTransfer"`
	TypeOperation         string  `json:"typeOperation" xml:"typeOperation"`
	CJBN                  string  `json:"cjbn" xml:"cjbn"`
	InternalOperationCode int     `json:"internalOperationCode" xml:"internalOperationCode"`
	CashBank              int     `json:"cashBank" xml:"cashBank"`
	CollectionPaid        float64 `json:"collectionPaid" xml:"collectionPaid"`
	PaymentCollectionDate string  `json:"paymentCollectionDate" xml:"paymentCollectionDate"`
	BankCashName          string  `json:"bankCashName" xml:"bankCashName"`
	Account               string  `json:"account" xml:"account"`
	AccountName           string  `json:"accountName" xml:"accountName"`
	StrReference          string  `json:"strReference" xml:"strReference"`
	ReferenceNumber       int64   `json:"referenceNumber" xml:"referenceNumber"`

	// Ticketless (mobile-captured) metadata.
	Ticketless              bool   `json:"ticketless" xml:"ticketless"`
	TicketlessCategory      string `json:"ticketlessCategory" xml:"ticketlessCategory"`
	TicketlessPaymentMethod string `json:"ticketlessPaymentMethod" xml:"ticketlessPaymentMethod"`
	TicketlessUser          string `json:"ticketlessUser" xml:"ticketlessUser"`
}

// TicketListResponse is the envelope returned by GET /ticket.
type TicketListResponse struct {
	apiResponse
	Tickets []Ticket `json:"tickets" xml:"tickets"`
}

// TicketsParams are the filters for listing tickets. Year is required; the date
// range is optional.
type TicketsParams struct {
	Year      int
	StartDate string
	EndDate   string
	CompanyID int64
}

// TicketsService groups the ticket operations.
type TicketsService struct {
	client *Client
}

// List returns the company expense tickets (GET /ticket).
func (s *TicketsService) List(ctx context.Context, params TicketsParams) (*TicketListResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setStringOpt("startDate", params.StartDate).
		setStringOpt("endDate", params.EndDate).
		values()
	out := &TicketListResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/ticket", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// File returns the file attached to a ticket, Base64-encoded (GET /ticket/file).
func (s *TicketsService) File(ctx context.Context, ticketID int64) (string, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return "", err
	}
	q := newQuery().setInt("companyId", companyID).setInt("ticketId", ticketID).values()
	var out string
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/ticket/file", query: q}, &out); err != nil {
		return "", err
	}
	return out, nil
}
