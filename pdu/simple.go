package pdu

// Unbind defines unbind PDU.
type Unbind struct{}

// CommandID implements pdu.PDU interface.
func (p Unbind) CommandID() CommandID {
	return UnbindID
}

// Response creates new UnbindResp.
func (p Unbind) Response() *UnbindResp {
	return &UnbindResp{}
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p Unbind) MarshalBinary() ([]byte, error) {
	return nil, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p Unbind) UnmarshalBinary(body []byte) error {
	return nil
}

// UnbindResp defines unbind_resp PDU.
type UnbindResp struct{}

// CommandID implements pdu.PDU interface.
func (p UnbindResp) CommandID() CommandID {
	return UnbindRespID
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p UnbindResp) MarshalBinary() ([]byte, error) {
	return nil, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p UnbindResp) UnmarshalBinary(body []byte) error {
	return nil
}

// EnquireLink PDU.
type EnquireLink struct{}

// CommandID implements pdu.PDU interface.
func (p EnquireLink) CommandID() CommandID {
	return EnquireLinkID
}

// Response creates new EnquireLinkResp.
func (p EnquireLink) Response() *EnquireLinkResp {
	return &EnquireLinkResp{}
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p EnquireLink) MarshalBinary() ([]byte, error) {
	return nil, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p EnquireLink) UnmarshalBinary(body []byte) error {
	return nil
}

// EnquireLinkResp PDU response.
type EnquireLinkResp struct{}

// CommandID implements pdu.PDU interface.
func (p EnquireLinkResp) CommandID() CommandID {
	return EnquireLinkRespID
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p EnquireLinkResp) MarshalBinary() ([]byte, error) {
	return nil, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p EnquireLinkResp) UnmarshalBinary(body []byte) error {
	return nil
}

// GenericNack PDU.
type GenericNack struct{}

// CommandID implements pdu.PDU interface.
func (p GenericNack) CommandID() CommandID {
	return GenericNackID
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (p GenericNack) MarshalBinary() ([]byte, error) {
	return nil, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (p GenericNack) UnmarshalBinary(body []byte) error {
	return nil
}
