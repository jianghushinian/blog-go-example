//go:build startorder_demo

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Start() 会把内存服务器切成 loopback 真实监听，此时 URL 变成 127.0.0.1:port，
// 任意 client 都能访问。顺序是"先 Start、后 Client"。
//
//	go test -tags startorder_demo -run TestStartBeforeClient -v .
func TestStartBeforeClient(t *testing.T) {
	ts := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	ts.Start()

	client := ts.Client()
	t.Logf("after Start: url=%q listener=%v", ts.URL, ts.Listener != nil)

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("loopback server should be reachable: %v", err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

// 反过来先 Client() 再 Start() 会 panic "Server already started"。
//
//	go test -tags startorder_demo -run TestClientThenStartPanics -v .
func TestClientThenStartPanics(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered panic: %v", r)
			return
		}
		t.Fatal("expected a panic when calling Start after Client")
	}()

	ts := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	_ = ts.Client()
	ts.Start()
}
