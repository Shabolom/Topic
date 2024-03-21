package tools

import (
	"Arkadiy_Servis_authorization/config"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

func ParsTokenClaims(strToken string) (Claims, error) {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(strToken, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Env.SecretKey), nil
	})

	if err != nil {
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, errors.New("valid error")
	}

	return claims, nil
}
