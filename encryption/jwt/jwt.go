package jwt

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Comment
func JwtEncode(secret string, claims jwt.Claims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return token, nil
}

// Comment
func JwtDecode[C interface{}](secret string, token string, claims C) (*C, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	data, err := json.Marshal(jwtToken.Claims)

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	return &claims, nil
}
