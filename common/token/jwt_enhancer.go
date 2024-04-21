package token

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidMethodHS256 = errors.New("invalid method HS256")
)

type JWTEnhancer struct {
	secret []byte
}

func NewJWTEnhancer(secret []byte) Enhancer {
	return &JWTEnhancer{secret: secret}
}

func (enhancer JWTEnhancer) Sign(claims *Claims) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenClaims.SignedString(enhancer.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (enhancer JWTEnhancer) Verify(value string) (*Claims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidMethodHS256
		}
		return enhancer.secret, nil
	}

	tokenClaims, err := jwt.ParseWithClaims(value, &Claims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims := tokenClaims.Claims.(*Claims)

	return claims, nil
}
