package utils

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type TokenMetadata struct {
	UserId  string
	Expires int64
}

func ExtractTokenMetadata(tokenStr string) (*TokenMetadata, error) {
	token, err := verifyToken(tokenStr)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID := claims["id"].(string)

		expires := int64(claims["expires"].(float64))

		return &TokenMetadata{
			UserId:  userID,
			Expires: expires,
		}, nil
	}

	return nil, err
}

func ExtractToken(bearToken string) string {
	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func verifyToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
