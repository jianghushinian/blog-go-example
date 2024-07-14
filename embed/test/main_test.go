package main

import (
	_ "embed"
	"testing"
)

//go:embed testdata/test.txt
var testF string

func TestEmbed(t *testing.T) {
	t.Log(testF)
}
