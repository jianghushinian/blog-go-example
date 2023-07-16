package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jianghushinian/blog-go-example/test/mysql/store"
)

type fakeUserStore struct{}

func (f *fakeUserStore) Create(user *store.User) error {
	return nil
}

func (f *fakeUserStore) Get(id int) (*store.User, error) {
	return &store.User{ID: id, Name: "test"}, nil
}

func TestUserHandler_CreateUser_by_fake(t *testing.T) {
	handler := &UserHandler{store: &fakeUserStore{}}
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	body := `{"name": "user2"}`
	reader := strings.NewReader(body)
	req := httptest.NewRequest("POST", "/users", reader)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "", w.Body.String())
}

func TestUserHandler_GetUser_by_fake(t *testing.T) {
	handler := &UserHandler{store: &fakeUserStore{}}
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":1,"name":"test"}`, w.Body.String())
}
