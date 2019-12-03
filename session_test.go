package smpp_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/ajankovic/smpp"
	"github.com/ajankovic/smpp/internal/mock"
	"github.com/ajankovic/smpp/pdu"
)

type testSequencer struct {
	seq  uint32
	skip bool
}

func (ts *testSequencer) Next() uint32 {
	if !ts.skip {
		ts.seq++
	} else {
		ts.skip = false
	}
	return ts.seq
}

func (ts *testSequencer) skipNext() {
	ts.skip = true
}

type testEncoder struct {
	buf *bytes.Buffer
	enc *pdu.Encoder
	seq *testSequencer
}

func newTestEncoder(i int) *testEncoder {
	buf := bytes.NewBuffer(nil)
	seq := &testSequencer{seq: uint32(i)}
	return &testEncoder{
		buf: buf,
		seq: seq,
		enc: pdu.NewEncoder(buf, seq),
	}
}

// Encode by incrementing counter.
func (te *testEncoder) i(p pdu.PDU, status ...pdu.Status) []byte {
	te.buf.Reset()
	st := pdu.StatusOK
	if len(status) > 0 {
		st = status[0]
	}
	_, err := te.enc.Encode(p, st)
	if err != nil {
		panic(err.Error())
	}
	out := make([]byte, te.buf.Len())
	copy(out, te.buf.Bytes())
	return out
}

// Encode by skipping increment.
func (te *testEncoder) s(p pdu.PDU, status ...pdu.Status) []byte {
	te.buf.Reset()
	st := pdu.StatusOK
	if len(status) > 0 {
		st = status[0]
	}
	te.seq.skipNext()
	_, err := te.enc.Encode(p, st)
	if err != nil {
		panic(err.Error())
	}
	out := make([]byte, te.buf.Len())
	copy(out, te.buf.Bytes())
	return out
}

func TestESMESession(t *testing.T) {
	bindTRx := &pdu.BindTRx{
		SystemID:         "ESME",
		Password:         "password",
		SystemType:       "type",
		InterfaceVersion: smpp.Version,
		AddressRange:     "111111",
	}
	bindTRxResp := bindTRx.Response("SMSC")
	bindTRxResp.Options = pdu.NewOptions().SetScInterfaceVersion(smpp.Version)
	submitSm := &pdu.SubmitSm{
		SourceAddr:      "source",
		DestinationAddr: "destination",
		ShortMessage:    "this is the message",
	}
	submitSmResp := submitSm.Response("id0")
	unbind := pdu.Unbind{}
	unbindResp := pdu.UnbindResp{}
	e := newTestEncoder(0)
	conn := mock.NewConn().
		ByteWrite(e.i(bindTRx)).ByteRead(e.s(bindTRxResp)).
		ByteWrite(e.i(submitSm)).ByteRead(e.s(submitSmResp)).
		Wait(1).
		ByteWrite(e.i(unbind)).ByteRead(e.s(unbindResp)).
		Wait(1).
		Closed()
	conf := smpp.SessionConf{
		SystemID: "TestingESME",
	}
	sess := smpp.NewSession(conn, conf)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	resp, err := sess.Send(ctx, bindTRx)
	if err != nil {
		t.Fatal(err)
	}
	if resp.CommandID() != pdu.BindTransceiverRespID {
		t.Errorf("expected BindTransceiverRespID got %d", resp.CommandID())
	}
	resp, err = sess.Send(ctx, submitSm)
	if err != nil {
		t.Fatal(err)
	}
	if resp.CommandID() != pdu.SubmitSmRespID {
		t.Errorf("expected SubmitSmRespID got %d", resp.CommandID())
	}
	resp, err = sess.Send(ctx, unbind)
	if err != nil {
		t.Fatal(err)
	}
	if resp.CommandID() != pdu.UnbindRespID {
		t.Errorf("expected UnbindRespID got %d", resp.CommandID())
	}
	if err := sess.Close(); err != nil {
		t.Errorf("Got error during session close %+v", err)
	}
	errors := conn.Validate()
	if errors != nil {
		for _, err := range errors {
			t.Error(err)
		}
	}
}

func TestESMESessionInvalidStatus(t *testing.T) {
	bindTRx := &pdu.BindTRx{
		SystemID: "ESME",
	}
	bindTRxResp := bindTRx.Response("SMSC")
	submitSm := &pdu.SubmitSm{
		SourceAddr:      "source",
		DestinationAddr: "destination",
		ShortMessage:    "this is the message",
	}
	submitSmResp := submitSm.Response("id0")
	e := newTestEncoder(0)
	conn := mock.NewConn().
		ByteWrite(e.i(bindTRx)).ByteRead(e.s(bindTRxResp)).
		ByteWrite(e.i(submitSm)).ByteRead(e.s(submitSmResp, pdu.StatusInvDstAdr)).
		Wait(1).
		Closed()
	conf := smpp.SessionConf{}
	sess := smpp.NewSession(conn, conf)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	resp, err := sess.Send(ctx, bindTRx)
	if err != nil {
		t.Fatal(err)
	}
	if resp.CommandID() != pdu.BindTransceiverRespID {
		t.Errorf("expected BindTransceiverRespID got %d", resp.CommandID())
	}
	resp, err = sess.Send(ctx, submitSm)
	if err == nil {
		t.Errorf("Expected status error got nil")
	}
	if resp.CommandID() != pdu.SubmitSmRespID {
		t.Errorf("expected SubmitSmRespID got %d", resp.CommandID())
	}
	if serr, ok := err.(smpp.StatusError); !ok {
		t.Errorf("Expected StatusError type")
	} else {
		expected := "Invalid Destination Address '0xB'"
		if serr.Error() != expected {
			t.Errorf("Status error: %v, expected %s", err, expected)
		}
	}
	if err := sess.Close(); err != nil {
		t.Errorf("Got error during session close %+v", err)
	}
	errors := conn.Validate()
	if errors != nil {
		for _, err := range errors {
			t.Error(err)
		}
	}
}

func TestSMSCSession(t *testing.T) {
	bindTRx := &pdu.BindTRx{
		SystemID:         "ESME",
		Password:         "password",
		SystemType:       "type",
		InterfaceVersion: smpp.Version,
		AddressRange:     "111111",
	}
	bindTRxResp := bindTRx.Response("SMSC")
	bindTRxResp.Options = pdu.NewOptions().SetScInterfaceVersion(smpp.Version)

	submitSm := &pdu.SubmitSm{
		SourceAddr:      "source",
		DestinationAddr: "destination",
		ShortMessage:    "this is the message",
	}
	submitSmResp := submitSm.Response("id0")

	sync := make(chan struct{})
	e := newTestEncoder(0)
	conn := mock.NewConn().
		ByteRead(e.i(bindTRx, pdu.StatusOK)).ByteWrite(e.s(bindTRxResp, pdu.StatusOK)).
		ByteRead(e.i(submitSm, pdu.StatusOK)).ByteWrite(e.s(submitSmResp, pdu.StatusOK)).Wait(1).
		Closed()
	conf := smpp.SessionConf{
		SystemID: "TestingSMSC",
		Type:     smpp.SMSC,
		Handler: smpp.HandlerFunc(func(ctx *smpp.Context) {
			switch ctx.CommandID() {
			case pdu.BindTransceiverID:
				btrx, err := ctx.BindTRx()
				if err != nil {
					t.Errorf("Handler can't get BindTRx request %v", err)
				}
				resp := btrx.Response("SMSC")
				resp.Options = pdu.NewOptions().SetScInterfaceVersion(smpp.Version)
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					t.Errorf("Handler can't respond to bind request %v", err)
				}
			case pdu.SubmitSmID:
				defer close(sync)
				sm, err := ctx.SubmitSm()
				if err != nil {
					t.Errorf("Handler can't get BindTRx request %v", err)
				}
				resp := sm.Response("id0")
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					t.Errorf("Handler can't respond to SubmitSm request %v", err)
				}
			}
		}),
	}
	sess := smpp.NewSession(conn, conf)
	select {
	case <-time.After(50 * time.Millisecond):
		t.Fatal("timeout waiting for response")
	case <-sync:
	}
	sess.Close()
	errors := conn.Validate()
	if errors != nil {
		for _, err := range errors {
			t.Error(err)
		}
	}
}
