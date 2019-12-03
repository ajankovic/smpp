package pdu

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"strings"
	"testing"
)

var pduTT = []struct {
	desc   string
	hexStr string
	pdu    PDU
	err    bool
}{
	{
		"valid submit_sm pdu",
		"00|00|00|7465737400|00|00|746573743200|00|00|00|00|00|00|00|00|00|03|6d7367",
		&SubmitSm{
			SourceAddr:      "test",
			DestinationAddr: "test2",
			ShortMessage:    "msg",
		},
		false,
	},
	{
		"valid submit_sm with long message",
		"00010161736466000101333831363331323334353400000001000000000100f76161736466617364666173646661736466206173646661736466617364666173646661207364666173642066612073646620617364206661207364666173642066612064666173646661736466617364666173646620617364666173646661736466617364666120736466617364206661207364662061736420666120736466617364206661206466617364666173646661736466617364666173646661736431313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313102040002006f",
		&SubmitSm{
			SourceAddrTon:   0x01,
			SourceAddrNpi:   0x01,
			SourceAddr:      "asdf",
			DestAddrTon:     0x01,
			DestAddrNpi:     0x01,
			DestinationAddr: "38163123454",
			PriorityFlag:    0x01,
			DataCoding:      0x01,
			ShortMessage:    "aasdfasdfasdfasdf asdfasdfasdfasdfa sdfasd fa sdf asd fa sdfasd fa dfasdfasdfasdfasdf asdfasdfasdfasdfa sdfasd fa sdf asd fa sdfasd fa dfasdfasdfasdfasdfasdfasd111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			Options:         NewOptions().SetUserMessageReference(0x6F),
		},
		false,
	},
	{
		"valid deliver_sm with long message",
		"00010161736466000101333831363331323334353400000001000000000100f76161736466617364666173646661736466206173646661736466617364666173646661207364666173642066612073646620617364206661207364666173642066612064666173646661736466617364666173646620617364666173646661736466617364666120736466617364206661207364662061736420666120736466617364206661206466617364666173646661736466617364666173646661736431313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313131313102040002006f",
		&DeliverSm{
			SourceAddrTon:   0x01,
			SourceAddrNpi:   0x01,
			SourceAddr:      "asdf",
			DestAddrTon:     0x01,
			DestAddrNpi:     0x01,
			DestinationAddr: "38163123454",
			PriorityFlag:    0x01,
			DataCoding:      0x01,
			ShortMessage:    "aasdfasdfasdfasdf asdfasdfasdfasdfa sdfasd fa sdf asd fa sdfasd fa dfasdfasdfasdfasdf asdfasdfasdfasdfa sdfasd fa sdf asd fa sdfasd fa dfasdfasdfasdfasdfasdfasd111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			Options:         NewOptions().SetUserMessageReference(0x6F),
		},
		false,
	},
	{
		"valid bind_trx pdu",
		"7465737400|746573743200|00|00|01|01|00",
		&BindTRx{
			SystemID: "test",
			Password: "test2",
			AddrTon:  1,
			AddrNpi:  1,
		},
		false,
	},
	{
		"valid query_sm pdu",
		"7465737400|01|01|6173646600",
		&QuerySm{
			MessageID:     "test",
			SourceAddrTon: 0x01,
			SourceAddrNpi: 0x01,
			SourceAddr:    "asdf",
		},
		false,
	},
	{
		"valid empty unbind pdu",
		"",
		&Unbind{},
		false,
	},
	{
		"valid bind_trx_resp pdu",
		"7465737400|0210|0001|34",
		&BindTRxResp{
			SystemID: "test",
			Options:  NewOptions().SetScInterfaceVersion(0x34),
		},
		false,
	},
	// Always append new cases to avoid messing up Encoding/Decoding tests which
	// rely on indexes in this table.
}

func toHexStr(s string) string {
	return strings.Replace(s, "|", "", -1)
}

func TestMarshalBinary(t *testing.T) {
	for _, row := range pduTT {
		t.Run(row.desc, func(t *testing.T) {
			b, err := row.pdu.MarshalBinary()
			if err != nil {
				if !row.err {
					t.Fatalf("unexpected error %s", err)
				}
				return
			}
			written := hex.EncodeToString(b)
			if written != toHexStr(row.hexStr) {
				t.Errorf("MarshalBinary() => %q\nExpected: %q\nErr: %v", written, toHexStr(row.hexStr), err)
			}
		})
	}
}

func TestUnmarshalBinary(t *testing.T) {
	for _, row := range pduTT {
		t.Run(row.desc, func(t *testing.T) {
			data, _ := hex.DecodeString(toHexStr(row.hexStr))
			p := reflect.New(reflect.TypeOf(row.pdu).Elem()).Interface().(PDU)
			err := p.UnmarshalBinary(data)
			if err != nil {
				if !row.err {
					t.Fatalf("unexpected error %s", err)
				}
				return
			}
			if !reflect.DeepEqual(p, row.pdu) {
				t.Errorf("UnmarshalBinary(p) => \n%+v\nExpected: \n%+v", p, row.pdu)
			}
		})
	}
}

func BenchmarkSubmitSm_MarshalBinary(b *testing.B) {
	b.SetBytes(285)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bin, err := pduTT[1].pdu.MarshalBinary()
		if err != nil {
			b.Fatalf("error with marshaling %v", err)
		}
		_ = bin
	}
}

func BenchmarkSubmitSm_UnmarshalBinary(b *testing.B) {
	in, _ := hex.DecodeString(toHexStr(pduTT[1].hexStr))
	b.SetBytes(int64(len(in)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pdu := &SubmitSm{}
		err := pdu.UnmarshalBinary(in)
		if err != nil {
			b.Fatalf("error with unmarshaling %v", err)
		}
		_ = pdu
	}
}

func TestSeparateUDH(t *testing.T) {
	udhtest, _ := hex.DecodeString("0B0504158200000003AA0301")
	b, _ := hex.DecodeString("0B0504158200000003AA030174657374")
	udh, content, err := SeparateUDH(b)
	if err != nil {
		t.Fatalf("separate udh %v", err)
	}
	if !bytes.Equal(udh, udhtest) {
		t.Errorf("separate udh got %X expected %X", udh, udhtest)
	}
	if string(content) != "test" {
		t.Errorf("separate udh got %X expected %X", content, "test")
	}
}

var codingTT = []struct {
	desc      string
	headerHex string
	sequencer Sequencer
	pduIndex  int
	status    Status
	seq       uint32
	err       bool
}{
	{
		"submit_sm with default sequencer",
		"0000002D|00000004|00000000|00000001",
		nil, // Default sequencer.
		0,   // From pduTT.
		StatusOK,
		1,
		false,
	},
	{
		"submit_sm with custom sequencer",
		"0000002D|00000004|00000000|00000003",
		NewSequencer(3),
		0,
		StatusOK,
		3,
		false,
	},
	{
		"submit_sm with sequence number",
		"0000002D|00000004|00000000|00000004",
		nil,
		0,
		StatusOK,
		4,
		false,
	},
	{
		"unbind with empty body",
		"00000010|00000006|00000000|00000001",
		nil,
		5,
		StatusOK,
		1,
		false,
	},
	{
		"unbind with custom status",
		"00000010|00000006|00000004|00000001",
		nil,
		5,
		StatusInvBnd,
		1,
		false,
	},
	{
		"bindtrx resp with options",
		"0000001A|80000009|00000000|00000001",
		nil,
		6,
		StatusOK,
		1,
		false,
	},
}

func TestPDUEncoding(t *testing.T) {
	for _, row := range codingTT {
		t.Run(row.desc, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			enc := NewEncoder(buf, row.sequencer)

			opts := []EncoderOption{EncodeStatus(row.status)}
			if row.sequencer == nil {
				opts = append(opts, EncodeSeq(row.seq))
			}
			i, err := enc.Encode(pduTT[row.pduIndex].pdu, opts...)
			if err != nil {
				if !row.err {
					t.Fatalf("unexpected error %s", err)
				}
				return
			}
			if i != row.seq {
				t.Errorf("Encode() => seq %d expected %d", i, row.seq)
			}
			expected, _ := hex.DecodeString(toHexStr(row.headerHex + pduTT[row.pduIndex].hexStr))
			got := buf.Bytes()
			if !bytes.Equal(expected, got) {
				t.Errorf("Encode() => bytes\n%X\nexpected \n%X", got, expected)
			}
		})
	}
}

func TestPDUDecoding(t *testing.T) {
	for _, row := range codingTT {
		t.Run(row.desc, func(t *testing.T) {
			expected, _ := hex.DecodeString(toHexStr(row.headerHex + pduTT[row.pduIndex].hexStr))
			buf := bytes.NewBuffer(expected)
			dec := NewDecoder(buf)
			h, p, err := dec.Decode()
			if err != nil {
				if !row.err {
					t.Fatalf("unexpected error %s", err)
				}
				return
			}
			if h.Sequence() != row.seq {
				t.Errorf("Decode() => seq %d expected %d", h.Sequence(), row.seq)
			}
			if h.Status() != row.status {
				t.Errorf("Decode() => status %d expected %d", h.Status(), row.status)
			}
			if !reflect.DeepEqual(p, pduTT[row.pduIndex].pdu) {
				t.Errorf("Decode() => pdu\n%+v\nexpected \n%+v", p, pduTT[row.pduIndex].pdu)
			}
		})
	}
}
