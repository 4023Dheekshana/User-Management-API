package utils

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "SuperKey"

func GenerateToken(username, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Println("Error parsing token", err)
		return "", errors.New("error parsing token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Invalid token claims", err)
		return "", errors.New("invalid token claims")
	}
	username, ok := claims["username"].(string)
	if !ok {
		log.Println("Username not found in token claims")
		return "", errors.New("invalid token user claims")
	}
	return username, nil

}
