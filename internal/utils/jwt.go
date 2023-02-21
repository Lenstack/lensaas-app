package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type IJwt interface {
	GenerateToken(userId string, expiration time.Duration) (string, error)
	ValidateToken(token string) (string, error)
}

type Jwt struct {
	secret         string
	ExpirationTime time.Duration
}

func NewJwt(secret string, expiration string) *Jwt {
	expirationTime, err := time.ParseDuration(expiration)
	if err != nil {
		panic(err)
	}
	return &Jwt{secret, expirationTime}
}

func (j *Jwt) GenerateToken(userId string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(expiration).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.secret))
}

func (j *Jwt) ValidateToken(token string) (string, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	parsedToken, err := parser.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}
	return claims["id"].(string), nil
}
