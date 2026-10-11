//go:build panic_demo

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// panicHandler 有一个真实的 bug：路径没带 id 时 parts[1] 越界，handler 会 panic。
func panicHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	_, _ = w.Write([]byte("id=" + parts[1]))
}

// TestOldServerSwallowsPanic：老写法 NewServer，handler panic 只进 stderr，
// 客户端只看到 EOF，测试"期望出错"照样通过 → PASS。
//
//	go test -tags panic_demo -run TestOldServerSwallowsPanic -v .
func TestOldServerSwallowsPanic(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(panicHandler))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users")
	if err == nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
	if err == nil {
		t.Fatal("expected an error for missing id")
	}
	t.Logf("client err = %v (looks like a normal failure; the panic only went to stderr)", err)
}

// TestNewTestServerCatchesPanic：新写法 NewTestServer，同一个 handler panic
// 被 httptest 报成 t.Errorf → FAIL。
//
//	go test -tags panic_demo -run TestNewTestServerCatchesPanic -v .
func TestNewTestServerCatchesPanic(t *testing.T) {
	ts := httptest.NewTestServer(t, http.HandlerFunc(panicHandler))
	client := ts.Client()

	resp, err := client.Get("http://example.com/users")
	if err == nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
	if err == nil {
		t.Fatal("expected an error for missing id")
	}
	t.Logf("client err = %v (the server-side panic is also reported as t.Errorf)", err)
}
