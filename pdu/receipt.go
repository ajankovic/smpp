package pdu

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// DeliveryReceipt in format
// “id:IIIIIIIIII sub:SSS dlvrd:DDD submit date:YYMMDDhhmm done date:YYMMDDhhmm stat:DDDDDDD err:E Text: ...”
type DeliveryReceipt struct {
	Id         string
	Sub        string
	Dlvrd      string
	SubmitDate time.Time
	DoneDate   time.Time
	Stat       DelStat
	Err        string
	Text       string
}

type DelStat string

const (
	DelStatEnRoute       DelStat = "ENROUTE"
	DelStatDelivered     DelStat = "DELIVRD"
	DelStatExpired       DelStat = "EXPIRED"
	DelStatDeleted       DelStat = "DELETED"
	DelStatUndeliverable DelStat = "UNDELIV"
	DelStatAccepted      DelStat = "ACCEPTD"
	DelStatUnknown       DelStat = "UNKNOWN"
	DelStatRejected      DelStat = "REJECTD"
)

var DelStatMap = map[uint8]DelStat{
	1: DelStatEnRoute,
	2: DelStatDelivered,
	3: DelStatExpired,
	4: DelStatDeleted,
	5: DelStatUndeliverable,
	6: DelStatAccepted,
	7: DelStatUnknown,
	8: DelStatRejected,
}

func (dr *DeliveryReceipt) String() string {
	return fmt.Sprintf(
		"id:%s sub:%s dlvrd:%s submit date:%s done date:%s stat:%s err:%s text:%s",
		dr.Id, dr.Sub, dr.Dlvrd, dr.SubmitDate.Format(RecDateLayout), dr.DoneDate.Format(RecDateLayout), dr.Stat, dr.Err, dr.Text,
	)
}

var deliveryReceipt = regexp.MustCompile(`(\w+ ?\w+)+:([\w\-]+)`)

// YYMMDDhhmm
var RecDateLayout = "0601021504"
var SecRecDateLayout = "060102150405"

var dateFormats = []string{"20060102150405", "0601021504", "060102150405"}

func ParseDateTime(value string) (time.Time, error) {
	for _, df := range dateFormats {
		if result, err := time.ParseInLocation(value, df, time.Local); err == nil {
			return result, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse time %s", value)
}

// ParseDeliveryReceipt parses delivery receipt format defined in smpp 3.4 specification
func ParseDeliveryReceipt(sm string) (*DeliveryReceipt, error) {
	e := errors.New("smpp: invalid receipt format")
	i := strings.Index(sm, "text:")
	if i == -1 {
		i = strings.Index(sm, "Text:")
		if i == -1 {
			return nil, e
		}
	}
	delRec := DeliveryReceipt{}
	match := deliveryReceipt.FindAllStringSubmatch(sm[:i], -1)
	for idx, m := range match {
		if len(m) != 3 {
			return nil, e
		}
		// TODO improve error with more details
		switch idx {
		case 0:
			if m[1] != "id" {
				return nil, e
			}
			delRec.Id = m[2]
		case 1:
			if m[1] != "sub" {
				return nil, e
			}
			delRec.Sub = m[2]
		case 2:
			if m[1] != "dlvrd" {
				return nil, e
			}
			delRec.Dlvrd = m[2]
		case 3:
			if m[1] != "submit date" {
				return nil, e
			}
			t, err := ParseDateTime(m[2])
			if err != nil {
				return nil, e
			}
			delRec.SubmitDate = t
		case 4:
			if m[1] != "done date" {
				return nil, e
			}
			t, err := ParseDateTime(m[2])
			if err != nil {
				return nil, e
			}
			delRec.DoneDate = t
		case 5:
			if m[1] != "stat" {
				return nil, e
			}
			// TODO validate status value
			delRec.Stat = DelStat(m[2])
		case 6:
			if m[1] != "err" {
				return nil, e
			}
			delRec.Err = m[2]
		default:
			return nil, e
		}
	}
	delRec.Text = sm[i+5:]
	return &delRec, nil
}
