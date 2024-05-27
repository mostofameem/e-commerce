package middlewares

import (
	"ecommerce/auth"
	"ecommerce/web/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Middleware func(http.Handler) http.Handler

type Manager struct {
	globalMiddlewares []Middleware
}

func NewManager() *Manager {
	return &Manager{
		globalMiddlewares: make([]Middleware, 0),
	}
}

func (m Manager) Use(middlewares ...Middleware) Manager {
	m.globalMiddlewares = append(m.globalMiddlewares, middlewares...)
	return m
}

func (m *Manager) With(handler http.Handler, middlewares ...Middleware) http.Handler {
	var h http.Handler
	h = handler

	for _, m := range middlewares {
		h = m(h)
	}

	for _, m := range m.globalMiddlewares {
		h = m(h)
	}

	return h
}

func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendError(w, http.StatusForbidden, fmt.Errorf("authorization header is missing"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, _ := auth.ParseToken(tokenString)
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
		refreshToken, err := auth.ParseToken(refreshString)
		if err != nil {
			utils.SendError(w, http.StatusUnauthorized, fmt.Errorf("invalid refresh token: %v", err))
			return
		}

		if refreshToken.Valid {
			// Generate new access token
			claims := token.Claims.(jwt.MapClaims)
			claims["exp"] = time.Now().Add(1 * time.Minute).Unix()
			newToken, err := auth.GenerateAccessTokenFromClaims(claims)
			if err != nil {
				utils.SendError(w, http.StatusInternalServerError, fmt.Errorf("error generating new token: %v", err))
				return
			}
			//w.Header().Set("Authorization", "Bearer "+newToken)
			utils.SendData(w, newToken)
			// Continue with the request
			next.ServeHTTP(w, r)
			return
		}

		utils.SendError(w, http.StatusUnauthorized, fmt.Errorf("refresh token invalid"))
	})
}
