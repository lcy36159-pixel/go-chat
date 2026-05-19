package jwt

import (
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
)

const defaultTokenTTLHours = 24
const minSecretLength = 32

var (
	loadSecretOnce sync.Once
	secretBytes    []byte
	secretErr      error
	loadTTLOnce    sync.Once
	ttlDuration    time.Duration
)

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
			ExpiresAt: jwtgo.NewNumericDate(time.Now().Add(tokenTTL())),
		},
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}
	return token.SignedString(secret)
}

func ParseToken(tokenStr string) (uint, error) {
	if tokenStr == "" {
		return 0, errors.New("token is required")
	}

	claims := &Claims{}
	token, err := jwtgo.ParseWithClaims(tokenStr, claims, func(token *jwtgo.Token) (interface{}, error) {
		method, ok := token.Method.(*jwtgo.SigningMethodHMAC)
		if !ok || method.Alg() != jwtgo.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret()
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

func jwtSecret() ([]byte, error) {
	loadSecretOnce.Do(func() {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secretErr = errors.New("JWT_SECRET is required")
			return
		}
		if len(secret) < minSecretLength {
			secretErr = errors.New("JWT_SECRET must be at least 32 characters")
			return
		}
		secretBytes = []byte(secret)
	})
	return secretBytes, secretErr
}

func tokenTTL() time.Duration {
	loadTTLOnce.Do(func() {
		ttlDuration = defaultTokenTTLHours * time.Hour
		ttlHoursStr := os.Getenv("JWT_TTL_HOURS")
		if ttlHoursStr == "" {
			return
		}
		ttlHours, err := strconv.Atoi(ttlHoursStr)
		if err != nil || ttlHours <= 0 {
			return
		}
		ttlDuration = time.Duration(ttlHours) * time.Hour
	})

	return ttlDuration
}
