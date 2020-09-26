package pdu

import (
	"fmt"
	"io/ioutil"
	"time"

	smpptime "github.com/daominah/smpp/time"
)

// DeliverSm contains mandatory fields for delivering short message.
// There is no need to set SmLength it will be automatically set when
// encoding pdu to binary representation.
type DeliverSm struct {
	ServiceType          string
	SourceAddrTon        int
	SourceAddrNpi        int
	SourceAddr           string
	DestAddrTon          int
	DestAddrNpi          int
	DestinationAddr      string
	EsmClass             EsmClass
	ProtocolID           int
	PriorityFlag         int
	ScheduleDeliveryTime time.Time
	ValidityPeriod       time.Time
	RegisteredDelivery   RegisteredDelivery
	ReplaceIfPresentFlag int
	DataCoding           int
	SmDefaultMsgID       int
	ShortMessage         string
	Options              *Options
}

// CommandID implements pdu.PDU interface.
func (p DeliverSm) CommandID() CommandID {
	return DeliverSmID
}

// Response creates new DeliverSmResp.
func (p DeliverSm) Response(msgID string) *DeliverSmResp {
	return &DeliverSmResp{
		MessageID: msgID,
	}
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p DeliverSm) MarshalBinary() ([]byte, error) {
	out := append(
		[]byte(p.ServiceType),
		0,
		byte(p.SourceAddrTon),
		byte(p.SourceAddrNpi),
	)
	out = append(out, append([]byte(p.SourceAddr), 0)...)
	out = append(out, byte(p.DestAddrTon), byte(p.DestAddrNpi))
	out = append(out, append([]byte(p.DestinationAddr), 0)...)
	out = append(out, p.EsmClass.Byte(), byte(p.ProtocolID), byte(p.PriorityFlag))
	tm, err := writeTime(smpptime.Absolute, p.ScheduleDeliveryTime)
	if err != nil {
		return nil, err
	}
	out = append(out, tm...)
	tm, err = writeTime(smpptime.Absolute, p.ValidityPeriod)
	if err != nil {
		return nil, err
	}
	out = append(out, tm...)
	l := len(p.ShortMessage)
	out = append(out, p.RegisteredDelivery.Byte(), byte(p.ReplaceIfPresentFlag), byte(p.DataCoding), byte(p.SmDefaultMsgID), byte(l))
	if l > 0 {
		out = append(out, []byte(p.ShortMessage)...)
	}
	if p.Options == nil {
		return out, nil
	}
	opts, err := p.Options.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return append(out, opts...), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p *DeliverSm) UnmarshalBinary(body []byte) error {
	if len(body) < 25 {
		return fmt.Errorf("smpp/pdu: deliver_sm body too short: %d", len(body))
	}
	buf := newBuffer(body)
	res, err := buf.ReadCString(6)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding service_type %s", err)
	}
	p.ServiceType = string(res)
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
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding dest_addr_ton %s", err)
	}
	p.DestAddrTon = int(b)
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding dest_addr_npi %s", err)
	}
	p.DestAddrNpi = int(b)
	res, err = buf.ReadCString(21)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding dest_addr %s", err)
	}
	p.DestinationAddr = string(res)
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding esm_class %s", err)
	}
	p.EsmClass = ParseEsmClass(b)
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding protocol_id %s", err)
	}
	p.ProtocolID = int(b)
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding priority_flag %s", err)
	}
	p.PriorityFlag = int(b)
	res, err = buf.ReadCString(17)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding schedule_delivery_time %s", err)
	}
	t, err := smpptime.Parse(res)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding schedule_delivery_time %s", err)
	}
	p.ScheduleDeliveryTime = t
	res, err = buf.ReadCString(17)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding validity_period %s", err)
	}
	t, err = smpptime.Parse(res)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding validity_period %s", err)
	}
	p.ValidityPeriod = t
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding registered_delivery %s", err)
	}
	p.RegisteredDelivery = ParseRegisteredDelivery(b)
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding replace_if_present_flag %s", err)
	}
	p.ReplaceIfPresentFlag = int(b)
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding data_coding %s", err)
	}
	p.DataCoding = int(b)
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding sm_default_msg_id %s", err)
	}
	p.SmDefaultMsgID = int(b)
	sm, err := buf.ReadString(254)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding short_message %s", err)
	}
	p.ShortMessage = string(sm)
	if buf.Len() == 0 {
		return nil
	}
	if p.Options == nil {
		p.Options = NewOptions()
	}
	opts, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	if err := p.Options.UnmarshalBinary(opts); err != nil {
		return err
	}
	return nil
}

// DeliverSmResp contains mandatory fields for deliver_sm response.
type DeliverSmResp struct {
	MessageID string
}

// CommandID implements pdu.PDU interface.
func (p DeliverSmResp) CommandID() CommandID {
	return DeliverSmRespID
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p DeliverSmResp) MarshalBinary() ([]byte, error) {
	return []byte{0}, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p *DeliverSmResp) UnmarshalBinary(body []byte) error {
	return nil
}
