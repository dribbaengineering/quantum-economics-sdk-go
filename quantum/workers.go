package quantum

import (
	"context"
	"net/http"
	"net/url"
)

// Absence is a worker absence record, used both as a request body (send/delete)
// and as a response item.
type Absence struct {
	ID             int64   `json:"id,omitempty" xml:"id,omitempty"`
	RegidWorker    int64   `json:"regid_worker,omitempty" xml:"regid_worker,omitempty"`
	WorkerName     string  `json:"worker_name,omitempty" xml:"worker_name,omitempty"`
	DateStart      string  `json:"date_start,omitempty" xml:"date_start,omitempty"`
	DateEnd        string  `json:"date_end,omitempty" xml:"date_end,omitempty"`
	Duration       float64 `json:"duration,omitempty" xml:"duration,omitempty"`
	Reason         int     `json:"reason,omitempty" xml:"reason,omitempty"`
	ReasonName     string  `json:"reason_name,omitempty" xml:"reason_name,omitempty"`
	Description    string  `json:"description,omitempty" xml:"description,omitempty"`
	Validated      int     `json:"validated,omitempty" xml:"validated,omitempty"`
	Year           int     `json:"year,omitempty" xml:"year,omitempty"`
	ValidationDesc string  `json:"validation_desc,omitempty" xml:"validation_desc,omitempty"`
	ValidationDate string  `json:"validation_date,omitempty" xml:"validation_date,omitempty"`
	ValidationUser string  `json:"validation_user,omitempty" xml:"validation_user,omitempty"`
	// File is the attachment contents, Base64-encoded.
	File     string `json:"file,omitempty" xml:"file,omitempty"`
	Filename string `json:"filename,omitempty" xml:"filename,omitempty"`
}

// WorkTimeRegistry is a work-time clock record, used both as a request body
// (send/edit) and as a response item.
type WorkTimeRegistry struct {
	Regid           int64   `json:"regid,omitempty" xml:"regid,omitempty"`
	RegidWorker     int64   `json:"regid_worker,omitempty" xml:"regid_worker,omitempty"`
	RegTime         string  `json:"reg_time,omitempty" xml:"reg_time,omitempty"`
	SentTime        string  `json:"sent_time,omitempty" xml:"sent_time,omitempty"`
	Type            int     `json:"type,omitempty" xml:"type,omitempty"`
	Modified        int     `json:"modified,omitempty" xml:"modified,omitempty"`
	Validated       int     `json:"validated,omitempty" xml:"validated,omitempty"`
	ValidatedWorker int     `json:"validated_worker,omitempty" xml:"validated_worker,omitempty"`
	Latitude        float64 `json:"latitude,omitempty" xml:"latitude,omitempty"`
	Longitude       float64 `json:"longitude,omitempty" xml:"longitude,omitempty"`
	Telematic       bool    `json:"telematic,omitempty" xml:"telematic,omitempty"`
	Preregistered   bool    `json:"preregistered,omitempty" xml:"preregistered,omitempty"`
}

// AbsencesListParams are the filters for listing worker absences. Year is
// required; Worker and Page are optional.
type AbsencesListParams struct {
	Year      int
	Worker    string
	Page      int
	CompanyID int64
}

// WorkingDaysParams are the parameters for calculating working days between two
// dates. DateStart, DateEnd and Year are required; Worker is optional.
type WorkingDaysParams struct {
	Worker    string
	DateStart string
	DateEnd   string
	Year      int
	CompanyID int64
}

// CalendarParams are the filters for the per-worker calendar. Year is required;
// the rest are optional.
type CalendarParams struct {
	Year           int
	Worker         string
	DateStart      string
	DateEnd        string
	ReturnWorkTime *bool
	CompanyID      int64
}

// CalendarAllParams are the filters for the all-workers calendar. Year is
// required; the rest are optional.
type CalendarAllParams struct {
	Year           int
	Incidents      *bool
	DateStart      string
	DateEnd        string
	ReturnWorkTime *bool
	CompanyID      int64
}

// WorkTimeListParams are the filters for listing work-time records. Year is
// required; the rest are optional.
type WorkTimeListParams struct {
	Year      int
	DateStart string
	DateEnd   string
	Page      int
	Worker    string
	CompanyID int64
}

// ValidationParams are the parameters for the absence/work-time validation
// endpoints. Validation, Year, Worker and Regids are required; Reason is
// optional. Regids is a comma-separated list of record ids.
type ValidationParams struct {
	Validation int
	Reason     string
	Year       int
	Worker     string
	Regids     string
	CompanyID  int64
}

// WorkersService groups the worker time-tracking and absence operations. Most
// read endpoints are only loosely specified upstream and therefore return
// RawResponse; decode the payload into []Absence, []WorkTimeRegistry or your
// own type with RawResponse.Decode.
type WorkersService struct {
	client *Client
}

// --- Absences ---

// Absences lists a worker's absences for a year (GET /worker/absences).
func (s *WorkersService) Absences(ctx context.Context, params AbsencesListParams) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setStringOpt("worker", params.Worker).
		setIntOpt("page", int64(params.Page)).
		values()
	return s.raw(ctx, http.MethodGet, "/worker/absences", q, nil)
}

// AbsencesSummary returns a summary of all workers' absences for a year
// (GET /worker/absences/summary).
func (s *WorkersService) AbsencesSummary(ctx context.Context, year int) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).setInt("year", int64(year)).values()
	return s.raw(ctx, http.MethodGet, "/worker/absences/summary", q, nil)
}

// AbsenceFile returns the Base64 file of an absence (GET /worker/absences/file).
func (s *WorkersService) AbsenceFile(ctx context.Context, hash string) (*DocumentResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).setString("hash", hash).values()
	out := &DocumentResponse{}
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/worker/absences/file", query: q}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// WorkingDays calculates the working days between two dates
// (GET /worker/absences/getWorkingDays).
func (s *WorkersService) WorkingDays(ctx context.Context, params WorkingDaysParams) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setStringOpt("worker", params.Worker).
		setString("dateStart", params.DateStart).
		setString("dateEnd", params.DateEnd).
		setInt("year", int64(params.Year)).
		values()
	return s.raw(ctx, http.MethodGet, "/worker/absences/getWorkingDays", q, nil)
}

// SendAbsence submits a new absence request (POST /worker/absences/send).
func (s *WorkersService) SendAbsence(ctx context.Context, worker string, year int, absence Absence) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setIntOpt("year", int64(year)).
		setStringOpt("worker", worker).
		values()
	return s.raw(ctx, http.MethodPost, "/worker/absences/send", q, absence)
}

// DeleteAbsence deletes an absence from a worker (POST /worker/absences/delete).
func (s *WorkersService) DeleteAbsence(ctx context.Context, worker string, year int, absence Absence) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(year)).
		setString("worker", worker).
		values()
	return s.raw(ctx, http.MethodPost, "/worker/absences/delete", q, absence)
}

// ValidateAbsences validates a list of a worker's absences
// (POST /worker/absences/validate).
func (s *WorkersService) ValidateAbsences(ctx context.Context, params ValidationParams) (*RawResponse, error) {
	return s.validate(ctx, "/worker/absences/validate", params)
}

// --- Calendar ---

// Calendar returns a worker's work-time and absence calendar events
// (GET /worker/calendar).
func (s *WorkersService) Calendar(ctx context.Context, params CalendarParams) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setStringOpt("worker", params.Worker).
		setStringOpt("dateStart", params.DateStart).
		setStringOpt("dateEnd", params.DateEnd).
		setBoolOpt("returnWorkTime", params.ReturnWorkTime).
		values()
	return s.raw(ctx, http.MethodGet, "/worker/calendar", q, nil)
}

// CalendarAll returns the calendar events of all workers
// (GET /worker/calendar/all).
func (s *WorkersService) CalendarAll(ctx context.Context, params CalendarAllParams) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setBoolOpt("incidents", params.Incidents).
		setStringOpt("dateStart", params.DateStart).
		setStringOpt("dateEnd", params.DateEnd).
		setBoolOpt("returnWorkTime", params.ReturnWorkTime).
		values()
	return s.raw(ctx, http.MethodGet, "/worker/calendar/all", q, nil)
}

// CalendarIncidents returns the calendar incidents of all workers
// (GET /worker/calendar/incidents).
func (s *WorkersService) CalendarIncidents(ctx context.Context, params CalendarAllParams) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setStringOpt("dateStart", params.DateStart).
		setStringOpt("dateEnd", params.DateEnd).
		setBoolOpt("returnWorkTime", params.ReturnWorkTime).
		values()
	return s.raw(ctx, http.MethodGet, "/worker/calendar/incidents", q, nil)
}

// --- Work time ---

// WorkTime lists work-time records for a year (GET /worker/workTime).
func (s *WorkersService) WorkTime(ctx context.Context, params WorkTimeListParams) (*RawResponse, error) {
	return s.workTimeList(ctx, "/worker/workTime", params)
}

// WorkTimeWeek lists work-time records for a week (GET /worker/workTime/week).
// DateStart is required.
func (s *WorkersService) WorkTimeWeek(ctx context.Context, year int, dateStart, worker string) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(year)).
		setString("dateStart", dateStart).
		setStringOpt("worker", worker).
		values()
	return s.raw(ctx, http.MethodGet, "/worker/workTime/week", q, nil)
}

// WorkTimePending lists work-time records pending validation
// (GET /worker/workTime/pending/worker).
func (s *WorkersService) WorkTimePending(ctx context.Context, params WorkTimeListParams) (*RawResponse, error) {
	return s.workTimeList(ctx, "/worker/workTime/pending/worker", params)
}

// WorkTimePre lists the missing work-time records inferred from schedule and
// records (GET /worker/workTime/pre).
func (s *WorkersService) WorkTimePre(ctx context.Context, params WorkTimeListParams) (*RawResponse, error) {
	return s.workTimeList(ctx, "/worker/workTime/pre", params)
}

// WorkTimeStatus returns a worker's work-time status
// (GET /worker/workTime/status).
func (s *WorkersService) WorkTimeStatus(ctx context.Context, year int, maxDate, maxTime string) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(year)).
		setStringOpt("maxDate", maxDate).
		setStringOpt("maxTime", maxTime).
		values()
	return s.raw(ctx, http.MethodGet, "/worker/workTime/status", q, nil)
}

// WorkTimeSummary returns a weekly summary of all workers' records
// (GET /worker/workTime/summary). DateStart is required.
func (s *WorkersService) WorkTimeSummary(ctx context.Context, year int, dateStart, dateEnd string) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(year)).
		setString("dateStart", dateStart).
		setStringOpt("dateEnd", dateEnd).
		values()
	return s.raw(ctx, http.MethodGet, "/worker/workTime/summary", q, nil)
}

// SendWorkTime registers a new work-time record for a worker
// (POST /worker/workTime/send).
func (s *WorkersService) SendWorkTime(ctx context.Context, record WorkTimeRegistry) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().setInt("companyId", companyID).values()
	return s.raw(ctx, http.MethodPost, "/worker/workTime/send", q, record)
}

// EditWorkTime edits a work-time record of a worker
// (POST /worker/workTime/edit).
func (s *WorkersService) EditWorkTime(ctx context.Context, id int64, year int, record WorkTimeRegistry) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("id", id).
		setInt("year", int64(year)).
		values()
	return s.raw(ctx, http.MethodPost, "/worker/workTime/edit", q, record)
}

// WorkTimeReminder sends a clock-in reminder push notification to a worker
// (POST /worker/workTime/reminder).
func (s *WorkersService) WorkTimeReminder(ctx context.Context, year int, worker string) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(0)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(year)).
		setString("worker", worker).
		values()
	return s.raw(ctx, http.MethodPost, "/worker/workTime/reminder", q, nil)
}

// ValidateWorkTime validates a list of a worker's work-time records
// (POST /worker/workTime/validate).
func (s *WorkersService) ValidateWorkTime(ctx context.Context, params ValidationParams) (*RawResponse, error) {
	return s.validate(ctx, "/worker/workTime/validate", params)
}

// ValidateWorkTimeByWorker validates a worker's own work-time records
// (POST /worker/workTime/validate/worker).
func (s *WorkersService) ValidateWorkTimeByWorker(ctx context.Context, params ValidationParams) (*RawResponse, error) {
	return s.validate(ctx, "/worker/workTime/validate/worker", params)
}

// ValidateAllWorkTime validates work-time records across a date range
// (POST /worker/workTime/validateAll). Validation and Year are required;
// DateStart/DateEnd/Reason are optional.
func (s *WorkersService) ValidateAllWorkTime(ctx context.Context, params ValidationParams, dateStart, dateEnd string) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("validation", int64(params.Validation)).
		setStringOpt("reason", params.Reason).
		setInt("year", int64(params.Year)).
		setStringOpt("dateStart", dateStart).
		setStringOpt("dateEnd", dateEnd).
		values()
	return s.raw(ctx, http.MethodPost, "/worker/workTime/validateAll", q, nil)
}

// --- helpers ---

func (s *WorkersService) workTimeList(ctx context.Context, path string, params WorkTimeListParams) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("year", int64(params.Year)).
		setStringOpt("dateStart", params.DateStart).
		setStringOpt("dateEnd", params.DateEnd).
		setIntOpt("page", int64(params.Page)).
		setStringOpt("worker", params.Worker).
		values()
	return s.raw(ctx, http.MethodGet, path, q, nil)
}

func (s *WorkersService) validate(ctx context.Context, path string, params ValidationParams) (*RawResponse, error) {
	companyID, err := s.client.resolveCompanyID(params.CompanyID)
	if err != nil {
		return nil, err
	}
	q := newQuery().
		setInt("companyId", companyID).
		setInt("validation", int64(params.Validation)).
		setStringOpt("reason", params.Reason).
		setInt("year", int64(params.Year)).
		setString("worker", params.Worker).
		setString("regids", params.Regids).
		values()
	return s.raw(ctx, http.MethodPost, path, q, nil)
}

func (s *WorkersService) raw(ctx context.Context, method, path string, q url.Values, body any) (*RawResponse, error) {
	out := &RawResponse{}
	if err := s.client.do(ctx, request{method: method, path: path, query: q, body: body}, out); err != nil {
		return nil, err
	}
	return out, nil
}
