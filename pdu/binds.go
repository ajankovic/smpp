package pdu

import (
	"fmt"
)

// BindTx binding pdu in transmitter mode.
type BindTx struct {
	SystemID         string
	Password         string
	SystemType       string
	InterfaceVersion int
	AddrTon          int
	AddrNpi          int
	AddressRange     string
}

// CommandID implements pdu.PDU interface.
func (p BindTx) CommandID() CommandID {
	return BindTransmitterID
}

// Response creates new BindTxResp.
func (p BindTx) Response(sysID string) *BindTxResp {
	return &BindTxResp{
		SystemID: sysID,
	}
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p BindTx) MarshalBinary() ([]byte, error) {
	return marshalBind(
		p.SystemID,
		p.Password,
		p.SystemType,
		p.InterfaceVersion,
		p.AddrTon,
		p.AddrNpi,
		p.AddressRange,
	)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p *BindTx) UnmarshalBinary(body []byte) error {
	return unmarshalBind(
		body,
		&p.SystemID,
		&p.Password,
		&p.SystemType,
		&p.InterfaceVersion,
		&p.AddrTon,
		&p.AddrNpi,
		&p.AddressRange,
	)
}

// BindTxResp bind response.
type BindTxResp struct {
	SystemID string
	Options  *Options
}

// CommandID implements pdu.PDU interface.
func (p BindTxResp) CommandID() CommandID {
	return BindTransmitterRespID
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p BindTxResp) MarshalBinary() ([]byte, error) {
	return cStringOptsRespMarshal(p.SystemID, p.Options)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p *BindTxResp) UnmarshalBinary(body []byte) error {
	var err error
	p.SystemID, p.Options, err = cStringOptsRespUnmarshal(body)
	return err
}

// BindRx binding pdu in receiver mode.
type BindRx struct {
	SystemID         string
	Password         string
	SystemType       string
	InterfaceVersion int
	AddrTon          int
	AddrNpi          int
	AddressRange     string
}

// CommandID implements pdu.PDU interface.
func (p BindRx) CommandID() CommandID {
	return BindReceiverID
}

// Response creates new BindRxResp.
func (p BindRx) Response(sysID string) *BindRxResp {
	return &BindRxResp{
		SystemID: sysID,
	}
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p BindRx) MarshalBinary() ([]byte, error) {
	return marshalBind(
		p.SystemID,
		p.Password,
		p.SystemType,
		p.InterfaceVersion,
		p.AddrTon,
		p.AddrNpi,
		p.AddressRange,
	)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p *BindRx) UnmarshalBinary(body []byte) error {
	return unmarshalBind(
		body,
		&p.SystemID,
		&p.Password,
		&p.SystemType,
		&p.InterfaceVersion,
		&p.AddrTon,
		&p.AddrNpi,
		&p.AddressRange,
	)
}

// BindRxResp bind response.
type BindRxResp struct {
	SystemID string
	Options  *Options
}

// CommandID implements pdu.PDU interface.
func (p BindRxResp) CommandID() CommandID {
	return BindReceiverRespID
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p BindRxResp) MarshalBinary() ([]byte, error) {
	return cStringOptsRespMarshal(p.SystemID, p.Options)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p *BindRxResp) UnmarshalBinary(body []byte) error {
	var err error
	p.SystemID, p.Options, err = cStringOptsRespUnmarshal(body)
	return err
}

// BindTRx binding PDU in receiver mode.
type BindTRx struct {
	SystemID         string
	Password         string
	SystemType       string
	InterfaceVersion int
	AddrTon          int
	AddrNpi          int
	AddressRange     string
}

// CommandID implements pdu.PDU interface.
func (p BindTRx) CommandID() CommandID {
	return BindTransceiverID
}

// Response creates new BindTRxResp.
func (p BindTRx) Response(sysID string) *BindTRxResp {
	return &BindTRxResp{
		SystemID: sysID,
	}
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p BindTRx) MarshalBinary() ([]byte, error) {
	return marshalBind(
		p.SystemID,
		p.Password,
		p.SystemType,
		p.InterfaceVersion,
		p.AddrTon,
		p.AddrNpi,
		p.AddressRange,
	)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p *BindTRx) UnmarshalBinary(body []byte) error {
	return unmarshalBind(
		body,
		&p.SystemID,
		&p.Password,
		&p.SystemType,
		&p.InterfaceVersion,
		&p.AddrTon,
		&p.AddrNpi,
		&p.AddressRange,
	)
}

// BindTRxResp bind response.
type BindTRxResp struct {
	SystemID string
	Options  *Options
}

// CommandID implements pdu.PDU interface.
func (p BindTRxResp) CommandID() CommandID {
	return BindTransceiverRespID
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p BindTRxResp) MarshalBinary() ([]byte, error) {
	return cStringOptsRespMarshal(p.SystemID, p.Options)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p *BindTRxResp) UnmarshalBinary(body []byte) error {
	var err error
	p.SystemID, p.Options, err = cStringOptsRespUnmarshal(body)
	return err
}

func marshalBind(systemID, password, systemType string, interfaceVer, addrTon, addrNpi int, addrRange string) ([]byte, error) {
	out := append([]byte(systemID), 0)
	out = append(out, append([]byte(password), 0)...)
	out = append(out, append([]byte(systemType), 0)...)
	out = append(out, byte(interfaceVer), byte(addrTon), byte(addrNpi))
	out = append(out, append([]byte(addrRange), 0)...)
	return out, nil
}

func unmarshalBind(body []byte, systemID, password, systemType *string, interfaceVer, addrTon, addrNpi *int, addrRange *string) error {
	if len(body) < 7 {
		return fmt.Errorf("smpp/pdu: bind body too short: %d", len(body))
	}
	buf := newBuffer(body)
	res, err := buf.ReadCString(16)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding system_id %s", err)
	}
	*systemID = string(res)
	res, err = buf.ReadCString(9)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding password %s", err)
	}
	*password = string(res)
	res, err = buf.ReadCString(13)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding system_type %s", err)
	}
	*systemType = string(res)
	b, err := buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding interface_version %s", err)
	}
	*interfaceVer = int(b)
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding addr_ton %s", err)
	}
	*addrTon = int(b)
	b, err = buf.ReadByte()
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding addr_npi %s", err)
	}
	*addrNpi = int(b)
	res, err = buf.ReadCString(41)
	if err != nil {
		return fmt.Errorf("smpp/pdu: decoding addr_range %s", err)
	}
	*addrRange = string(res)
	return nil
}
