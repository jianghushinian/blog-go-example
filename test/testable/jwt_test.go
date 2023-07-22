package main

import (
	_ "embed"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// 生成测试用的密钥对
// openssl genrsa -out private.pem 2048
// openssl rsa -pubout -in private.pem -out public.pem

//go:embed  testdata/private.pem
var privateKey string

func TestJWT_GenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	assert.NoError(t, err)

	j := NewJWT("jianghushinian", key)
	j.nowFunc = func() time.Time {
		return time.Unix(1689815972, 0)
	}

	actual, err := j.GenerateToken("1234", 2*time.Hour)
	assert.NoError(t, err)

	expected := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOjE2ODk4MjMxNzIsImlzc3VlZEF0IjoxNjg5ODE1OTcyLCJpc3N1ZXIiOiJqaWFuZ2h1c2hpbmlhbiIsInN1YmplY3QiOiIxMjM0In0.NmCDxFaBfAPPgWQ0zVMl8ON1UQMeIVNgFCn1vtbppsunb-VrOMCdnJlguvPnNc6fMD9EkzMYM3Ux8zFnTiICDMRX23UlhAo2Zb3DorThdrBcNWHMUd26DBNI9n_oUY5B6NPqtrutvqCex9lQH0vUYOt2O5dOyZ-H9cVNY1r3fJHNkYuNWxmoZRfka5o1oSWvUw8hBJfgjANOzZ5ACIi0q5hnou5hQ8VljjFsP4zj2a2lU6w5Db8_rOA04BxilkfurdExcPeaAVCtA-Km0zNwL3gGwJB21gwyb4MRHsEf-ra-4-V7O5_JGiSOQgfkNB63RoASljRXpD6q-gakm0e0fA"
	assert.Equal(t, expected, actual)
}

func TestGenerateJWT(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	assert.NoError(t, err)

	nowFunc := func() time.Time {
		return time.Unix(1689815972, 0)
	}

	actual, err := GenerateJWT("jianghushinian", "1234", nowFunc, 2*time.Hour, key)
	assert.NoError(t, err)

	expected := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOjE2ODk4MjMxNzIsImlzc3VlZEF0IjoxNjg5ODE1OTcyLCJpc3N1ZXIiOiJqaWFuZ2h1c2hpbmlhbiIsInN1YmplY3QiOiIxMjM0In0.NmCDxFaBfAPPgWQ0zVMl8ON1UQMeIVNgFCn1vtbppsunb-VrOMCdnJlguvPnNc6fMD9EkzMYM3Ux8zFnTiICDMRX23UlhAo2Zb3DorThdrBcNWHMUd26DBNI9n_oUY5B6NPqtrutvqCex9lQH0vUYOt2O5dOyZ-H9cVNY1r3fJHNkYuNWxmoZRfka5o1oSWvUw8hBJfgjANOzZ5ACIi0q5hnou5hQ8VljjFsP4zj2a2lU6w5Db8_rOA04BxilkfurdExcPeaAVCtA-Km0zNwL3gGwJB21gwyb4MRHsEf-ra-4-V7O5_JGiSOQgfkNB63RoASljRXpD6q-gakm0e0fA"
	assert.Equal(t, expected, actual)
}
