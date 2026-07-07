package quantum

// ResponseError is the "error" object present in every Quantum response
// envelope. An ErrorCode of 0 with message "OK" indicates success.
type ResponseError struct {
	Message   string `json:"message" xml:"message"`
	ErrorCode int    `json:"errorCode" xml:"errorCode"`
}

// apiResponse is embedded in every typed response. It carries the fields shared
// by all Quantum envelopes and, via apiError, exposes the business-level error
// in a uniform way. Embedding it (rather than duplicating the fields) keeps the
// domain types focused on their own payload.
type apiResponse struct {
	Err        *ResponseError `json:"error,omitempty" xml:"error,omitempty"`
	APIVersion float64        `json:"apiVersion,omitempty" xml:"apiVersion,omitempty"`
}

// apiError implements the errorCarrier interface. It returns a non-nil error
// only when Quantum reported a non-zero error code.
func (r apiResponse) apiError() *ResponseError {
	if r.Err != nil && r.Err.ErrorCode != 0 {
		return r.Err
	}
	return nil
}

// errorCarrier is satisfied by any response type able to report a Quantum
// business error. Every response embeds apiResponse, so pointers to those
// types satisfy this interface automatically. The client uses it to detect
// envelope-level failures regardless of the concrete response type.
type errorCarrier interface {
	apiError() *ResponseError
}
