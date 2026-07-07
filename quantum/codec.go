package quantum

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

// ContentType identifies a wire format understood by the Quantum API. The API
// accepts and returns both JSON and XML; JSON is the default and recommended
// format.
type ContentType string

const (
	// ContentTypeJSON selects application/json.
	ContentTypeJSON ContentType = "application/json"
	// ContentTypeXML selects application/xml.
	ContentTypeXML ContentType = "application/xml"
)

// mediaType returns the value for the Content-Type / Accept headers.
func (c ContentType) mediaType() string {
	if c == "" {
		return string(ContentTypeJSON)
	}
	return string(c)
}

// codec encapsulates the marshalling strategy for a single wire format. Keeping
// encode/decode behind an interface isolates the rest of the client from the
// concrete serialization library (single responsibility + dependency
// inversion) and makes it trivial to add another format later.
type codec interface {
	// Marshal serializes v into request-body bytes.
	Marshal(v any) ([]byte, error)
	// Unmarshal deserializes data into v, which is always a non-nil pointer.
	Unmarshal(data []byte, v any) error
}

// codecFor returns the codec matching a content type. Unknown values fall back
// to JSON, matching the server behaviour ("Si no se indica por defecto la
// respuesta será en formato JSON").
func codecFor(ct ContentType) codec {
	switch ct {
	case ContentTypeXML:
		return xmlCodec{}
	default:
		return jsonCodec{}
	}
}

type jsonCodec struct{}

func (jsonCodec) Marshal(v any) ([]byte, error) { return json.Marshal(v) }

func (jsonCodec) Unmarshal(data []byte, v any) error {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

type xmlCodec struct{}

func (xmlCodec) Marshal(v any) ([]byte, error) {
	body, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	// Prepend the declaration Quantum's examples use.
	return append([]byte(xml.Header), body...), nil
}

func (xmlCodec) Unmarshal(data []byte, v any) error {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	return xml.Unmarshal(data, v)
}

// decodeInto reads all of r and decodes it into v using the given codec,
// returning a helpful error that includes a snippet of the payload on failure.
func decodeInto(c codec, r io.Reader, v any) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	if v == nil {
		return nil
	}
	if err := c.Unmarshal(data, v); err != nil {
		return &DecodeError{Err: err, Body: data}
	}
	return nil
}
