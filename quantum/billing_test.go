package quantum

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestInvoiceCreateEncodesBody(t *testing.T) {
	var body map[string]any
	var gotContentType, gotMethod, gotPath string
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		gotMethod = r.Method
		gotPath = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = io.WriteString(w, `{"error":{"message":"OK","errorCode":0},"id":501}`)
	})

	res, err := client.Invoices.Create(context.Background(), Invoice{
		Type:                    InvoiceTypeIssued,
		CustomerProviderID:      2,
		TotalAmountWithoutTaxes: 100,
		TotalAmount:             121,
		Line: []InvoiceLine{{
			Description: "Order 1",
			Quantity:    1,
			Amount:      121,
			Base:        100,
			Percentage:  21,
			TaxCode:     21,
		}},
		Installment:    []InvoiceInstallment{{Date: "22/09/2025", Amount: 121}},
		DescriptionSII: "test invoice",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if res.ID != 501 {
		t.Errorf("id = %d, want 501", res.ID)
	}
	if gotMethod != http.MethodPost || gotPath != "/contabilidad/ws/invoice" {
		t.Errorf("got %s %s", gotMethod, gotPath)
	}
	if gotContentType != "application/json" {
		t.Errorf("Content-Type = %q", gotContentType)
	}
	if body["type"] != "C" || body["customerProviderId"].(float64) != 2 {
		t.Errorf("unexpected body: %v", body)
	}
	// omitempty must drop zero-valued optional fields.
	if _, present := body["rental"]; present {
		t.Errorf("expected omitempty to drop rental, body: %v", body)
	}
}

func TestInvoiceGetUnwrapsFirst(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"error":{"errorCode":0},"invoices":[{"id":7,"type":"C","totalAmount":121}]}`)
	})
	inv, err := client.Invoices.Get(context.Background(), 7)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if inv.ID != 7 || inv.TotalAmount != 121 {
		t.Errorf("unexpected invoice: %+v", inv)
	}
}

func TestInvoiceGetNotFound(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"error":{"errorCode":0},"invoices":[]}`)
	})
	_, err := client.Invoices.Get(context.Background(), 7)
	if err != ErrNotFound {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestInvoiceListQueryParams(t *testing.T) {
	var q url.Values
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		q = r.URL.Query()
		_, _ = io.WriteString(w, `{"error":{"errorCode":0},"invoices":[]}`)
	})
	_, err := client.Invoices.List(context.Background(), ListInvoicesParams{
		Type:      InvoiceTypeReceived,
		StartDate: "01-01-2025",
		EndDate:   "31-12-2025",
		Page:      2,
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if q.Get("type") != "P" || q.Get("startDate") != "01-01-2025" || q.Get("page") != "2" {
		t.Errorf("unexpected query: %v", q)
	}
}

func TestCustomerGetByNIF(t *testing.T) {
	var gotPath string
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = io.WriteString(w, `{"error":{"errorCode":0},"customers":[{"customerId":"10","nif":"B00000000","name":"ACME"}]}`)
	})
	c, err := client.Customers.GetByNIF(context.Background(), "B00000000")
	if err != nil {
		t.Fatalf("GetByNIF: %v", err)
	}
	if c.CustomerID != "10" || c.Name != "ACME" {
		t.Errorf("unexpected customer: %+v", c)
	}
	if gotPath != "/contabilidad/ws/customer/nif/B00000000" {
		t.Errorf("path = %q", gotPath)
	}
}

func TestXMLContentNegotiation(t *testing.T) {
	var gotContentType, gotAccept, bodyStr string
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		gotAccept = r.Header.Get("Accept")
		raw, _ := io.ReadAll(r.Body)
		bodyStr = string(raw)
		w.Header().Set("Content-Type", "application/xml")
		_, _ = io.WriteString(w, `<Id><error><errorCode>0</errorCode></error><id>9</id></Id>`)
	}, WithXML())

	res, err := client.Customers.Create(context.Background(), Customer{NIF: "00000011s", Name: "name", CityCode: "08100"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if res.ID != 9 {
		t.Errorf("id = %d, want 9", res.ID)
	}
	if gotContentType != "application/xml" || gotAccept != "application/xml" {
		t.Errorf("content negotiation headers wrong: ct=%q accept=%q", gotContentType, gotAccept)
	}
	if !strings.Contains(bodyStr, "<Customer>") || !strings.Contains(bodyStr, "<nif>00000011s</nif>") {
		t.Errorf("unexpected XML body: %s", bodyStr)
	}
}

func TestRawResponseDecode(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"error":{"errorCode":0},"records":[{"regid":1,"type":2}]}`)
	})
	res, err := client.Workers.WorkTime(context.Background(), WorkTimeListParams{Year: 2025})
	if err != nil {
		t.Fatalf("WorkTime: %v", err)
	}
	var payload struct {
		Records []WorkTimeRegistry `json:"records"`
	}
	if err := res.Decode(&payload); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if len(payload.Records) != 1 || payload.Records[0].Regid != 1 {
		t.Errorf("unexpected payload: %+v", payload)
	}
}
