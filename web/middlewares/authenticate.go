package middlewares

import (
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/web/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	accessTokenString, err := accessToken.SignedString([]byte(config.GetConfig().JwtSecret))
	if err != nil {
		return "", "", err
	}

	// Create refresh token
	refreshTokenClaims := jwt.MapClaims{
		"Id":  usr.Id,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // Refresh token valid for 7 days
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.GetConfig().JwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JwtSecret), nil
	})
}
func GenerateAccessTokenFromClaims(claims jwt.Claims) (string, error) {
	// Create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(config.GetConfig().JwtSecret))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendError(w, http.StatusForbidden, fmt.Errorf("authorization header is missing"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, _ := ParseToken(tokenString)
		if token.Valid {
			next.ServeHTTP(w, r) // Token is valid, continue with the request
			return
		}
		// Token has expired, check for refresh token
		refreshHeader := r.Header.Get("Refresh-Token")
		if refreshHeader == "" {
			utils.SendError(w, http.StatusUnauthorized, fmt.Errorf("refresh token missing"))
			return
		}

		// Validate and parse the refresh token
		refreshString := strings.TrimPrefix(refreshHeader, "Bearer ")
		refreshToken, err := ParseToken(refreshString)
		if err != nil {
			utils.SendError(w, http.StatusUnauthorized, fmt.Errorf("invalid refresh token: %v", err))
			return
		}

		if refreshToken.Valid {
			// Generate new access token
			claims := token.Claims.(jwt.MapClaims)
			claims["exp"] = time.Now().Add(1 * time.Minute).Unix()
			newToken, err := GenerateAccessTokenFromClaims(claims)
			if err != nil {
				utils.SendError(w, http.StatusInternalServerError, fmt.Errorf("error generating new token: %v", err))
				return
			}
			log.Println(newToken)
			// Continue with the request
			next.ServeHTTP(w, r)
			return
		}

		utils.SendError(w, http.StatusUnauthorized, fmt.Errorf("refresh token invalid"))
	})
}

func GetIdFromHeader(r string) (string, error) {

	authHeader := r
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := ParseToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("error Parsing")
	}

	// Extract user ID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	userID, _ := claims["Id"].(string)
	return userID, nil
}
