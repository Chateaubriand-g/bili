package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Chateaubriand-g/bili/gateway/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	expiredHour = 2
	waitHour    = 0
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

func GenerateToken(userID uint64, secret string) (string, error) {
	expiredTime := time.Now().Add(time.Duration(expiredHour) * time.Hour)
	waitTime := time.Now().Add(time.Duration(waitHour) * time.Hour)

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

func ParseToken(tokenstring, secret string) (*Claims, error) {
	/*
		jwt.Token{
			Raw string
			Method SigningMethod
			Header map[string]interface{}
			Claim Claims
			Signature string
			Valid bool
		}
	*/
	token, err := jwt.ParseWithClaims(tokenstring, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		} else if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, ErrSignatureInvalid
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims, err := ParseToken(parts[1], cfg.Jwt.Secret)
		if err != nil {
			switch err {
			case ErrTokenExpired:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			case ErrTokenNotValidYet:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token not valid yet"})
			case ErrSignatureInvalid:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token signature error"})
			case ErrTokenInvalid:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalid"})
			}
			return
		}

		userID := strconv.FormatUint(claims.UserID, 10)
		//X-xxx-xx 自定义请求头
		c.Request.Header.Set("X-User-ID", userID)
		c.Next()
	}
}
