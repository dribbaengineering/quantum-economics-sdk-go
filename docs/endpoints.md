# Endpoint coverage

Every documented Quantum API endpoint is reachable through a service method.
Methods take a `context.Context` first and return a typed response (or a
`RawResponse` for the few under-specified endpoints — see
[architecture.md](architecture.md#response-typing-policy)).

Paths are relative to `https://app.quantumeconomics.es/contabilidad/ws`.

## Invoices — `client.Invoices`

| Method | HTTP | Path |
| --- | --- | --- |
| `Create` | POST | `/invoice` |
| `CreateWithDocument` | POST | `/invoice/invoiceWithFile` |
| `CreateWithFacturae` | POST | `/invoice/invoiceWithFacturae` |
| `List` | GET | `/invoice` |
| `ListFull` | GET | `/invoice/full` |
| `Get` | GET | `/invoice/{id}` |
| `Dashboard` | GET | `/invoice/dashboard` |
| `PaymentState` | GET | `/invoice/state/{id}` |
| `PDFURL` | GET | `/invoice/pdf` |
| `PDFBase64` | GET | `/invoice/pdfB64` |
| `Document` | GET | `/invoice/document` |

## Pro forma invoices — `client.Proforma`

| Method | HTTP | Path |
| --- | --- | --- |
| `Create` | POST | `/proforma` |
| `List` | GET | `/proforma` |
| `ListFull` | GET | `/proforma/full` |
| `Get` | GET | `/proforma/{id}` |
| `Dashboard` | GET | `/proforma/dashboard` |
| `PDFURL` | GET | `/proforma/pdf` |
| `PDFBase64` | GET | `/proforma/pdfB64` |
| `Document` | GET | `/proforma/document` |

## Customers — `client.Customers`

| Method | HTTP | Path |
| --- | --- | --- |
| `List` | GET | `/customer` |
| `GetByID` | GET | `/customer/{id}` |
| `GetByNIF` | GET | `/customer/nif/{id}` |
| `Create` | POST | `/customer` |
| `Update` | PUT | `/customer` |

## Suppliers — `client.Providers`

| Method | HTTP | Path |
| --- | --- | --- |
| `List` | GET | `/provider` |
| `GetByID` | GET | `/provider/{id}` |
| `GetByNIF` | GET | `/provider/nif/{id}` |
| `Create` | POST | `/provider` |

## Companies — `client.Companies`

| Method | HTTP | Path |
| --- | --- | --- |
| `List` | GET | `/company` |
| `GetFull` | GET | `/company/full` |
| `State` | GET | `/company/state` |

## Banking — `client.Banks`

| Method | HTTP | Path |
| --- | --- | --- |
| `List` | GET | `/bank` |
| `Get` | GET | `/bank/{id}` |
| `Movements` | GET | `/bank/{id}/movements` |

## Accounts — `client.Accounts`

| Method | HTTP | Path |
| --- | --- | --- |
| `Search` | GET | `/account` |
| `AccountingPlan` | GET | `/account/accountingPlan` |
| `SearchFull` | GET | `/account/full` |
| `Definitions` | GET | `/account/definition/full` |

## Taxes — `client.Taxes`

| Method | HTTP | Path |
| --- | --- | --- |
| `List` | GET | `/tax` |
| `Detail` | GET | `/tax/detail` |
| `Full` | GET | `/tax/full` |
| `PDF` | GET | `/tax/pdf` |
| `PDFBase64` | GET | `/tax/pdfB64` |
| `PDFOld` | GET | `/tax/pdfOld` |

## Tax types — `client.TaxTypes`

| Method | HTTP | Path |
| --- | --- | --- |
| `List` | GET | `/taxesTypes` |
| `Full` | GET | `/taxesTypes/full` |

## Listings — `client.Listings`

| Method | HTTP | Path |
| --- | --- | --- |
| `AccountStatement` | GET | `/listing/accountStatement` |
| `BalanceSheet` | GET | `/listing/balanceSheet` |
| `ProfitAndLoss` | GET | `/listing/profitAndLoss` |
| `TrialBalance` | GET | `/listing/trialBalance` |
| `RegistrationBook` | GET | `/listing/registrationBook` |

## Labour — `client.Labour`

| Method | HTTP | Path |
| --- | --- | --- |
| `Summary` | GET | `/labour` |
| `Detail` | GET | `/labour/detail` |
| `Document` | GET | `/labour/document` |
| `Payroll` | GET | `/labour/payroll` |
| `PayrollByCompany` | GET | `/labour/payroll/company` |
| `PayrollDocument` | GET | `/labour/payroll/document` |
| `ValidatePayroll` | POST | `/labour/payroll/validate` |
| `ManageBasicDocuments` | GET | `/labour/manageBasicDocuments` |
| `ManageCertDocuments` | GET | `/labour/manageCertDocuments` |
| `ManageDocuments` | GET | `/labour/manageDocuments` |
| `OtherDocuments` | GET | `/labour/otherDocuments` |

## Workers — `client.Workers`

| Method | HTTP | Path |
| --- | --- | --- |
| `Absences` | GET | `/worker/absences` |
| `AbsencesSummary` | GET | `/worker/absences/summary` |
| `AbsenceFile` | GET | `/worker/absences/file` |
| `WorkingDays` | GET | `/worker/absences/getWorkingDays` |
| `SendAbsence` | POST | `/worker/absences/send` |
| `DeleteAbsence` | POST | `/worker/absences/delete` |
| `ValidateAbsences` | POST | `/worker/absences/validate` |
| `Calendar` | GET | `/worker/calendar` |
| `CalendarAll` | GET | `/worker/calendar/all` |
| `CalendarIncidents` | GET | `/worker/calendar/incidents` |
| `WorkTime` | GET | `/worker/workTime` |
| `WorkTimeWeek` | GET | `/worker/workTime/week` |
| `WorkTimePending` | GET | `/worker/workTime/pending/worker` |
| `WorkTimePre` | GET | `/worker/workTime/pre` |
| `WorkTimeStatus` | GET | `/worker/workTime/status` |
| `WorkTimeSummary` | GET | `/worker/workTime/summary` |
| `SendWorkTime` | POST | `/worker/workTime/send` |
| `EditWorkTime` | POST | `/worker/workTime/edit` |
| `WorkTimeReminder` | POST | `/worker/workTime/reminder` |
| `ValidateWorkTime` | POST | `/worker/workTime/validate` |
| `ValidateWorkTimeByWorker` | POST | `/worker/workTime/validate/worker` |
| `ValidateAllWorkTime` | POST | `/worker/workTime/validateAll` |

## Tickets — `client.Tickets`

| Method | HTTP | Path |
| --- | --- | --- |
| `List` | GET | `/ticket` |
| `File` | GET | `/ticket/file` |

## Journal — `client.Diaries`

| Method | HTTP | Path |
| --- | --- | --- |
| `ObtainDiaries` | GET | `/diaries/obtainDiaries` |
| `DiaryDefinition` | GET | `/diaries/diaryDefinition` |

## DUA — `client.DUA`

| Method | HTTP | Path |
| --- | --- | --- |
| `List` | GET | `/dua` |

## Risk — `client.Risk`

| Method | HTTP | Path |
| --- | --- | --- |
| `Get` | GET | `/risk` |

## Portfolio — `client.Portfolio`

| Method | HTTP | Path |
| --- | --- | --- |
| `NewMovement` | POST | `/portfolio/movement` |

## Delivery notes — `client.DeliveryNotes`

| Method | HTTP | Path |
| --- | --- | --- |
| `Create` | POST | `/deliverynote` |

## QuantumBI (advisor) — `client.QuantumBI`

| Method | HTTP | Path |
| --- | --- | --- |
| `TotalHonorarium` | GET | `/advisorQuantumBi/totalHonorarium` |
| `SectionHonorary` | GET | `/advisorQuantumBi/sectionHonorary` |
| `ClientOGHonorary` | GET | `/advisorQuantumBi/clientOGHonorary` |
| `UserHonorariums` | GET | `/advisorQuantumBi/userHonorariums` |
| `UserClientHonorariums` | GET | `/advisorQuantumBi/userClientHonorariums` |
| `UserSectionHonorariums` | GET | `/advisorQuantumBi/userSectionHonorariums` |
| `UserTimeCost` | GET | `/advisorQuantumBi/userTimeCost` |
| `RateHonorary` | GET | `/advisorQuantumBi/rateHonorary` |
