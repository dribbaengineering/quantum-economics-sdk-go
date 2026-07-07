package quantum

// This file collects the small closed-set string codes used across the API as
// named constants. They document the accepted values inline and let callers
// avoid magic strings, while remaining plain string types so any value the
// server later introduces still round-trips.

// InvoiceType distinguishes issued (customer) from received (supplier) invoices.
type InvoiceType = string

const (
	// InvoiceTypeIssued is a customer invoice ("factura emitida").
	InvoiceTypeIssued InvoiceType = "C"
	// InvoiceTypeReceived is a supplier invoice ("factura recibida").
	InvoiceTypeReceived InvoiceType = "P"
)

// YesNo models the "S"/"N" flags used by tax fields.
type YesNo = string

const (
	Yes YesNo = "S"
	No  YesNo = "N"
)

// GoodsServices models the operation-kind flag on taxes ("B" goods, "S"
// services).
const (
	OperationGoods    = "B"
	OperationServices = "S"
)

// ReferenceType classifies an invoice line reference ("S" services, "M"
// articles).
const (
	ReferenceTypeService = "S"
	ReferenceTypeArticle = "M"
)

// AccountType selects the ledger family used by account and listing endpoints.
type AccountType = string

const (
	AccountTypeGeneral  AccountType = "G"
	AccountTypeCustomer AccountType = "C"
	AccountTypeProvider AccountType = "P"
)

// BankType distinguishes a bank account from a cash register.
type BankType = string

const (
	BankTypeBank BankType = "B"
	BankTypeCash BankType = "C"
)

// Language codes accepted by the document/PDF endpoints.
type Language = string

const (
	LanguageSpanish Language = "es"
	LanguageCatalan Language = "ca"
)

// TaxTypeKind classifies a defined tax type ("G" VAT/IVA general, "I" IRPF,
// "R" recargo).
type TaxTypeKind = string

const (
	TaxTypeKindVAT  TaxTypeKind = "G"
	TaxTypeKindIRPF TaxTypeKind = "I"
	TaxTypeKindRE   TaxTypeKind = "R"
)

// Period values accepted by tax and listing endpoints (monthly, quarterly,
// half-yearly and annual codes).
const (
	PeriodAnnual   = "0A"
	PeriodQ1       = "1T"
	PeriodQ2       = "2T"
	PeriodQ3       = "3T"
	PeriodQ4       = "4T"
	PeriodH1       = "1S"
	PeriodH2       = "2S"
	PeriodJanuary  = "01"
	PeriodDecember = "12"
)
