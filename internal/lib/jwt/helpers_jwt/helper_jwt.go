package helpers_jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenString string, secret string) (jwt.MapClaims, error) {

	ParsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("failed to parse token")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := ParsedToken.Claims.(jwt.MapClaims); ok && ParsedToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token claims")
}
