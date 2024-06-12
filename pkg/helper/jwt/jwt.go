package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	ID   string
	Role string
}

func GenerateToken(user UserClaims, secret string, duration time.Duration) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iss": "@j11er1",
		"aud": user.Role,
		"exp": time.Now().Add(time.Second * time.Duration(duration.Seconds())).Unix(),
		"iat": time.Now().Unix(),
	})

	s, err := claims.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return s, nil
}

func VerifyToken(token string, secret string) (jwt.Claims, error) {

	t, err := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return t.Claims, nil
}
