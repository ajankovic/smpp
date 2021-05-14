package time

import (
	"errors"
	"fmt"
	"time"
	gotime "time"
)

// Layout defines SMPP time layout in string representation.
// It can be Relative, Absolute, Simple.
type Layout int

const (
	// SimpleSeconds layout in seconds YYMMDDhhmmss.
	SimpleSeconds Layout = iota
	// SimpleMinutes layout in minutes YYMMDDhhmm.
	SimpleMinutes
	// Absolute layout YYMMDDhhmmsstnn[+-].
	Absolute
	// Relative layout YYMMDDhhmmss000[R].
	Relative
)

// Parse converts bytestring representation of time from SMPP format
// to standard time.Time. Relative layouts will be added to the current
// time and returned as time.Time.
func Parse(in []byte) (gotime.Time, error) {
	l := len(in)
	switch l {
	case 1, 0:
		// nil time
		return gotime.Time{}, nil
	case 12:
		// simple seconds
		return gotime.Parse("060102150405", string(in))
	case 14:
		return gotime.Parse("20060102150405", string(in))
	case 10:
		// simple minutes
		return gotime.Parse("0601021504", string(in))
	case 16:
		layoutIndicator := in[l-1]
		switch layoutIndicator {
		case 'R':
			// Relative layout.
			y := int((in[0]-48)*10 + (in[1] - 48))
			mo := int((in[2]-48)*10 + (in[3] - 48))
			d := int((in[4]-48)*10 + (in[5] - 48))
			h := int((in[6]-48)*10 + (in[7] - 48))
			mi := int((in[8]-48)*10 + (in[9] - 48))
			s := int((in[10]-48)*10 + (in[11] - 48))
			return gotime.Now().
				AddDate(y, mo, d).
				Add(time.Duration(h)*time.Hour +
					time.Duration(mi)*time.Minute +
					time.Duration(s)*time.Second), nil
		case '-', '+':
			// Absolute layout.
			nn := int((in[13]-48)*10 + (in[14] - 48))
			offset := nn * 900 // 15 min intervals in seconds.
			if layoutIndicator == '-' {
				offset = -offset
			}
			var loc *gotime.Location
			if offset != 0 {
				loc = gotime.FixedZone("Custom", offset)
			} else {
				loc = gotime.UTC
			}
			t, err := gotime.ParseInLocation("060102150405", string(in[:l-4]), loc)
			if err != nil {
				return time.Time{}, err
			}
			t = t.Add(time.Duration(in[12]-48) * 100 * time.Millisecond)
			return t, nil
		default:
			return gotime.Time{}, fmt.Errorf("smpp/time: invalid layout length %s", in)
		}
	default:
		return gotime.Time{}, fmt.Errorf("smpp/time: invalid layout length %s", in)
	}
}

// Format converts time.Time into string representation defined by smpp
// predefined layout.
func Format(layout Layout, t gotime.Time) (string, error) {
	switch layout {
	case SimpleSeconds:
		return t.Format("060102150405"), nil
	case SimpleMinutes:
		return t.Format("0601021504"), nil
	case Relative:
		y, mo, d, h, mi, s := diff(t, gotime.Now())
		return fmt.Sprintf("%02d%02d%02d%02d%02d%02d000R", y, mo, d, h, mi, s), nil
	case Absolute:
		sign := "+"
		_, z := t.Zone()
		offset := z / 900
		if offset < 0 {
			sign = "-"
			offset = -offset
		}
		return fmt.Sprintf("%s%d%02d%s", t.Format("060102150405"), t.Nanosecond()/100000000, offset, sign), nil
	default:
		return "", errors.New("smpp/time: invalid format layout")
	}
}

// Go supports only dif with hours so borrowing this from
// https://stackoverflow.com/questions/36530251/golang-time-since-with-months-and-years
func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// Days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
