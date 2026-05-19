package jwt

import (
	"errors"
	"os"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 24 * time.Hour

type Claims struct {
	UserID uint `json:"user_id"`
	jwtgo.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
	if userID == 0 {
		return "", errors.New("user_id is required")
	}

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(time.Now().Add(tokenTTL)),
		},
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

func ParseToken(tokenStr string) (uint, error) {
	if tokenStr == "" {
		return 0, errors.New("token is required")
	}

	claims := &Claims{}
	token, err := jwtgo.ParseWithClaims(tokenStr, claims, func(token *jwtgo.Token) (interface{}, error) {
		if token.Method != jwtgo.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret(), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}
	if claims.UserID == 0 {
		return 0, errors.New("invalid user_id in token")
	}

	return claims.UserID, nil
}

func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "go-chat-secret"
	}
	return []byte(secret)
}
