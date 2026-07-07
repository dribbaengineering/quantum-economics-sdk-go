package quantum

import "encoding/json"

// RawResponse is returned by endpoints whose response schema is not reliably
// described by the upstream API specification (some Worker, Labour, Diary and
// DUA listing endpoints reuse an unrelated schema in the spec). It still gives
// you the two things that always matter — uniform error handling via the
// embedded envelope, and the complete payload — without pretending to a shape
// that may be inaccurate.
//
// Decode the payload into your own type with Decode:
//
//	res, err := client.Workers.WorkTime(ctx, params)
//	if err != nil { return err }
//	var records []quantum.WorkTimeRegistry
//	if err := res.Decode(&records); err != nil { return err }
type RawResponse struct {
	apiResponse
	// Payload is the raw response body (JSON). It is always populated on
	// success, so callers can inspect or decode it however they need.
	Payload json.RawMessage
}

// UnmarshalJSON captures the entire body in Payload while still populating the
// embedded envelope so business errors are detected like any other response.
func (r *RawResponse) UnmarshalJSON(data []byte) error {
	r.Payload = append(r.Payload[:0], data...)

	// Decode the envelope fields via an alias to avoid recursing into this
	// method.
	var env struct {
		Err        *ResponseError `json:"error"`
		APIVersion float64        `json:"apiVersion"`
	}
	if err := json.Unmarshal(data, &env); err != nil {
		return err
	}
	r.Err = env.Err
	r.APIVersion = env.APIVersion
	return nil
}

// Decode unmarshals the raw payload into v, a non-nil pointer.
func (r *RawResponse) Decode(v any) error {
	return json.Unmarshal(r.Payload, v)
}
