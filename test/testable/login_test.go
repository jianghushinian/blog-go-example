package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	u := User{
		ID:     1,
		Name:   "test1",
		Mobile: "13800001111",
	}
	token, err := Login(u)
	assert.NoError(t, err)
	assert.Equal(t, 32, len(token))
	assert.Equal(t, "jCnuqKnsN5UAM9-LgEGS_COvJWp15RDv", token)
}

func init() {
	GenerateToken = func(n int) (string, error) {
		return "jCnuqKnsN5UAM9-LgEGS_COvJWp15RDv", nil
	}
}
