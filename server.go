package smpp

import (
	"context"
	"net"
	"sync"
	"time"
)

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// Server implements SMPP SMSC server.
type Server struct {
	Addr        string
	SessionConf *SessionConf

	wg         sync.WaitGroup
	mu         sync.Mutex
	listeners  map[net.Listener]struct{}
	doneChan   chan struct{}
	activeSess map[*Session]struct{}
}

// NewServer creates new SMPP server for managing SMSC sessions.
// Sessions will use provided SessionConf as template configuration.
func NewServer(addr string, conf SessionConf) *Server {
	return &Server{
		Addr:        addr,
		SessionConf: &conf,
	}
}

// ListenAndServe starts server listening. Blocking function.
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":2775"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

// Serve accepts incoming connections and starts SMPP sessions.
func (srv *Server) Serve(ln net.Listener) error {
	defer ln.Close()
	srv.trackListener(ln, true)
	// How long to sleep on accept failure.
	var tempDelay time.Duration
	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-srv.getDoneChan():
				return nil
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0

		srv.wg.Add(1)
		go func(conf SessionConf) {
			defer srv.wg.Done()
			conf.Type = SMSC
			sess := NewSession(conn, conf)
			srv.trackSess(sess, true)
			select {
			case <-sess.NotifyClosed():
			case <-srv.getDoneChan():
				sess.Close()
			}
			srv.trackSess(sess, false)
		}(*srv.SessionConf)
	}
}

// Unbind gracefully closes server by sending Unbind requests to all connected peers.
func (srv *Server) Unbind(ctx context.Context) error {
	srv.mu.Lock()
	for sess := range srv.activeSess {
		Unbind(ctx, sess)
	}
	srv.mu.Unlock()
	return srv.Close()
}

// Close implements closer interface.
func (srv *Server) Close() error {
	srv.mu.Lock()
	srv.closeDoneChanLocked()
	err := srv.closeListenersLocked()
	srv.mu.Unlock()
	srv.wg.Wait()
	return err
}

func (srv *Server) getDoneChan() <-chan struct{} {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.getDoneChanLocked()
}

func (srv *Server) getDoneChanLocked() chan struct{} {
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	return srv.doneChan
}

func (srv *Server) closeDoneChanLocked() {
	ch := srv.getDoneChanLocked()
	select {
	case <-ch:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded by srv.mu.
		close(ch)
	}
}

func (srv *Server) trackListener(ln net.Listener, add bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.listeners == nil {
		srv.listeners = make(map[net.Listener]struct{})
	}
	if add {
		// If the *Server is being reused after a previous
		// Close or Shutdown, reset its doneChan:
		if len(srv.listeners) == 0 && len(srv.activeSess) == 0 {
			srv.doneChan = nil
		}
		srv.listeners[ln] = struct{}{}
	} else {
		delete(srv.listeners, ln)
	}
}

func (srv *Server) closeListenersLocked() error {
	var err error
	for ln := range srv.listeners {
		if cerr := ln.Close(); cerr != nil && err == nil {
			err = cerr
		}
		delete(srv.listeners, ln)
	}
	return err
}

func (srv *Server) trackSess(sess *Session, add bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.activeSess == nil {
		srv.activeSess = make(map[*Session]struct{})
	}
	if add {
		srv.activeSess[sess] = struct{}{}
	} else {
		delete(srv.activeSess, sess)
	}
}
