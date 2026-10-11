package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserClientA(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		response string
		status   int
		wantErr  bool
	}{
		{
			name:     "success",
			user:     &User{ID: 1, Name: "江湖十年"},
			response: `{"id":1,"name":"江湖十年"}`,
			status:   200,
			wantErr:  false,
		},
		{
			name:     "notfound",
			user:     &User{ID: 2},
			response: "",
			status:   404,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.status)
					if tt.response != "" {
						w.Write([]byte(tt.response))
					}
				}),
			)
			defer ts.Close()

			t.Logf("listener=%v addr=%s", ts.Listener != nil, ts.Listener.Addr()) // true, 127.0.0.1:random port
			if os.Getenv("HOLD") != "" {
				time.Sleep(30 * time.Second) // HOLD=1: hold the process to observe the real port in another terminal
			}

			client := NewUserClient(ts.URL, nil)
			user, err := client.GetUser(tt.user.ID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.user, user)
			}
		})
	}
}

func TestUserClientB(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		response string
		status   int
		wantErr  bool
	}{
		{
			name:     "success",
			user:     &User{ID: 1, Name: "江湖十年"},
			response: `{"id":1,"name":"江湖十年"}`,
			status:   200,
			wantErr:  false,
		},
		{
			name:     "notfound",
			user:     &User{ID: 2},
			response: "",
			status:   404,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewTestServer(
				t,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.status)
					if tt.response != "" {
						w.Write([]byte(tt.response))
					}
				}),
			)

			t.Logf("listener=%v url=%q (before Client call)", ts.Listener != nil, ts.URL)
			if ts.Listener != nil {
				t.Fatal("NewTestServer should not create a real listener")
			}

			client := ts.Client()
			t.Logf("url=%q (after Client call, in-memory network address)", ts.URL)

			uc := NewUserClient(ts.URL, client)
			user, err := uc.GetUser(tt.user.ID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.user, user)
			}
		})
	}
}
