package quantum

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultBaseURL is the production Quantum API root, including the /ws context
// path. All service paths are resolved relative to it.
const DefaultBaseURL = "https://app.quantumeconomics.es/contabilidad/ws/"

// defaultUserAgent identifies this SDK to the server.
const defaultUserAgent = "quantum-economics-sdk-go/1.0"

// defaultTimeout is applied to the default HTTP client when the caller does not
// provide one.
const defaultTimeout = 30 * time.Second

// Client is the entry point of the SDK. It is safe for concurrent use by
// multiple goroutines once constructed. Business functionality is reached
// through the service fields, each of which owns a single API domain.
//
// A Client is immutable after construction: configure it entirely through the
// options passed to NewClient.
type Client struct {
	baseURL        *url.URL
	httpClient     Doer
	auth           Authenticator
	companyID      int64
	userAgent      string
	logger         Logger
	requestFormat  ContentType
	responseFormat ContentType
	defaultHeaders map[string]string

	// Services grouped by API domain.
	Invoices      *InvoicesService
	Proforma      *ProformaService
	Customers     *CustomersService
	Providers     *ProvidersService
	Companies     *CompaniesService
	Banks         *BanksService
	Accounts      *AccountsService
	Taxes         *TaxesService
	TaxTypes      *TaxTypesService
	Listings      *ListingsService
	Labour        *LabourService
	Workers       *WorkersService
	Tickets       *TicketsService
	Diaries       *DiariesService
	DUA           *DUAService
	Risk          *RiskService
	Portfolio     *PortfolioService
	DeliveryNotes *DeliveryNotesService
	QuantumBI     *QuantumBIService
}

// NewClient builds a Client from the given options. At minimum an API key (or a
// custom Authenticator) is required; a company id is strongly recommended since
// almost every endpoint needs one.
func NewClient(opts ...Option) (*Client, error) {
	cfg := &config{
		baseURL:        DefaultBaseURL,
		userAgent:      defaultUserAgent,
		requestFormat:  ContentTypeJSON,
		responseFormat: ContentTypeJSON,
		requestTimeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	base, err := parseBaseURL(cfg.baseURL)
	if err != nil {
		return nil, err
	}

	auth := cfg.auth
	if auth == nil {
		if strings.TrimSpace(cfg.apiKey) == "" {
			return nil, ErrMissingAPIKey
		}
		auth = APIKeyAuthenticator{Key: cfg.apiKey}
	}

	httpClient := cfg.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: cfg.requestTimeout}
	}

	logger := cfg.logger
	if logger == nil {
		logger = nopLogger{}
	}

	c := &Client{
		baseURL:        base,
		httpClient:     httpClient,
		auth:           auth,
		companyID:      cfg.companyID,
		userAgent:      cfg.userAgent,
		logger:         logger,
		requestFormat:  cfg.requestFormat,
		responseFormat: cfg.responseFormat,
		defaultHeaders: cfg.defaultHeaders,
	}

	c.registerServices()
	return c, nil
}

// registerServices instantiates every service. Each service embeds a back
// reference to the client so it can reuse the shared request pipeline.
func (c *Client) registerServices() {
	c.Invoices = &InvoicesService{client: c}
	c.Proforma = &ProformaService{client: c}
	c.Customers = &CustomersService{client: c}
	c.Providers = &ProvidersService{client: c}
	c.Companies = &CompaniesService{client: c}
	c.Banks = &BanksService{client: c}
	c.Accounts = &AccountsService{client: c}
	c.Taxes = &TaxesService{client: c}
	c.TaxTypes = &TaxTypesService{client: c}
	c.Listings = &ListingsService{client: c}
	c.Labour = &LabourService{client: c}
	c.Workers = &WorkersService{client: c}
	c.Tickets = &TicketsService{client: c}
	c.Diaries = &DiariesService{client: c}
	c.DUA = &DUAService{client: c}
	c.Risk = &RiskService{client: c}
	c.Portfolio = &PortfolioService{client: c}
	c.DeliveryNotes = &DeliveryNotesService{client: c}
	c.QuantumBI = &QuantumBIService{client: c}
}

// CompanyID returns the default company id configured on the client.
func (c *Client) CompanyID() int64 { return c.companyID }

// BaseURL returns a copy of the configured base URL.
func (c *Client) BaseURL() string { return c.baseURL.String() }

// resolveCompanyID returns the company id to use for a request: the per-call
// override when non-zero, otherwise the client default. It errors when neither
// is set, turning a common misconfiguration into an explicit, early failure.
func (c *Client) resolveCompanyID(override int64) (int64, error) {
	if override != 0 {
		return override, nil
	}
	if c.companyID != 0 {
		return c.companyID, nil
	}
	return 0, ErrMissingCompanyID
}

// parseBaseURL validates and normalizes the base URL, ensuring it ends with a
// trailing slash so relative path resolution works predictably.
func parseBaseURL(raw string) (*url.URL, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return nil, ErrInvalidBaseURL
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, ErrInvalidBaseURL
	}
	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}
	return u, nil
}
