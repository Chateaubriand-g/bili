package util

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ATokenexpiredHour = 2
	ATokenwaitHour    = 0
	RTokenexpiredHour = 24
	RTokenwaitHour    = 0
)

var (
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenInvalid     = errors.New("token invaild")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrSignatureInvalid = errors.New("signatrue invaild")
)

type Claims struct {
	UserID uint64 `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uint64, secret string) (string, error) {
	expiredTime := time.Now().Add(time.Duration(ATokenexpiredHour) * time.Hour)
	waitTime := time.Now().Add(time.Duration(ATokenwaitHour) * time.Hour)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(waitTime),
			Issuer:    "g",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", nil
	}
	return hex.EncodeToString(b), nil
}
