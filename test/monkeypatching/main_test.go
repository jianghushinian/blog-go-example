package main

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserHandler_CreateUser(t *testing.T) {
	mysqlDB := &gorm.DB{}
	handler := NewUserHandler(mysqlDB)
	router := setupRouter(handler)

	// 为 mysqlDB 打上猴子补丁，替换其 Create 方法
	patches := gomonkey.ApplyMethod(reflect.TypeOf(mysqlDB), "Create",
		func(in *gorm.DB, value interface{}) (tx *gorm.DB) {
			expected := &User{
				Name: "user1",
			}
			actual := value.(*User)
			assert.Equal(t, expected, actual)
			return in
		})
	// 测试执行完成后将猴子补丁复原
	defer patches.Reset()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/users", strings.NewReader(`{"name": "user1"}`))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "", w.Body.String())
}

func TestUserHandler_GetUser(t *testing.T) {
	mysqlDB := &gorm.DB{}
	handler := NewUserHandler(mysqlDB)
	router := setupRouter(handler)

	patches := gomonkey.ApplyMethod(reflect.TypeOf(mysqlDB), "First",
		func(in *gorm.DB, dest interface{}, conds ...interface{}) (tx *gorm.DB) {
			assert.Equal(t, dest, &User{})
			assert.Equal(t, len(conds), 1)
			assert.Equal(t, conds[0], 1)

			u := dest.(*User)
			u.ID = 1
			u.Name = "user1"
			return in
		})
	defer patches.Reset()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":1,"name":"user1"}`, w.Body.String())
}
