package main

import "errors"

func Positive(n int) (bool, bool) {
	if n == 0 {
		return false, false
	}
	return n > -1, true
}

func Positive1(n int) (bool, error) {
	if n == 0 {
		return false, errors.New("undefined")
	}
	return n > -1, nil
}

// NOTE: 不要返回指针

func Positive2(n int) *bool {
	if n == 0 {
		return nil
	}
	r := n > -1
	return &r
}
