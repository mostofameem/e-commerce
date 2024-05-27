package auth

import (
	"ecommerce/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(usr models.User) (string, string, error) {
	// Create access token
	accessTokenClaims := jwt.MapClaims{
		"Id":    usr.Id,
		"Name":  usr.Name,
		"Email": usr.Email,
		"exp":   time.Now().Add(1 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	// Create refresh token
	refreshTokenClaims := jwt.MapClaims{
		"Id":  usr.Id,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // Refresh token valid for 7 days
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

var secretKey = []byte("M4q1t8i7eK2oQp5vF0u9Xs6BvG3hT1rD")

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
}
func GenerateAccessTokenFromClaims(claims jwt.Claims) (string, error) {
	// Create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}
