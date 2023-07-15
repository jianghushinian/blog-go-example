package main

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserHandler(t *testing.T) {
	cleanup := setupTestUser()
	defer cleanup()

	w := httptest.NewRecorder()

	body := strings.NewReader(`{"name": "user2"}`)
	req := httptest.NewRequest("POST", "/users", body)

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "", w.Body.String())

	assert.Equal(t, 2, len(users))
	u2, _ := json.Marshal(users[1])
	assert.Equal(t, `{"id":2,"name":"user2"}`, string(u2))
}

func TestGetUserHandler(t *testing.T) {
	cleanup := setupTestUser()
	defer cleanup()

	type want struct {
		code int
		body string
	}
	tests := []struct {
		name string
		args int
		want want
	}{
		{
			name: "get test-user1",
			args: 1,
			want: want{
				code: 200,
				body: `{"id":1,"name":"test-user1"}`,
			},
		},
		{
			name: "get user not found",
			args: 2,
			want: want{
				code: 404,
				body: `{"msg":"notfound"}`,
			},
		},
	}

	router := setupRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d", tt.args), nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}

func setupTestUser() func() {
	defaultUsers := users
	users = []User{
		{ID: 1, Name: "test-user1"},
	}
	return func() {
		users = defaultUsers
	}
}
