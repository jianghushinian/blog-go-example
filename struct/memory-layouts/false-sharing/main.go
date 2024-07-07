package main

import "golang.org/x/sys/cpu"

var S struct {
	a string
	_ cpu.CacheLinePad
}
