package pdu

import (
	"encoding/binary"
	"fmt"
)

// Options maps all optional values and provides simple API for access.
// Only comonly used parameters have helpers, others have to be created
// by the users.
type Options struct {
	fields map[TagID][]byte
}

// NewOptions creates new options map.
func NewOptions() *Options {
	return &Options{
		fields: make(map[TagID][]byte),
	}
}

// Set assigns new TLV field.
func (o *Options) Set(tag TagID, val []byte) *Options {
	o.fields[tag] = val
	return o
}

// SetSingle assigns new TLV field with one byte value.
func (o *Options) SetSingle(tag TagID, val int) *Options {
	o.fields[tag] = []byte{byte(val)}
	return o
}

// SetDouble assigns new TLV field with two bytes value.
func (o *Options) SetDouble(tag TagID, val int) *Options {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(val))
	o.fields[tag] = b
	return o
}

// SetString assigns new TLV field with string value.
func (o *Options) SetString(tag TagID, val string) *Options {
	o.fields[tag] = []byte(val)
	return o
}

// SetCString assigns new TLV field with string value.
func (o *Options) SetCString(tag TagID, val string) *Options {
	o.fields[tag] = append([]byte(val), 0)
	return o
}

// Get tries to get byte value out of TLV field if present. If it's not it
// returns ok as false.
func (o *Options) Get(tag TagID) ([]byte, bool) {
	val, ok := o.fields[tag]
	return val, ok
}

// GetSingle returns tag value as one byte integer.
func (o *Options) GetSingle(tag TagID) (int, bool) {
	val, ok := o.fields[tag]
	if !ok {
		return 0, false
	}
	return int(val[0]), true
}

// GetDouble returns tag value as two byte integer.
func (o *Options) GetDouble(tag TagID) (int, bool) {
	b, ok := o.fields[tag]
	if !ok {
		return 0, false
	}
	return int(binary.BigEndian.Uint16(b)), true
}

// GetString returns tag value as string.
func (o *Options) GetString(tag TagID) (string, bool) {
	b, ok := o.fields[tag]
	if !ok {
		return "", false
	}
	return string(b), true
}

// GetCString returns tag value as string.
func (o *Options) GetCString(tag TagID) (string, bool) {
	b, ok := o.fields[tag]
	if !ok || len(b) == 0 {
		return "", false
	}
	return string(b[:len(b)-1]), true
}

// UserMessageReference is helper function for getting this option.
func (o *Options) UserMessageReference() int {
	val, ok := o.GetDouble(TagUserMessageReference)
	if !ok {
		return 0
	}
	return val
}

// SarMsgRefNum is helper function for getting this option.
func (o *Options) SarMsgRefNum() int {
	val, ok := o.GetDouble(TagSarMsgRefNum)
	if !ok {
		return 0
	}
	return val
}

// SarTotalSegments is helper function for getting this option.
func (o *Options) SarTotalSegments() int {
	val, ok := o.GetSingle(TagSarTotalSegments)
	if !ok {
		return 0
	}
	return val
}

// SarSegmentSeqnum is helper function for getting this option.
func (o *Options) SarSegmentSeqnum() int {
	val, ok := o.GetSingle(TagSarSegmentSeqnum)
	if !ok {
		return 0
	}
	return val
}

// ScInterfaceVersion is helper function for getting this option.
func (o *Options) ScInterfaceVersion() int {
	val, ok := o.GetSingle(TagScInterfaceVersion)
	if !ok {
		return 0
	}
	return val
}

// MessagePayload is helper function for getting this option.
func (o *Options) MessagePayload() string {
	val, ok := o.GetString(TagMessagePayload)
	if !ok {
		return ""
	}
	return val
}

// MessageState is helper function for getting this option.
func (o *Options) MessageState() int {
	val, ok := o.GetSingle(TagMessageState)
	if !ok {
		return 0
	}
	return val
}

// ReceiptedMessageID is helper function for getting this option.
func (o *Options) ReceiptedMessageID() string {
	val, ok := o.GetCString(TagReceiptedMessageID)
	if !ok {
		return ""
	}
	return val
}

// SetUserMessageReference is helper function for setting this option.
func (o *Options) SetUserMessageReference(val int) *Options {
	return o.SetDouble(TagUserMessageReference, val)
}

// SetSarMsgRefNum is helper function for setting this option.
func (o *Options) SetSarMsgRefNum(val int) *Options {
	return o.SetDouble(TagSarMsgRefNum, val)
}

// SetSarTotalSegments is helper function for setting this option.
func (o *Options) SetSarTotalSegments(val int) *Options {
	return o.SetSingle(TagSarTotalSegments, val)
}

// SetSarSegmentSeqnum is helper function for setting this option.
func (o *Options) SetSarSegmentSeqnum(val int) *Options {
	return o.SetSingle(TagSarSegmentSeqnum, val)
}

// SetScInterfaceVersion is helper function for setting this option.
func (o *Options) SetScInterfaceVersion(val int) *Options {
	return o.SetSingle(TagScInterfaceVersion, val)
}

// SetMessagePayload is helper function for setting this option.
func (o *Options) SetMessagePayload(val string) *Options {
	return o.SetString(TagMessagePayload, val)
}

// SetMessageState is helper function for setting this option.
func (o *Options) SetMessageState(val int) *Options {
	return o.SetSingle(TagMessageState, val)
}

// SetReceiptedMessageID is helper function for setting this option.
func (o *Options) SetReceiptedMessageID(val string) *Options {
	return o.SetCString(TagReceiptedMessageID, val)
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (o *Options) MarshalBinary() ([]byte, error) {
	var out []byte
	for tag, val := range o.fields {
		tlv := make([]byte, 4+len(val))
		binary.BigEndian.PutUint16(tlv[:2], uint16(tag))
		binary.BigEndian.PutUint16(tlv[2:4], uint16(len(val)))
		copy(tlv[4:], val)
		out = append(out, tlv...)
	}
	return out, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (o *Options) UnmarshalBinary(buf []byte) error {
	n := 0
	for n < len(buf) {
		if len(buf)-n <= 4 {
			return fmt.Errorf("smpp/pdu: invalid optional body length")
		}
		tag := TagID(binary.BigEndian.Uint16(buf[n : n+2]))
		l := int(binary.BigEndian.Uint16(buf[n+2 : n+4]))
		if n+4+l >= len(buf)+1 {
			return fmt.Errorf("smpp/pdu: invalid optional field length (%s %d)", tag, l)
		}
		o.fields[tag] = buf[n+4 : n+4+l]
		n += 4 + l
	}
	return nil
}
