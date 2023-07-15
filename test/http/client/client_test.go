package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

var ts *httptest.Server

func Test_monitor(t *testing.T) {
	type args struct {
		pid     int
		webhook string
	}
	tests := []struct {
		name    string
		args    args
		want    *Result
		wantErr error
	}{
		{
			name: "process exited and send feishu success",
			args: args{
				pid:     10000000,
				webhook: ts.URL + "/success",
			},
			want: &Result{
				StatusCode:    0,
				StatusMessage: "success",
				Code:          0,
				Data:          make(map[string]interface{}),
				Msg:           "success",
			},
		},
		{
			name: "process exited and send feishu error",
			args: args{
				pid:     20000000,
				webhook: ts.URL + "/error",
			},
			wantErr: errors.New("code: 19001, error: param invalid: incoming webhook access token invalid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("WEBHOOK", tt.args.webhook)

			got, err := monitor(tt.args.pid)
			if err != nil {
				if tt.wantErr == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("monitor() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("monitor() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func newTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.RequestURI {
		case "/success":
			_, _ = fmt.Fprintf(w, `{"StatusCode":0,"StatusMessage":"success","code":0,"data":{},"msg":"success"}`)
		case "/error":
			_, _ = fmt.Fprintf(w, `{"code":19001,"data":{},"msg":"param invalid: incoming webhook access token invalid"}`)
		}
	}))
}

func TestMain(m *testing.M) {
	ts = newTestServer()
	m.Run()
	ts.Close()
}
