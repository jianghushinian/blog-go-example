package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jianghushinian/blog-go-example/test/mysql/store"
	"github.com/jianghushinian/blog-go-example/test/mysql/store/mocks"
)

func TestUserHandler_CreateUser_by_mock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := mocks.NewMockUserStore(ctrl)
	mockUserStore.EXPECT().Create(&store.User{
		Name: "user1",
	}).Return(nil)

	handler := &UserHandler{store: mockUserStore}
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	body := `{"name": "user1"}`
	reader := strings.NewReader(body)
	req := httptest.NewRequest("POST", "/users", reader)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "", w.Body.String())
}

func TestUserHandler_GetUser_by_mock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := mocks.NewMockUserStore(ctrl)
	mockUserStore.EXPECT().Get(2).Return(&store.User{
		ID:   2,
		Name: "user2",
	}, nil)

	handler := &UserHandler{store: mockUserStore}
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":2,"name":"user2"}`, w.Body.String())
}
