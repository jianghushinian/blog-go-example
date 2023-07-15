package main

import (
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func Test_monitor_by_gock(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("http://localhost:8080").
		Post("/webhook").
		Reply(200).
		JSON(map[string]interface{}{
			"StatusCode":    0,
			"StatusMessage": "success",
			"Code":          0,
			"Data":          make(map[string]interface{}),
			"Msg":           "success",
		})

	_ = os.Setenv("WEBHOOK", "http://localhost:8080/webhook")
	got, err := monitor(30000000)
	assert.NoError(t, err)
	assert.Equal(t, &Result{
		StatusCode:    0,
		StatusMessage: "success",
		Code:          0,
		Data:          make(map[string]interface{}),
		Msg:           "success",
	}, got)

	assert.True(t, gock.IsDone())
}
