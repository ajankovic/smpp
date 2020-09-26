package smpp_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/daominah/smpp"
	"github.com/daominah/smpp/pdu"
)

const (
	TestAddr = ":30303"
)

func TestSMPPServer(t *testing.T) {
	sessConf := smpp.SessionConf{
		Handler: smpp.HandlerFunc(func(ctx *smpp.Context) {
			switch ctx.CommandID() {
			case pdu.BindTransceiverID:
				btrx, err := ctx.BindTRx()
				if err != nil {
					t.Errorf(err.Error())
				}
				resp := btrx.Response("TestingServer")
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					t.Errorf(err.Error())
				}
			}
		}),
	}
	srv := smpp.NewServer(TestAddr, sessConf)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			t.Errorf("Expected no error on server close %v", err)
		}
	}()
	time.Sleep(time.Millisecond * 10)
	sess1 := bindToServer(TestAddr, smpp.HandlerFunc(func(ctx *smpp.Context) {
		switch ctx.CommandID() {
		case pdu.UnbindID:
			ubd, err := ctx.Unbind()
			if err != nil {
				t.Errorf(err.Error())
			}
			resp := ubd.Response()
			if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
				t.Errorf(err.Error())
			}
		}
	}))
	sess2 := bindToServer(TestAddr, smpp.HandlerFunc(func(ctx *smpp.Context) {
		switch ctx.CommandID() {
		case pdu.UnbindID:
			ubd, err := ctx.Unbind()
			if err != nil {
				t.Errorf(err.Error())
			}
			resp := ubd.Response()
			if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
				t.Errorf(err.Error())
			}
		}
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := srv.Unbind(ctx)
	if err != nil {
		t.Error(err.Error())
	}
	select {
	case <-sess1.NotifyClosed():
	case <-time.After(100 * time.Millisecond):
		t.Errorf("session %s was not closed in time", sess1)
	}
	select {
	case <-sess2.NotifyClosed():
	case <-time.After(100 * time.Millisecond):
		t.Errorf("session %s was not closed in time", sess2)
	}
}

func bindToServer(bind string, hf smpp.HandlerFunc) *smpp.Session {
	bc := smpp.BindConf{
		Addr:     bind,
		SystemID: "Client",
		Password: "password",
	}
	sc := smpp.SessionConf{
		Handler: hf,
	}
	sess, err := smpp.BindTRx(sc, bc)
	if err != nil {
		log.Fatalf("error during bind %v", err)
	}
	return sess
}
