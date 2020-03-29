package security

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"little-diary-measurement-service/src/errors"
)

type JwtTokenReader struct {
	PublicKey string
}

func (v *JwtTokenReader) ReadUserUuid(tokenString string) (string, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		publicKeyBytes := []byte(v.PublicKey)
		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
		if err != nil {
			return nil, err
		}
		return verifyKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userClaims, ok := claims["user_claims"].(map[string]interface{})
		if !ok {
			return "", &errors.ForbiddenError{S: "token invalid"}
		}
		userUuid, ok := userClaims["uuid"].(string)
		if !ok {
			return "", &errors.ForbiddenError{S: "token invalid"}
		}
		return userUuid, nil
	} else {
		return "", &errors.ForbiddenError{S: "token invalid"}
	}
}
