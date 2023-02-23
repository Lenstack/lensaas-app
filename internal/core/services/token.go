package services

import (
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"math/rand"
	"time"
)

type ITokenService interface {
	GenerateToken(userId string, expiration time.Duration) (string, error)
	ValidateToken(token string) (string, error)
	NewRefreshToken() (string, error)
}

type TokenService struct {
	secret         string
	ExpirationTime time.Duration
}

func NewTokenService(secret string, expiration string) *TokenService {
	expirationTime, err := time.ParseDuration(expiration)
	if err != nil {
		panic(err)
	}
	return &TokenService{secret, expirationTime}
}

func (ts *TokenService) GenerateToken(userId string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(expiration).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(ts.secret))
}

func (ts *TokenService) ValidateToken(token string) (string, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	parsedToken, err := parser.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(ts.secret), nil
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

func (ts *TokenService) NewRefreshToken() (string, error) {
	token := make([]byte, 64)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(token), nil
}
