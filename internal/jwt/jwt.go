package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Claims содержит данные из JWT токена
type Claims struct {
	UserID int64
	Role   string
}

// ParseToken парсит JWT токен и проверяет его подпись и срок действия.
// Не обращается к базе данных, только валидирует токен.
func ParseToken(tokenStr string, secret []byte) (*Claims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	sub, ok := claims["sub"]
	if !ok {
		return nil, errors.New("token subject not found")
	}

	var userID int64
	switch v := sub.(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	default:
		return nil, errors.New("invalid token subject type")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("token role not found")
	}

	return &Claims{
		UserID: userID,
		Role:   role,
	}, nil
}
