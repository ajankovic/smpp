package mock

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestConnUnfinishedScenario(t *testing.T) {
	c := NewConn().
		ByteRead([]byte{1}).ByteWrite([]byte{1})
	errs := c.Validate()
	if errs == nil && len(errs) != 1 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestConnFinishedScenario(t *testing.T) {
	c := NewConn().
		ByteRead([]byte{1}).ByteWrite([]byte{1})
	out := make([]byte, 1)
	n, err := c.Read(out)
	if err != nil || n != 1 {
		t.Fatalf("Invalid read results %d %v", n, err)
	}
	n, err = c.Write(out)
	if err != nil || n != 1 {
		t.Fatalf("Invalid write results %d %v", n, err)
	}
	errs := c.Validate()
	if errs != nil {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestConnComplexScenario(t *testing.T) {
	c := NewConn().
		ByteRead([]byte{1}).ByteWrite([]byte{1}).
		ByteWrite([]byte{2}).NoResp().
		ByteWrite([]byte{3}).ByteRead([]byte{3}).
		ErrRead(nil).
		ErrWrite(errors.New("test")).
		Closed()
	out := make([]byte, 1)
	n, err := c.Read(out)
	if err != nil || n != 1 {
		t.Fatalf("Invalid read results %d %v", n, err)
	}
	n, err = c.Write(out)
	if err != nil || n != 1 {
		t.Fatalf("Invalid write results %d %v", n, err)
	}
	n, err = c.Write([]byte{2})
	if err != nil || n != 1 {
		t.Fatalf("Invalid write results %d %v", n, err)
	}
	n, err = c.Write([]byte{3})
	if err != nil || n != 1 {
		t.Fatalf("Invalid write results %d %v", n, err)
	}
	out = make([]byte, 1)
	n, err = c.Read(out)
	if err != nil || n != 1 {
		t.Fatalf("Invalid read results %d %v", n, err)
	}
	n, err = c.Read(out)
	if err == nil {
		t.Error("expected read error")
	}
	n, err = c.Write(out)
	if err == nil {
		t.Error("expected write error")
	}
	c.Close()
	errs := c.Validate()
	if errs != nil {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestConnPartialRead(t *testing.T) {
	read := []byte{1, 1, 2, 2}
	c := NewConn().
		ByteRead(read).NoResp()
	out := make([]byte, 2)
	n, err := c.Read(out)
	if err != nil || n != 2 {
		t.Fatalf("Invalid read results %d %v", n, err)
	}
	out2 := make([]byte, 2)
	n, err = c.Read(out2)
	if err != nil || n != 2 {
		t.Fatalf("Invalid read results %d %v", n, err)
	}
	out = append(out, out2...)
	if !reflect.DeepEqual(out, read) {
		t.Fatalf("Reads are not equal\n% x\n!=\n% x", out, read)
	}
	errs := c.Validate()
	if errs != nil {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestConnWait(t *testing.T) {
	c := NewConn().
		ByteRead([]byte{1}).ByteWrite([]byte{2}).
		ByteRead([]byte{3}).ByteWrite([]byte{4}).Wait(1).
		ByteRead([]byte{5}).ByteWrite([]byte{6}).Wait(2)
	sync := make(chan struct{})
	read := []byte{1, 3, 5}
	write := []byte{2, 4, 6}
	out := make([]byte, 3)
	go func() {
		for i := range read {
			r := make([]byte, 1)
			n, err := c.Read(r)
			if err != nil || n != 1 {
				t.Fatalf("Invalid read results %d %v", n, err)
			}
			out[i] = r[0]
			sync <- struct{}{}
		}
	}()
	for i := range write {
		select {
		case <-sync:
			// next sync must not execute before write
			select {
			case <-sync:
				t.Fatal("Wait was not respected")
			case <-time.After(5 * time.Millisecond):
				if _, err := c.Write([]byte{write[i]}); err != nil {
					t.Fatal(err)
				}
			}
		case <-time.After(10 * time.Millisecond):
			t.Fatal("Timeout waiting for sync")
		}
	}
	if !bytes.Equal(out, read) {
		t.Fatalf("Reads are not equal\n% x\n!=\n% x", out, read)
	}
	errs := c.Validate()
	if errs != nil {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestConnWaitWriteFirst(t *testing.T) {
	c := NewConn().
		ByteWrite([]byte{2}).ByteRead([]byte{1}).
		ByteWrite([]byte{4}).ByteRead([]byte{3}).Wait(1).
		ByteWrite([]byte{6}).ByteRead([]byte{5}).Wait(2)
	sync := make(chan struct{})
	read := []byte{1, 3, 5}
	write := []byte{2, 4, 6}
	out := make([]byte, 3)
	go func() {
		for i := range write {
			n, err := c.Write([]byte{write[i]})
			if err != nil || n != 1 {
				t.Fatalf("Invalid Write results %d %v", n, err)
			}
			sync <- struct{}{}
		}
	}()
	for i := range read {
		select {
		case <-sync:
			r := make([]byte, 1)
			if _, err := c.Read(r); err != nil {
				t.Fatal(err)
			}
			out[i] = r[0]
		case <-time.After(10 * time.Millisecond):
			t.Fatal("Timeout waiting for sync")
		}
	}
	if !bytes.Equal(out, read) {
		t.Fatalf("Reads are not equal\n% x\n!=\n% x", out, read)
	}
	errs := c.Validate()
	if errs != nil {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestConnProcessWrites(t *testing.T) {
	c := NewConn().
		ProcessWrite(func(step, count int) ([]byte, error) {
			return []byte{byte(count)}, nil
		}).
		ProcessRead(func(step, count int) ([]byte, error) {
			return []byte{byte(count)}, nil
		}).Times(10)
	sync := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			_, err := c.Write([]byte{byte(i + 1)})
			if err != nil {
				t.Fatal(err)
			}
		}
		sync <- struct{}{}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			out := make([]byte, 1)
			_, err := c.Read(out)
			if err != nil {
				t.Fatal(err)
			}
		}
		sync <- struct{}{}
	}()
	for i := 0; i < 2; i++ {
		select {
		case <-sync:
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("Read write timeout %d", i)
		}
	}
	errs := c.Validate()
	if errs != nil {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestConnProcessReads(t *testing.T) {
	c := NewConn().
		ProcessRead(func(step, count int) ([]byte, error) {
			return []byte{byte(count)}, nil
		}).
		ProcessWrite(func(step, count int) ([]byte, error) {
			return []byte{byte(count)}, nil
		}).Times(10)
	for i := 0; i < 10; i++ {
		out := make([]byte, 1)
		_, err := c.Read(out)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 10; i++ {
		_, err := c.Write([]byte{byte(i + 1)})
		if err != nil {
			t.Fatal(err)
		}
	}
	errs := c.Validate()
	if errs != nil {
		for _, err := range errs {
			t.Error(err)
		}
	}
}
