package pdu

import (
	"fmt"
	"time"

	smpptime "github.com/ajankovic/smpp/time"
)

// QuerySm represents quering PDU.
type QuerySm struct {
	MessageID     string
	SourceAddrTon int
	SourceAddrNpi int
	SourceAddr    string
}

// CommandID implements pdu.PDU interface.
func (p QuerySm) CommandID() CommandID {
	return QuerySmID
}

// Response creates new QuerySmResp.
func (p QuerySm) Response(date time.Time, state, err int) *QuerySmResp {
	return &QuerySmResp{
		MessageID:    p.MessageID,
		FinalDate:    date,
		MessageState: state,
		ErrorCode:    err,
	}
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p QuerySm) MarshalBinary() ([]byte, error) {
	out := append([]byte(p.MessageID), 0)
	out = append(out, byte(p.SourceAddrTon), byte(p.SourceAddrNpi))
	out = append(out, append([]byte(p.SourceAddr), 0)...)
	return out, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p *QuerySm) UnmarshalBinary(body []byte) error {
	if len(body) < 6 {
		return fmt.Errorf("smpp/pdu: query_sm body too short: %d", len(body))
	}
	buf := newBuffer(body)
	res, err := buf.ReadCString(65)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding message_id %s", err)
	}
	p.MessageID = string(res)
	b, err := buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding source_addr_ton %s", err)
	}
	p.SourceAddrTon = int(b)
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding source_addr_npi %s", err)
	}
	p.SourceAddrNpi = int(b)
	res, err = buf.ReadCString(21)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding source_addr %s", err)
	}
	p.SourceAddr = string(res)
	return nil
}

// QuerySmResp holds response to query_sm PDU.
type QuerySmResp struct {
	MessageID    string
	FinalDate    time.Time
	MessageState int
	ErrorCode    int
}

// CommandID implements pdu.PDU interface.
func (p QuerySmResp) CommandID() CommandID {
	return QuerySmRespID
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p QuerySmResp) MarshalBinary() ([]byte, error) {
	out := append([]byte(p.MessageID), 0)
	tm, err := writeTime(smpptime.Absolute, p.FinalDate)
	if err != nil {
		return nil, err
	}
	out = append(out, tm...)
	out = append(out, byte(p.MessageState), byte(p.ErrorCode))
	return out, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p *QuerySmResp) UnmarshalBinary(body []byte) error {
	if len(body) < 6 {
		return fmt.Errorf("smpp/pdu: query_sm body too short: %d", len(body))
	}
	buf := newBuffer(body)
	res, err := buf.ReadCString(65)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding message_id %s", err)
	}
	p.MessageID = string(res)
	res, err = buf.ReadCString(17)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding final_date %s", err)
	}
	t, err := smpptime.Parse(res)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding final_date %s", err)
	}
	p.FinalDate = t
	b, err := buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding message_state %s", err)
	}
	p.MessageState = int(b)
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding error_code %s", err)
	}
	p.ErrorCode = int(b)
	return nil
}
