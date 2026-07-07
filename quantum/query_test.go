package quantum

import (
	"testing"
	"time"
)

func TestQueryBuilderOmitsZeroValues(t *testing.T) {
	q := newQuery().
		setString("required", "x").
		setStringOpt("empty", "").
		setStringOpt("present", "y").
		setInt("num", 0).
		setIntOpt("zero", 0).
		setIntOpt("nonzero", 5).
		setBoolOpt("nilbool", nil).
		setBoolOpt("truebool", Bool(true)).
		values()

	if q.Get("required") != "x" || q.Get("present") != "y" {
		t.Errorf("required strings missing: %v", q)
	}
	if _, ok := q["empty"]; ok {
		t.Error("empty optional string should be omitted")
	}
	if q.Get("num") != "0" {
		t.Error("mandatory int should be written even when zero")
	}
	if _, ok := q["zero"]; ok {
		t.Error("zero optional int should be omitted")
	}
	if q.Get("nonzero") != "5" {
		t.Error("non-zero optional int should be written")
	}
	if _, ok := q["nilbool"]; ok {
		t.Error("nil *bool should be omitted")
	}
	if q.Get("truebool") != "true" {
		t.Error("non-nil *bool should be written")
	}
}

func TestDateHelpers(t *testing.T) {
	d := time.Date(2025, 9, 22, 0, 0, 0, 0, time.UTC)
	if got := FormatQueryDate(d); got != "22-09-2025" {
		t.Errorf("FormatQueryDate = %q", got)
	}
	if got := FormatBodyDate(d); got != "22/09/2025" {
		t.Errorf("FormatBodyDate = %q", got)
	}

	parsed, err := ParseBodyDate("22/09/2025")
	if err != nil {
		t.Fatalf("ParseBodyDate: %v", err)
	}
	if !parsed.Equal(d) {
		t.Errorf("ParseBodyDate = %v, want %v", parsed, d)
	}

	if _, err := ParseQueryDate("22-09-2025"); err != nil {
		t.Errorf("ParseQueryDate: %v", err)
	}
}
