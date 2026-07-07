# Authentication

Every request to the Quantum API carries two credentials:

1. An **`Authorization`** header of the form `API-KEY <key>`.
2. A **`companyId`** query parameter identifying the company.

## Obtaining the credentials

1. In the Quantum web application go to **Mi Configuración → Conectores → Q-API**.
2. The list shows every company the user can access. The number in front of the
   company name is the **`companyId`**.
3. Press generate to create an **API key** for that company. It can be
   regenerated at any time (for example if it is compromised).

An API key is bound to a single company, so in practice one `Client` maps to one
company.

## Configuring the client

```go
client, err := quantum.NewClient(
    quantum.WithAPIKey("RGowc3lHV2NFSTMxVzZ6cm1wTnc2OFFJVjR6UjBiOTM="),
    quantum.WithCompanyID(28218),
)
```

- `WithAPIKey` sets the key; the SDK prepends `API-KEY ` for you, so pass only
  the raw key.
- `WithCompanyID` sets the default `companyId` sent with every call. Most
  parameter structs also expose a `CompanyID` field to override it per call
  (rarely needed, since keys are per-company).

If the API key is missing, `NewClient` returns `ErrMissingAPIKey`. If a call is
made without any company id (neither on the client nor the params), it fails
early with `ErrMissingCompanyID`.

## Custom authentication

`WithAPIKey` installs the default `APIKeyAuthenticator`. To take full control —
for example to load the key from a secret manager, or rotate it — provide your
own `Authenticator`:

```go
type Authenticator interface {
    Authorize(req *http.Request) error
}

client, _ := quantum.NewClient(
    quantum.WithAuthenticator(myAuthenticator{}),
    quantum.WithCompanyID(28218),
)
```

`WithAuthenticator` takes precedence over `WithAPIKey`.

## Keeping the key safe

Never hard-code the key. Read it from the environment or a secret store:

```go
quantum.WithAPIKey(os.Getenv("QUANTUM_API_KEY"))
```
