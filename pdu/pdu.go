package pdu

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	smpptime "github.com/ajankovic/smpp/time"
)

// PDU defines interface for PDU structures
type PDU interface {
	CommandID() CommandID
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// EsmClass is used to indicate special message attributes associated with the short message.
type EsmClass struct {
	Mode    int
	Type    int
	Feature int
}

// Byte converts EsmClass into a single byte for pdu encoding.
func (ec EsmClass) Byte() byte {
	out := byte(0)
	out |= byte(ec.Mode)
	out |= byte(ec.Type) << 2
	out |= byte(ec.Feature) << 6
	return out
}

// ParseEsmClass parses esm class from pdu.
func ParseEsmClass(b byte) EsmClass {
	out := EsmClass{}
	out.Mode = int(b & 0x03)
	out.Type = int((b >> 2) & 0x0F)
	out.Feature = int(b >> 6)
	return out
}

const (
	DefaultEsmMode         = 0x0
	DatagramEsmMode        = 0x1
	ForwardEsmMode         = 0x2
	StoreAndForwardEsmMode = 0x3
	NotApplicableEsmMode   = 0x7
)

const (
	DefaultEsmType = 0x0
	DelRecEsmType  = 0x1
	DelAckEsmType  = 0x2
	UsrAckEsmType  = 0x4
	ConAbtEsmType  = 0x6
	IDNEsmType     = 0x8
)

const (
	NoEsmFeat          = 0x0
	UDHIEsmFeat        = 0x1
	RepPathEsmFeat     = 0x2
	UDHIRepPathEsmFeat = 0x3
)

// RegisteredDelivery is used to request an SMSC delivery receipt and/or SME
// originated acknowledgements.
type RegisteredDelivery struct {
	Receipt           int
	SMEAck            int
	InterNotification int
}

// Byte converts RegisteredDelivery into a single byte for pdu encoding.
func (rd RegisteredDelivery) Byte() byte {
	out := byte(0)
	out |= byte(rd.Receipt)
	out |= byte(rd.SMEAck) << 2
	out |= byte(rd.InterNotification) << 4
	return out
}

// ParseRegisteredDelivery parses registered_delivery from pdu.
func ParseRegisteredDelivery(b byte) RegisteredDelivery {
	out := RegisteredDelivery{}
	out.Receipt = int(b & 0x03)
	out.SMEAck = int((b >> 2) & 0x0F)
	out.InterNotification = int((b >> 4) & 0x01)
	return out
}

const (
	NoDeliveryReceipt   = 0x0
	YesDeliveryReceipt  = 0x1
	FailDeliveryReceipt = 0x2
)

const (
	NoSMEAck     = 0x0
	YesSMEAck    = 0x1
	ManualSMEAck = 0x2
	AllSMEAck    = 0x3
)

const (
	NoInterNotification  = 0x0
	YesInterNotification = 0x1
)

func writeTime(layout smpptime.Layout, t time.Time) ([]byte, error) {
	var schedDel []byte
	if !t.IsZero() {
		out, err := smpptime.Format(layout, t)
		if err != nil {
			return nil, err
		}
		schedDel = []byte(out)
	}
	return append(schedDel, 0), nil
}

type pduReader struct {
	*bytes.Buffer
}

func newBuffer(buf []byte) *pduReader {
	return &pduReader{
		Buffer: bytes.NewBuffer(buf),
	}
}

func (r *pduReader) ReadCString(limit int) ([]byte, error) {
	var out []byte
	i := 0
	for {
		i++
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == 0x0 {
			return out, nil
		}
		if i == limit {
			return nil, errors.New("invalid c string length")
		}
		out = append(out, b)
	}
}

func (r *pduReader) ReadString(limit int) ([]byte, error) {
	l, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if int(l) > limit {
		return nil, errors.New("invalid string length")
	}
	out := make([]byte, l)
	n, err := r.Read(out)
	if err != nil {
		return nil, err
	}
	if n != int(l) {
		return nil, errors.New("read count missmatch")
	}
	return out, nil
}

func cStringOptsRespUnmarshal(body []byte) (string, *Options, error) {
	n := -1
	for i := 0; i < len(body); i++ {
		if body[i] == 0 {
			n = i + 1
			break
		}
	}
	if n < 0 {
		return "", nil, errors.New("smpp/pdu: c string is not terminated")
	}
	var opts *Options
	if len(body[n:]) > 0 {
		opts = NewOptions()
		if err := opts.UnmarshalBinary(body[n:]); err != nil {
			return "", nil, err
		}
	}
	return string(body[:n-1]), opts, nil
}

func cStringOptsRespMarshal(str string, opts *Options) ([]byte, error) {
	out := append([]byte(str), 0)
	if opts == nil {
		return out, nil
	}
	o, err := opts.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return append(out, o...), nil
}

// Sequencer provides way of altering default PDU sequencing.
// This can be useful for load balancing requests.
type Sequencer interface {
	Next() uint32
}

// NewSequencer creates new sequencer with starting value set to n.
// Allowed range is 0x00000001 to 0x7FFFFFFF.
func NewSequencer(n uint32) Sequencer {
	if n == 0 {
		n = 1
	}
	return &defaultSequencer{n}
}

type defaultSequencer struct {
	n uint32
}

func (seq *defaultSequencer) Next() uint32 {
	n := seq.n
	seq.n++
	return n
}

// Encoder is responsible for encoding PDU structure to writer.
type Encoder struct {
	w   io.Writer
	seq Sequencer
}

// NewEncoder instantiates pdu encoder.
func NewEncoder(w io.Writer, seq Sequencer) *Encoder {
	if seq == nil {
		seq = NewSequencer(1)
	}
	return &Encoder{
		w:   w,
		seq: seq,
	}
}

type encoderOpts struct {
	seq    uint32
	status Status
}

// Encode PDU structure and write it to the assigned writer.
func (en *Encoder) Encode(p PDU, opts ...EncoderOption) (uint32, error) {
	// TODO consider introducing convention where pdu.MarshalBinary
	// should return slice with prepended space for header to avoid
	// allocation and copy.
	body, err := p.MarshalBinary()
	if err != nil {
		return 0, err
	}

	eOpts := encoderOpts{}
	for _, o := range opts {
		o(&eOpts)
	}

	l := len(body) + 16
	buf := make([]byte, l)
	binary.BigEndian.PutUint32(buf[:4], uint32(l))
	binary.BigEndian.PutUint32(buf[4:8], uint32(p.CommandID()))
	binary.BigEndian.PutUint32(buf[8:12], uint32(eOpts.status))
	if eOpts.seq == 0 {
		eOpts.seq = en.seq.Next()
	}
	binary.BigEndian.PutUint32(buf[12:16], eOpts.seq)
	copy(buf[16:], body)
	_, err = en.w.Write(buf)
	return eOpts.seq, err
}

type EncoderOption func(*encoderOpts)

func EncodeSeq(seq uint32) EncoderOption {
	return func(eOpts *encoderOpts) {
		eOpts.seq = seq
	}
}

func EncodeStatus(status Status) EncoderOption {
	return func(eOpts *encoderOpts) {
		eOpts.status = status
	}
}

// Decoder reads input from reader and marshals it into PDU.
type Decoder struct {
	r io.Reader
}

// NewDecoder initializes new PDU decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

// Decode reads data from reader and populates PDU.
func (d *Decoder) Decode() (Header, PDU, error) {
	// Read header first.
	var headerBytes [16]byte
	if _, err := io.ReadFull(d.r, headerBytes[:]); err != nil {
		return nil, nil, err
	}

	header := &header{}
	if err := header.UnmarshalBinary(headerBytes[:]); err != nil {
		return header, nil, err
	}
	// TODO: || header.length > data.MAX_PDU_LEN
	if header.length < 16 {
		return header, nil, fmt.Errorf("smpp: invalid pdu header byte length: %d", header.length)
	}

	pdu := NewPDU(header.commandID)
	if header.length == 16 {
		// not expecting body to read - we're done.
		return header, pdu, nil
	}

	// Read rest of the PDU.
	bodyBytes := make([]byte, header.length-16)
	if len(bodyBytes) > 0 {
		if _, err := io.ReadFull(d.r, bodyBytes); err != nil {
			return header, pdu, fmt.Errorf("smpp: pdu length doesn't match read body length %d != %d", header.length, len(bodyBytes))
		}
	}

	// Unmarshal binary
	if err := pdu.UnmarshalBinary(bodyBytes); err != nil {
		return header, pdu, err
	}

	return header, pdu, nil
}

// NewPDU creates new PDU from CommandID.
func NewPDU(commandID CommandID) PDU {
	switch commandID {
	case GenericNackID:
		return &GenericNack{}
	case BindReceiverID:
		return &BindRx{}
	case BindReceiverRespID:
		return &BindRxResp{}
	case BindTransmitterID:
		return &BindTx{}
	case BindTransmitterRespID:
		return &BindTxResp{}
	case BindTransceiverID:
		return &BindTRx{}
	case BindTransceiverRespID:
		return &BindTRxResp{}
	case EnquireLinkID:
		return &EnquireLink{}
	case EnquireLinkRespID:
		return &EnquireLinkResp{}
	case QuerySmID:
		return &QuerySm{}
	case QuerySmRespID:
		return &QuerySmResp{}
	case SubmitSmID:
		return &SubmitSm{}
	case SubmitSmRespID:
		return &SubmitSmResp{}
	case DeliverSmID:
		return &DeliverSm{}
	case DeliverSmRespID:
		return &DeliverSmResp{}
	case UnbindID:
		return &Unbind{}
	case UnbindRespID:
		return &UnbindResp{}
	case ReplaceSmID:
		return &ReplaceSm{}
	case ReplaceSmRespID:
		return &ReplaceSmResp{}
	case CancelSmID:
		return &CancelSm{}
	case CancelSmRespID:
		return &CancelSmResp{}
	case OutbindID:
		return &Outbind{}
	case SubmitMultiID:
		return &SubmitMulti{}
	case SubmitMultiRespID:
		return &SubmitMultiResp{}
	case AlertNotificationID:
		return &AlertNotification{}
	case DataSmID:
		return &DataSm{}
	case DataSmRespID:
		return &DataSmResp{}
	}
	panic("pdu: unsupported PDU command")
}

// IsRequest returns true if command is request.
func IsRequest(id CommandID) bool {
	switch id {
	default:
		return true
	case GenericNackID,
		BindReceiverRespID,
		BindTransmitterRespID,
		QuerySmRespID,
		SubmitSmRespID,
		DeliverSmRespID,
		UnbindRespID,
		ReplaceSmRespID,
		CancelSmRespID,
		BindTransceiverRespID,
		EnquireLinkRespID,
		SubmitMultiRespID,
		DataSmRespID:
		return false
	}
}

// SystemID extracts system id value from PDU if it has one.
func SystemID(p PDU) string {
	switch p.CommandID() {
	case BindReceiverID:
		if p, ok := p.(*BindRx); ok {
			return p.SystemID
		}
	case BindTransmitterID:
		if p, ok := p.(*BindTx); ok {
			return p.SystemID
		}
	case BindTransceiverID:
		if p, ok := p.(*BindTRx); ok {
			return p.SystemID
		}
	case BindReceiverRespID:
		if p, ok := p.(*BindRxResp); ok {
			return p.SystemID
		}
	case BindTransmitterRespID:
		if p, ok := p.(*BindTxResp); ok {
			return p.SystemID
		}
	case BindTransceiverRespID:
		if p, ok := p.(*BindTRxResp); ok {
			return p.SystemID
		}
	}
	return ""
}

// SeparateUDH takes input bytes and separates them into UDH header and content.
func SeparateUDH(c []byte) ([]byte, []byte, error) {
	if len(c) == 0 {
		return nil, c, errors.New("smpp: invalid udh length")
	}
	l := int(c[0])
	if l >= len(c) {
		return nil, c, errors.New("smpp: invalid udh length value")
	}
	return c[:l+1], c[l+1:], nil
}
