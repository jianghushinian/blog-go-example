package main

import (
	"crypto/rand"
	"encoding/base64"
)

// Untestable
/*
	func GenerateToken(n int) (string, error) {
		token := make([]byte, n)
		_, err := rand.Read(token)
		if err != nil {
			return "", err
		}
		return base64.URLEncoding.EncodeToString(token)[:n], nil
	}
*/

var GenerateToken = func(n int) (string, error) {
	token := make([]byte, n)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token)[:n], nil
}

type User struct {
	ID     int
	Name   string
	Mobile string
}

func Login(u User) (string, error) {
	// ...
	token, err := GenerateToken(32)
	if err != nil {
		// ...
	}
	// ...
	return token, nil
}
