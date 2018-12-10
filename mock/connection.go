// Package mock implements necessary mocking structures to allow easier
// testing for the smpp package.
// Only ReadWriteCloser is currently implemented aimed at mocking network
// connection.
package mock

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

const (
	readR          = "read"
	writeR         = "write"
	countWritesOff = -90000000
)

type step struct {
	request      string
	write        []byte
	read         []byte
	err          error
	closed       bool
	waiting      bool
	done         bool
	noResp       bool
	wait         int
	count        int
	processRead  func(step int, count int) ([]byte, error)
	processWrite func(step int, count int) ([]byte, error)
}

// Conn implements ReadWriteCloser interface and can be used to mock
// network connection in tests.
type Conn struct {
	io.ReadWriteCloser
	done   chan struct{}
	mux    sync.Mutex
	errors []error
	steps  []*step
}

// NewConn instantiates mocked connection.
func NewConn() *Conn {
	return &Conn{
		done: make(chan struct{}),
	}
}

func (c *Conn) Read(out []byte) (int, error) {
	for {
		i, err := c.read(out)
		if i != -1 {
			return i, err
		}
		select {
		// Check if there is anything to read in regular intervals.
		case <-time.After(2 * time.Millisecond):
		case <-c.done:
			return 0, io.EOF
		}
	}
}

func (c *Conn) read(out []byte) (int, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	d := 0
	for i, s := range c.steps {
		if s.done || s.closed {
			d++
			if d == len(c.steps) {
				return -1, io.EOF
			}
			continue
		}
		if s.processRead != nil && s.read == nil {
			var err error
			s.read, err = s.processRead(i, s.count)
			if err != nil {
				return 0, err
			}
		}
		// Check for read requests that are not waiting for write.
		if s.request == readR && !s.waiting {
			if s.wait > 0 && !c.steps[s.wait-1].done {
				continue
			}
			if s.err != nil {
				s.done = true
				return 0, s.err
			}
			n := copy(out, s.read)
			if n < len(s.read) {
				s.read = s.read[n:]
				return n, nil
			}
			if s.noResp {
				s.done = true
			} else {
				s.waiting = true
			}
			return n, nil
		}
		// Check writes which are waiting for read.
		if s.request == writeR && s.waiting {
			if s.err != nil {
				s.done = true
				return 0, s.err
			}
			n := copy(out, s.read)
			if n < len(s.read) {
				s.read = s.read[n:]
				return n, nil
			}
			s.done = true
			return n, nil
		}
	}
	return -1, nil
}

// Write implements io.Writer interface.
func (c *Conn) Write(in []byte) (int, error) {
	for {
		select {
		// Write timeout.
		case <-time.After(2 * time.Millisecond):
		case <-c.done:
			return 0, io.EOF
		}
		i, err := c.write(in)
		if i != -1 {
			return i, err
		}
	}
}

func (c *Conn) write(in []byte) (int, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	for i, s := range c.steps {
		if s.done || s.closed {
			continue
		}
		if s.processWrite != nil && s.write == nil {
			var err error
			s.write, err = s.processWrite(i, s.count)
			if err != nil {
				return 0, err
			}
		}
		// Handle responses to read requests.
		if s.request == readR && s.waiting {
			if !bytes.Equal(s.write, in) {
				continue
			}
			s.done = true
			return len(in), nil
		}
		// Handle write requests.
		if s.request == writeR && !s.waiting {
			if s.wait > 0 && !c.steps[s.wait-1].done {
				continue
			}
			if s.err != nil {
				s.done = true
				return 0, s.err
			}
			if s.write != nil {
				if !bytes.Equal(s.write, in) {
					continue
				}
			}
			if s.noResp {
				s.done = true
			} else {
				s.waiting = true
			}
			return len(in), nil
		}
	}
	err := fmt.Errorf("mock: unexpected write\n% X", in)
	c.errors = append(c.errors, err)
	return 0, err
}

// Close implements closer interface.
func (c *Conn) Close() error {
	c.mux.Lock()
	defer c.mux.Unlock()
	close(c.done)
	closed := false
	done := true
	for _, s := range c.steps {
		closed = closed || s.closed
		if s.closed {
			s.done = true
		}
		done = done && s.done
	}
	if !closed {
		err := errors.New("mock: unexpected call to Close")
		c.errors = append(c.errors, err)
		return err
	}
	if !done {
		err := errors.New("mock: closing unfinished scenario")
		c.errors = append(c.errors, err)
		return err
	}
	return nil
}

// ByteRead will set connection to respond with provided bytes.
func (c *Conn) ByteRead(read []byte) *Conn {
	c.mux.Lock()
	defer c.mux.Unlock()
	l := len(c.steps)
	if l != 0 && c.steps[l-1].read == nil && c.steps[l-1].processRead == nil && !c.steps[l-1].noResp && c.steps[l-1].err == nil {
		c.steps[l-1].read = read
	} else {
		c.steps = append(c.steps, &step{request: readR, read: read})
	}
	return c
}

// ErrRead will set connection to fail read with provided error.
func (c *Conn) ErrRead(err error) *Conn {
	c.mux.Lock()
	defer c.mux.Unlock()
	if err == nil {
		err = errors.New("mock: generic read error")
	}
	l := len(c.steps)
	if l != 0 && c.steps[l-1].read == nil && c.steps[l-1].processRead == nil && !c.steps[l-1].noResp && c.steps[l-1].err == nil {
		c.steps[l-1].err = err
	} else {
		c.steps = append(c.steps, &step{request: readR, err: err})
	}
	return c
}

// ByteWrite will set connection to expect provided bytes for write.
func (c *Conn) ByteWrite(write []byte) *Conn {
	c.mux.Lock()
	defer c.mux.Unlock()
	l := len(c.steps)
	if l != 0 && c.steps[l-1].write == nil && !c.steps[l-1].noResp && c.steps[l-1].err == nil {
		c.steps[l-1].write = write
	} else {
		c.steps = append(c.steps, &step{request: writeR, write: write})
	}
	return c
}

// ErrWrite will set connection to fail write with error.
func (c *Conn) ErrWrite(err error) *Conn {
	c.mux.Lock()
	defer c.mux.Unlock()
	if err == nil {
		err = errors.New("mock: generic write error")
	}
	l := len(c.steps)
	if l != 0 && c.steps[l-1].write == nil && !c.steps[l-1].noResp && c.steps[l-1].err == nil {
		c.steps[l-1].err = err
	} else {
		c.steps = append(c.steps, &step{request: writeR, err: err})
	}
	return c
}

// NoResp will set connection to not return response for the preceding read/write.
// Panics if write/read calls are unmatched.
func (c *Conn) NoResp() *Conn {
	c.mux.Lock()
	defer c.mux.Unlock()
	l := len(c.steps)
	if l == 0 {
		panic("mock: invalid call to NoResp")
	}
	if c.steps[l-1].write != nil && c.steps[l-1].read != nil && c.steps[l-1].processRead == nil {
		panic("mock: invalid call to NoResp")
	}
	c.steps[l-1].noResp = true
	return c
}

// Closed will set connection to expect call to Close.
func (c *Conn) Closed() *Conn {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.steps = append(c.steps, &step{closed: true})
	return c
}

// Wait acts on the last chained step and it instructs mock to
// wait for the indexed step to complete before the one Wait was
// called on.
func (c *Conn) Wait(s int) *Conn {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.steps[len(c.steps)-1].wait = s
	return c
}

// ProcessRead provides function to the connection that will process read.
func (c *Conn) ProcessRead(f func(step, count int) ([]byte, error)) *Conn {
	c.mux.Lock()
	defer c.mux.Unlock()
	l := len(c.steps)
	if l != 0 && c.steps[l-1].read == nil && c.steps[l-1].processRead == nil && !c.steps[l-1].noResp && c.steps[l-1].err == nil {
		c.steps[l-1].processRead = f
		c.steps[l-1].count = 1
	} else {
		c.steps = append(c.steps, &step{request: readR, processRead: f, count: 1})
	}
	return c
}

// ProcessWrite provides function to the connection that will process writes.
func (c *Conn) ProcessWrite(f func(step, count int) ([]byte, error)) *Conn {
	c.mux.Lock()
	defer c.mux.Unlock()
	l := len(c.steps)
	if l != 0 && c.steps[l-1].write == nil && c.steps[l-1].processWrite == nil && !c.steps[l-1].noResp && c.steps[l-1].err == nil {
		c.steps[l-1].processWrite = f
		c.steps[l-1].count = 1
	} else {
		c.steps = append(c.steps, &step{request: writeR, processWrite: f, count: 1})
	}
	return c
}

// Times multiplies last step n times.
func (c *Conn) Times(n int) *Conn {
	c.mux.Lock()
	defer c.mux.Unlock()
	l := len(c.steps)
	if l == 0 {
		panic("mock: invalid call to Times")
	}
	c.steps[l-1].count = 1
	for i := 2; i <= n; i++ {
		s := *c.steps[l-1]
		c.steps[l-1].count = i
		c.steps = append(c.steps, &s)
	}
	return c
}

// Validate will check executed scenario and return any errors in execution.
// It will return nil if scenario was valid.
func (c *Conn) Validate() []error {
	c.mux.Lock()
	defer c.mux.Unlock()
	for _, s := range c.steps {
		if !s.done {
			var val string
			if s.closed {
				val = "closing connection"
			} else if s.request == readR {
				val = fmt.Sprintf("%s % X", s.request, s.read)
			} else if s.request == writeR {
				val = fmt.Sprintf("%s % X", s.request, s.write)
			}
			c.errors = append(c.errors, fmt.Errorf("mock: step not finished: %s", val))
		}
	}
	return c.errors
}
