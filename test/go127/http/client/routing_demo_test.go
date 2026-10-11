//go:build inmemory_demo

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewTestServer 默认走内存网络，“无视地址、全部路由到测试服务器”这个魔法
// 只存在于 server.Client() 返回的 Transport 里（它在 DialContext 里直接
// 返回内存 listener 的连接）。一旦换成普通 client（比如 http.DefaultClient），
// 请求就会真的去拨 server.URL 里的地址——内存模式下是 example.com，直接打到公网。
//
//	go test -tags inmemory_demo -v .
func TestInMemoryRoutingGotcha(t *testing.T) {
	ts := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "served by test server, host=%s path=%s", r.Host, r.URL.Path)
	}))

	t.Logf("before Client(): url=%q listener=%v", ts.URL, ts.Listener)
	client := ts.Client()
	t.Logf("after  Client(): url=%q listener=%v", ts.URL, ts.Listener)

	// server.Client()：任意 host 都会被路由到内存服务器
	resp, err := client.Get("http://totally.other.host/anything")
	if err != nil {
		t.Fatalf("server.Client() should reach the in-memory server: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	t.Logf("server.Client() -> %q", body)

	// 普通 client：没有这层路由，会真的去拨 server.URL 的地址。
	// 这里用一个 DialContext 拦下来，证明"真的在拨公网"，同时不依赖外网。
	dialed := ""
	plain := &http.Client{Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialed = addr
			return nil, errors.New("blocked by demo")
		},
	}}
	if _, err := plain.Get(ts.URL); err == nil {
		t.Fatal("plain client should not reach the in-memory server")
	}
	t.Logf("plain client dialed %q (real network), not the in-memory listener", dialed)
}
