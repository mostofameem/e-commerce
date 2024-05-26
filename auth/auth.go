package auth

import (
	"ecommerce/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(usr models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":    usr.Id,
		"Name":  usr.Name,
		"Email": usr.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString(secretKey)
}

var secretKey = []byte("M4q1t8i7eK2oQp5vF0u9Xs6BvG3hT1rD")

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
}
