package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Tokens struct {
	Access                string
	TokenExpiresIn        string
	Refresh               string
	RefreshTokenExpiresIn string
}

func GenerateNewTokens(id string) (*Tokens, error) {
	accessToken, tokenExpiresIn, err := generateNewAccessToken(id)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenExpiresIn, err := generateNewRefreshToken()
	if err != nil {
		return nil, err
	}

	return &Tokens{
		Access:                accessToken,
		TokenExpiresIn:        tokenExpiresIn,
		Refresh:               refreshToken,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
	}, nil
}

func generateNewAccessToken(id string) (string, string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")

	minutesCount, _ := strconv.Atoi(os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))
	expiration := time.Now().Add(time.Minute * time.Duration(minutesCount))

	claims := jwt.MapClaims{}

	claims["id"] = id
	claims["expires"] = expiration.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return t, expiration.Format(time.RFC3339), nil
}

func generateNewRefreshToken() (string, string, error) {
	hash := sha256.New()

	refresh := os.Getenv("JWT_REFRESH_KEY") + time.Now().String()

	_, err := hash.Write([]byte(refresh))
	if err != nil {
		return "", "", err
	}

	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))
	expiration := time.Now().Add(time.Hour * time.Duration(hoursCount))
	expireTime := fmt.Sprint(expiration.Unix())

	t := hex.EncodeToString(hash.Sum(nil)) + "." + expireTime

	return t, expiration.Format(time.RFC3339), nil
}

func ParseRefreshToken(refreshToken string) (int64, error) {
	return strconv.ParseInt(strings.Split(refreshToken, ".")[1], 0, 64)
}
