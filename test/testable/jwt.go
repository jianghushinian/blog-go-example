package main

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 依赖注入

type JWT struct {
	privateKey *rsa.PrivateKey
	issuer     string
	// nowFunc is used to mock time in tests
	nowFunc func() time.Time
}

func NewJWT(issuer string, privateKey *rsa.PrivateKey) *JWT {
	return &JWT{
		privateKey: privateKey,
		issuer:     issuer,
		nowFunc:    time.Now,
	}
}

func (j *JWT) GenerateToken(userId string, expire time.Duration) (string, error) {
	nowSec := j.nowFunc().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		// map 会对其进行重新排序，排序结果影响签名结果，签名结果验证网址：https://jwt.io/
		"issuer":    j.issuer,
		"issuedAt":  nowSec,
		"expiresAt": nowSec + int64(expire.Seconds()),
		"subject":   userId,
	})

	return token.SignedString(j.privateKey)
}

func GenerateJWT(issuer string, userId string, nowFunc func() time.Time, expire time.Duration, privateKey *rsa.PrivateKey) (string, error) {
	nowSec := nowFunc().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"expiresAt": nowSec + int64(expire.Seconds()),
		"issuedAt":  nowSec,
		"issuer":    issuer,
		"subject":   userId,
	})

	return token.SignedString(privateKey)
}
