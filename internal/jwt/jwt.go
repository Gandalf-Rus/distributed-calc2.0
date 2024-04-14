package jwt

import (
	"fmt"
	"time"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

const hmacSampleSecret = "ultra_secret_signature"

func CreateToken(userID int64) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"nbf":    now.Unix(),
		"exp":    now.Add(time.Duration(config.Cfg.JwtTokenTimeout)).Unix(),
		"iat":    now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckTokenAndGetUserID(tokenString string) (int64, error) {
	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		return -1, err
	}

	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		return claims["userID"].(int64), nil
	}

	return -1, err
}
