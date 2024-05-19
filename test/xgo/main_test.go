package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xhd2015/xgo/runtime/mock"
	"gorm.io/gorm"
)

func TestUserHandler_CreateUser(t *testing.T) {
	mysqlDB := &gorm.DB{}
	handler := NewUserHandler(mysqlDB)
	router := setupRouter(handler)

	// 为 mysqlDB 打上猴子补丁，替换其 Create 方法
	mock.Patch(mysqlDB.Create, func(value interface{}) (tx *gorm.DB) {
		expected := &User{
			Name: "user1",
		}
		actual := value.(*User)
		assert.Equal(t, expected, actual)
		return mysqlDB
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/users", strings.NewReader(`{"name": "user1"}`))
	router.ServeHTTP(w, req)

	// 断言成功响应
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "", w.Body.String())
}

func TestUserHandler_GetUser(t *testing.T) {
	mysqlDB := &gorm.DB{}
	handler := NewUserHandler(mysqlDB)
	router := setupRouter(handler)

	// 为 mysqlDB 打上猴子补丁，替换其 First 方法
	mock.Patch(mysqlDB.First, func(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
		assert.Equal(t, dest, &User{})
		assert.Equal(t, len(conds), 1)
		assert.Equal(t, conds[0], 1)

		u := dest.(*User)
		u.ID = 1
		u.Name = "user1"
		return mysqlDB
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":1,"name":"user1"}`, w.Body.String())
}

func TestDemo(t *testing.T) {
	t.Log("---------- TestDemo ----------")

	// 测试 trace 功能
	// atoi, err := strconv.Atoi("a")
	// assert.NoError(t, err)
	// assert.Equal(t, atoi, 65)
}
