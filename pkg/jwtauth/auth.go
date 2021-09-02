package jwtauth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTAuth struct {
	secretKey string
}

type Claims struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

func NewJWTAuth() *JWTAuth {
	return &JWTAuth{
		secretKey: "so-secr3t!!",
	}
}

func (j *JWTAuth) Generate(claims *Claims, expiry time.Duration) (string, error) {
	claims.ExpiresAt = time.Now().Add(expiry).Unix()
	claims.IssuedAt = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTAuth) Verify(accessToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(j.secretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
