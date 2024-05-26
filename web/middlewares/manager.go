package middlewares

import (
	"ecommerce/auth"
	"ecommerce/web/utils"
	"fmt"
	"net/http"
	"strings"
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
		token, err := auth.ParseToken(tokenString)
		if err != nil || !token.Valid {
			utils.SendError(w, http.StatusForbidden, fmt.Errorf("invalid token"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
