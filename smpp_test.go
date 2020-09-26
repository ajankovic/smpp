package smpp_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/daominah/smpp"
	"github.com/daominah/smpp/pdu"
)

type mockServer struct {
	Addr    string
	Respond func(c net.Conn, in pdu.PDU, i int) []byte
}

func startServer(server *mockServer, n int) {
	l, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	tcpConn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer tcpConn.Close()

	for i := 0; i < n; i++ {
		server.Serve(tcpConn, i)
	}
}

func (this *mockServer) Serve(c net.Conn, i int) {
	d := pdu.NewDecoder(c)
	_, p, err := d.Decode()
	if err != nil {
		if err != io.EOF {
			log.Fatalf("serve decode %v %d", err, i)
		}
		return
	}
	if p == nil {
		log.Fatal("decode returned nil")
	}
	res := this.Respond(c, p, i)
	if res == nil {
		return
	}
	if _, err := c.Write(res); err != nil {
		log.Fatalf("connection write %v", err)
	}
}

func newBindingServer() *mockServer {
	b := &bytes.Buffer{}
	e := pdu.NewEncoder(b, nil)
	return &mockServer{
		Addr: "localhost:2222",
		Respond: func(c net.Conn, in pdu.PDU, i int) []byte {
			var res pdu.PDU
			switch in.CommandID() {
			case pdu.BindTransceiverID:
				res = &pdu.BindTRxResp{
					SystemID: "testing",
					Options:  pdu.NewOptions().SetScInterfaceVersion(0x34),
				}
			case pdu.UnbindID:
				res = &pdu.UnbindResp{}
			}
			b.Reset()
			if _, err := e.Encode(res); err != nil {
				panic("Can't encode pdu")
			}
			return b.Bytes()
		},
	}
}

func TestBindingUnbinding(t *testing.T) {
	finished := make(chan struct{})
	server := newBindingServer()
	go func() {
		startServer(server, 2)
		finished <- struct{}{}
	}()
	time.Sleep(time.Millisecond * 10)
	conf := smpp.BindConf{
		Addr: "localhost:2222",
	}
	sess, err := smpp.BindTRx(smpp.SessionConf{}, conf)
	if err != nil {
		t.Errorf("bind error %s", err)
	}
	if sess.SystemID() != "testing" {
		t.Errorf("Invalid SystemID after bind %s", sess.SystemID())
	}
	err = smpp.Unbind(context.Background(), sess)
	if err != nil {
		t.Errorf("unbind error %s", err)
	}
	select {
	case <-sess.NotifyClosed():
	case <-time.After(100 * time.Millisecond):
		t.Error("session close timeout")
	}
	select {
	case <-finished:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("mock server didn't close")
	}
}

func TestBindToDeadEnd(t *testing.T) {
	conf := smpp.BindConf{
		Addr: "localhost:8484",
	}
	sess, err := smpp.BindTRx(smpp.SessionConf{}, conf)
	if err == nil {
		t.Errorf("expected error bot got nil")
	}
	if sess != nil {
		t.Errorf("expected session to be nil got %s", sess)
	}
}
