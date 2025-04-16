package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwt_secret_key = "019639a9-e3a3-7c2d-92f1-23e62c3ecdc0"

type claims struct {
	ClientID string `json:"client_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(clientID string) (string, error) {
	claims := claims{
		ClientID: clientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwt_secret_key))
}

func IsJwTValid(token string) (bool, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(jwt_secret_key), nil
	})
	if err != nil {
		return false, err
	}

	if _, ok := jwtToken.Claims.(*claims); ok && jwtToken.Valid {
		return true, nil
	}

	return false, nil
}

func DecodeBasicAuth(authHeader string) (username, password string, err error) {
	prefix := "Basic "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", "", errors.New("authorization header does not have Basic prefix")
	}

	encoded := authHeader[len(prefix):]
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode base64 credentials: %w", err)
	}
	decoded := string(decodedBytes)

	parts := strings.Split(decoded, ":")
	if len(parts) != 2 {
		return "", "", errors.New("invalid basic auth format, missing colon separator")
	}

	return parts[0], parts[1], nil
}
