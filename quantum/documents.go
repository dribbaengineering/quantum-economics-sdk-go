package quantum

import "encoding/json"

// DocumentResponse carries a document returned inline, Base64-encoded, together
// with its filename. It is used by the "document" endpoints (invoice, proforma,
// labour, payroll) that embed the file in a JSON envelope.
type DocumentResponse struct {
	apiResponse
	// Document is the file contents, Base64-encoded.
	Document string `json:"document" xml:"document"`
	// Filename is the suggested filename, including extension.
	Filename string `json:"filename" xml:"filename"`
}

// FileBase64Response carries a generated report both as a Base64 blob (usually a
// PDF) and, when available, as a structured JSON breakdown. It is used by the
// listing endpoints (account statement, balance sheet, profit & loss, trial
// balance).
type FileBase64Response struct {
	apiResponse
	// Base64 is the rendered document (PDF), Base64-encoded.
	Base64 string `json:"base64" xml:"base64"`
	// JSON is the structured representation of the same report, when provided.
	JSON []PDFReportEntry `json:"json" xml:"json"`
}

// URLResponse carries a temporary URL pointing at a generated document.
type URLResponse struct {
	apiResponse
	URL string `json:"url" xml:"url"`
}

// PDFReportEntry is a single node of the structured breakdown accompanying a
// FileBase64Response. Entries can nest through Sublevel to represent report
// hierarchies.
type PDFReportEntry struct {
	Title            string `json:"title" xml:"title"`
	StartPeriodMonth int    `json:"startPeriodMonth" xml:"startPeriodMonth"`
	EndPeriodMonth   int    `json:"endPeriodMonth" xml:"endPeriodMonth"`
	// Fields holds report-specific key/value figures; its shape depends on the
	// report, so it is exposed as raw JSON for the caller to interpret.
	Fields   json.RawMessage  `json:"fields" xml:"fields"`
	Sublevel []PDFReportEntry `json:"sublevel" xml:"sublevel"`
}
