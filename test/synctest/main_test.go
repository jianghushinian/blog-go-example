package main

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

func TestAfterFunc1(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	calledCh := make(chan struct{}) // closed when AfterFunc is called
	context.AfterFunc(ctx, func() {
		close(calledCh)
	})

	// TODO: Assert that the AfterFunc has not been called.

	cancel()

	// TODO: Assert that the AfterFunc has been called.
}

func TestAfterFunc2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	calledCh := make(chan struct{}) // closed when AfterFunc is called
	context.AfterFunc(ctx, func() {
		close(calledCh)
	})

	// funcCalled reports whether the function was called.
	funcCalled := func() bool {
		select {
		case <-calledCh:
			return true
		case <-time.After(10 * time.Millisecond):
			return false
		}
	}

	if funcCalled() {
		t.Fatalf("AfterFunc function called before context is canceled")
	}

	cancel()

	if !funcCalled() {
		t.Fatalf("AfterFunc function not called after context is canceled")
	}
}

// GOEXPERIMENT=synctest go test . -v -run=^TestAfterFunc$
func TestAfterFunc(t *testing.T) {
	synctest.Run(func() {
		ctx, cancel := context.WithCancel(context.Background())

		funcCalled := false
		context.AfterFunc(ctx, func() {
			funcCalled = true
		})

		synctest.Wait()
		if funcCalled {
			t.Fatalf("AfterFunc function called before context is canceled")
		}

		cancel()

		synctest.Wait()
		if !funcCalled {
			t.Fatalf("AfterFunc function not called after context is canceled")
		}
	})
}

// GOEXPERIMENT=synctest go test . -v -run=^TestWithTimeout$
func TestWithTimeout(t *testing.T) {
	synctest.Run(func() {
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Wait just less than the timeout.
		time.Sleep(timeout - time.Nanosecond)
		synctest.Wait()
		if err := ctx.Err(); err != nil {
			t.Fatalf("before timeout, ctx.Err() = %v; want nil", err)
		}

		// Wait the rest of the way until the timeout.
		time.Sleep(time.Nanosecond)
		synctest.Wait()
		if err := ctx.Err(); err != context.DeadlineExceeded {
			t.Fatalf("after timeout, ctx.Err() = %v; want DeadlineExceeded", err)
		}
	})
}

// GOEXPERIMENT=synctest go test . -v -run=^Test$
func Test(t *testing.T) {
	synctest.Run(func() {
		srvConn, cliConn := net.Pipe()
		defer srvConn.Close()
		defer cliConn.Close()
		tr := &http.Transport{
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				return cliConn, nil
			},
			// Setting a non-zero timeout enables "Expect: 100-continue" handling.
			// Since the following test does not sleep,
			// we will never encounter this timeout,
			// even if the test takes a long time to run on a slow machine.
			ExpectContinueTimeout: 5 * time.Second,
		}

		body := "request body"
		go func() {
			req, _ := http.NewRequest("PUT", "http://test.tld/", strings.NewReader(body))
			req.Header.Set("Expect", "100-continue")
			resp, err := tr.RoundTrip(req)
			if err != nil {
				t.Errorf("RoundTrip: unexpected error %v", err)
			} else {
				resp.Body.Close()
			}
		}()

		req, err := http.ReadRequest(bufio.NewReader(srvConn))
		if err != nil {
			t.Fatalf("ReadRequest: %v", err)
		}

		var gotBody strings.Builder
		go io.Copy(&gotBody, req.Body)
		synctest.Wait()
		if got := gotBody.String(); got != "" {
			t.Fatalf("before sending 100 Continue, unexpectedly read body: %q", got)
		}

		srvConn.Write([]byte("HTTP/1.1 100 Continue\r\n\r\n"))
		synctest.Wait()
		if got := gotBody.String(); got != body {
			t.Fatalf("after sending 100 Continue, read body %q, want %q", got, body)
		}

		srvConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	})
}
