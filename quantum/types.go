package quantum

// IDResponse is returned by the create endpoints (customers, providers,
// invoices, proforma, delivery notes, portfolio movements). ID is the
// identifier of the newly created or matched resource; for invoice-family
// endpoints it is the identifier of the created invoice.
type IDResponse struct {
	apiResponse
	ID int64 `json:"id" xml:"id"`
}
