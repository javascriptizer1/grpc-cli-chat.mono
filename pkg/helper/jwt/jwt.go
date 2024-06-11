package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain"
)

// TODO: pkg is independent on internal

func GenerateToken(user domain.User, secret string, duration time.Duration) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iss": "@j11er1",
		"aud": strconv.Itoa(int(user.Role)),
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
