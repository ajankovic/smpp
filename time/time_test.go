package time_test

import (
	"testing"
	gotime "time"

	"github.com/ajankovic/smpp/time"
)

func TestParseRelative(t *testing.T) {
	in := []byte("020610233429000R")
	future := gotime.Now().UTC().AddDate(2, 6, 12)
	past := gotime.Now().UTC().AddDate(2, 6, 9)
	out, err := time.Parse(in)
	if err != nil {
		t.Error(err)
	}
	if !out.Before(future) {
		t.Errorf("parsed time %s is not before expected %s", out, future)
	}
	if !out.After(past) {
		t.Errorf("parsed time %s is not after expected %s", out, past)
	}
}

func TestParseAbsolute(t *testing.T) {
	in := []byte("020610233429120-")
	loc := gotime.FixedZone("Custom", -5*3600)
	expected := gotime.Date(2002, gotime.June, 10, 23, 34, 29, 100000000, loc)
	out, err := time.Parse(in)
	if err != nil {
		t.Error(err)
	}
	if !out.Equal(expected) {
		t.Errorf("time not expected %s", out)
	}
}

func TestParseSimpleMinutes(t *testing.T) {
	in := []byte("0206102334")
	expected := gotime.Date(2002, gotime.June, 10, 23, 34, 0, 0, gotime.UTC)
	out, err := time.Parse(in)
	if err != nil {
		t.Error(err)
	}
	if !out.Equal(expected) {
		t.Errorf("time not expected %s", out)
	}
}

func TestParseSimpleSecs(t *testing.T) {
	in := []byte("020610233413")
	expected := gotime.Date(2002, gotime.June, 10, 23, 34, 13, 0, gotime.UTC)
	out, err := time.Parse(in)
	if err != nil {
		t.Error(err)
	}
	if !out.Equal(expected) {
		t.Errorf("time not expected %s", out)
	}
}

func TestParseInvalidFormat(t *testing.T) {
	in := []byte("invalidformat")
	_, err := time.Parse(in)
	if err == nil {
		t.Error("expected error got nil")
	}
	in = []byte("invalid")
	_, err = time.Parse(in)
	if err == nil {
		t.Error("expected error got nil")
	}
}

func TestFormatSecs(t *testing.T) {
	d := gotime.Date(2002, gotime.June, 10, 23, 34, 13, 0, gotime.UTC)
	expected := "020610233413"
	out, err := time.Format(time.SimpleSeconds, d)
	if err != nil {
		t.Error(err)
	}
	if out != expected {
		t.Errorf("format not expected %s", out)
	}
}

func TestFormatMins(t *testing.T) {
	d := gotime.Date(2002, gotime.June, 10, 23, 34, 0, 0, gotime.UTC)
	expected := "0206102334"
	out, err := time.Format(time.SimpleMinutes, d)
	if err != nil {
		t.Error(err)
	}
	if out != expected {
		t.Errorf("format not expected %s", out)
	}
}

func TestFormatAbsolute(t *testing.T) {
	d := gotime.Date(2002, gotime.June, 10, 23, 34, 13, 100000000, gotime.UTC)
	expected := "020610233413100+"
	out, err := time.Format(time.Absolute, d)
	if err != nil {
		t.Error(err)
	}
	if out != expected {
		t.Errorf("format not expected %s", out)
	}
}

func TestFormatRelative(t *testing.T) {
	d := gotime.Now().UTC().Add(10 * gotime.Hour)
	expected := "000000100000000R"
	out, err := time.Format(time.Relative, d)
	if err != nil {
		t.Error(err)
	}
	if out != expected {
		t.Errorf("format not expected %s", out)
	}
}
