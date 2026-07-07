package quantum

import (
	"net/url"
	"strconv"
	"time"
)

// queryBuilder is a small helper for assembling url.Values without repeating
// nil/zero checks in every service method. Optional parameters are only written
// when they carry a meaningful value, matching the API's "omit to accept the
// default" behaviour.
type queryBuilder struct {
	v url.Values
}

func newQuery() *queryBuilder { return &queryBuilder{v: url.Values{}} }

// setString writes a mandatory string parameter unconditionally.
func (q *queryBuilder) setString(key, val string) *queryBuilder {
	q.v.Set(key, val)
	return q
}

// setStringOpt writes a string parameter only when it is non-empty.
func (q *queryBuilder) setStringOpt(key, val string) *queryBuilder {
	if val != "" {
		q.v.Set(key, val)
	}
	return q
}

// setInt writes a mandatory integer parameter unconditionally.
func (q *queryBuilder) setInt(key string, val int64) *queryBuilder {
	q.v.Set(key, strconv.FormatInt(val, 10))
	return q
}

// setIntOpt writes an integer parameter only when it is non-zero.
func (q *queryBuilder) setIntOpt(key string, val int64) *queryBuilder {
	if val != 0 {
		q.v.Set(key, strconv.FormatInt(val, 10))
	}
	return q
}

// setBoolOpt writes a boolean parameter only when the pointer is non-nil,
// allowing callers to distinguish "unset" from "false".
func (q *queryBuilder) setBoolOpt(key string, val *bool) *queryBuilder {
	if val != nil {
		q.v.Set(key, strconv.FormatBool(*val))
	}
	return q
}

// setFloatOpt writes a float parameter only when the pointer is non-nil.
func (q *queryBuilder) setFloatOpt(key string, val *float64) *queryBuilder {
	if val != nil {
		q.v.Set(key, strconv.FormatFloat(*val, 'f', -1, 64))
	}
	return q
}

// values returns the accumulated parameters.
func (q *queryBuilder) values() url.Values { return q.v }

// Bool is a convenience helper for taking the address of a boolean literal when
// populating optional *bool parameters.
//
//	params := quantum.TaxTypesParams{Valid: quantum.Bool(true)}
func Bool(b bool) *bool { return &b }

// Float64 is a convenience helper for taking the address of a float literal.
func Float64(f float64) *float64 { return &f }

// queryDateLayout is the layout used by date query parameters (dd-mm-yyyy),
// e.g. startDate=01-01-2025.
const queryDateLayout = "02-01-2006"

// bodyDateLayout is the layout used inside request/response bodies (dd/mm/yyyy),
// e.g. an installment date "22/09/2025".
const bodyDateLayout = "02/01/2006"

// FormatQueryDate formats a time as a Quantum date query parameter (dd-mm-yyyy).
func FormatQueryDate(t time.Time) string { return t.Format(queryDateLayout) }

// FormatBodyDate formats a time as a Quantum body date (dd/mm/yyyy).
func FormatBodyDate(t time.Time) string { return t.Format(bodyDateLayout) }

// ParseQueryDate parses a Quantum date query parameter (dd-mm-yyyy).
func ParseQueryDate(s string) (time.Time, error) { return time.Parse(queryDateLayout, s) }

// ParseBodyDate parses a Quantum body date (dd/mm/yyyy).
func ParseBodyDate(s string) (time.Time, error) { return time.Parse(bodyDateLayout, s) }
