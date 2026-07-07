package quantum

import "time"

// Option configures a Client. Options follow the functional-options pattern so
// the client can grow new knobs without breaking existing call sites and
// without exposing its internal fields (open/closed principle).
type Option func(*config)

// config holds everything NewClient needs to assemble a Client. It is unexported
// so the only supported way to build a client is through NewClient + Option.
type config struct {
	baseURL        string
	apiKey         string
	auth           Authenticator
	companyID      int64
	httpClient     Doer
	userAgent      string
	logger         Logger
	requestFormat  ContentType
	responseFormat ContentType
	requestTimeout time.Duration
	defaultHeaders map[string]string
}

// WithBaseURL overrides the API base URL. Defaults to DefaultBaseURL. Useful for
// pointing the client at a staging environment or a local mock server.
func WithBaseURL(rawURL string) Option {
	return func(c *config) { c.baseURL = rawURL }
}

// WithAPIKey sets the Quantum API key used for the default APIKeyAuthenticator.
// Ignored if WithAuthenticator is also supplied.
func WithAPIKey(key string) Option {
	return func(c *config) { c.apiKey = key }
}

// WithAuthenticator injects a custom Authenticator, taking full control of how
// requests are authorized. Overrides WithAPIKey.
func WithAuthenticator(a Authenticator) Option {
	return func(c *config) { c.auth = a }
}

// WithCompanyID sets the default company id sent as the companyId query
// parameter. Individual calls may override it through their params.
func WithCompanyID(companyID int64) Option {
	return func(c *config) { c.companyID = companyID }
}

// WithHTTPClient injects the HTTP client used to execute requests. Any type
// satisfying Doer works, including a customized *http.Client or a test double.
func WithHTTPClient(d Doer) Option {
	return func(c *config) { c.httpClient = d }
}

// WithHTTPTimeout sets the timeout of the default *http.Client. It has no effect
// when a custom Doer is supplied through WithHTTPClient.
func WithHTTPTimeout(d time.Duration) Option {
	return func(c *config) { c.requestTimeout = d }
}

// WithUserAgent overrides the User-Agent header sent with every request.
func WithUserAgent(ua string) Option {
	return func(c *config) { c.userAgent = ua }
}

// WithLogger installs a logger for request/response diagnostics. Defaults to a
// no-op logger.
func WithLogger(l Logger) Option {
	return func(c *config) { c.logger = l }
}

// WithRequestFormat selects the wire format used to encode request bodies
// (JSON by default). Quantum accepts both application/json and application/xml.
func WithRequestFormat(ct ContentType) Option {
	return func(c *config) { c.requestFormat = ct }
}

// WithResponseFormat sets the Accept header, asking Quantum to answer in the
// given format (JSON by default).
func WithResponseFormat(ct ContentType) Option {
	return func(c *config) { c.responseFormat = ct }
}

// WithXML is a shortcut that sends and receives XML.
func WithXML() Option {
	return func(c *config) {
		c.requestFormat = ContentTypeXML
		c.responseFormat = ContentTypeXML
	}
}

// WithDefaultHeader adds a header sent with every request. Repeated calls
// accumulate; later calls override earlier ones for the same key.
func WithDefaultHeader(key, value string) Option {
	return func(c *config) {
		if c.defaultHeaders == nil {
			c.defaultHeaders = map[string]string{}
		}
		c.defaultHeaders[key] = value
	}
}
