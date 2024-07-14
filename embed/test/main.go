package main

import "embed"

//go:embed testdata
var testFS embed.FS

func main() {
	_ = testFS
}
