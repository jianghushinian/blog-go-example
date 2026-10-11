package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	r := setupRouter()

	body := bytes.NewBufferString(`{"name": "江湖十年"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var u User
	err := json.Unmarshal(w.Body.Bytes(), &u)
	assert.NoError(t, err)
	assert.Equal(t, "江湖十年", u.Name)
	assert.Equal(t, 1, u.ID)
}
