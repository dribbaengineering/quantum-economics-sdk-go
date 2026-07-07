package quantum

import (
	"net/http"
	"strings"
)

// Authenticator applies credentials to an outgoing request. Implementations
// must be safe for concurrent use because a single client may issue many
// requests in parallel.
//
// The interface is intentionally tiny (interface segregation): the client only
// needs a way to authorize a request, nothing more. Custom implementations can
// be supplied through WithAuthenticator, which makes it possible to plug in
// token rotation, secret managers or test doubles without touching the client.
type Authenticator interface {
	// Authorize mutates the request so it carries valid credentials.
	Authorize(req *http.Request) error
}

// APIKeyAuthenticator authorizes requests using a Quantum API key.
//
// Quantum expects the header "Authorization: API-KEY <key>". The word
// "API-KEY", a single space and the key are concatenated for you, so callers
// pass only the raw key returned by the Quantum web application.
type APIKeyAuthenticator struct {
	// Key is the API key generated from Quantum (Mi Configuración → Conectores
	// → Q-API). It is used verbatim; do not prefix it with "API-KEY".
	Key string
}

// authorizationScheme is the fixed prefix Quantum requires in the header value.
const authorizationScheme = "API-KEY"

// Authorize implements Authenticator.
func (a APIKeyAuthenticator) Authorize(req *http.Request) error {
	key := strings.TrimSpace(a.Key)
	if key == "" {
		return ErrMissingAPIKey
	}
	req.Header.Set("Authorization", authorizationScheme+" "+key)
	return nil
}
