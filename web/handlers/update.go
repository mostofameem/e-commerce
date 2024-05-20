package handlers

import (
	"ecommerce/auth"
	"ecommerce/web/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func Update(w http.ResponseWriter, r *http.Request) {
	authHandler := auth.AuthenticateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusForbidden)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.ParseToken(tokenString)
		if err != nil {
			utils.SendError(w, http.StatusForbidden, err, "Error parsing")
		}
		// Extract user ID from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusForbidden)
			return
		}
		userID, err := strconv.Atoi(fmt.Sprintf("%.0f", claims["Id"].(float64)))
		if err != nil {
			http.Error(w, "Invalid user ID in token", http.StatusForbidden)
			return
		}
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		Users[userID].Username = user.Username
		Users[userID].Password = user.Password

		utils.SendData(w, Users[userID])

	}))
	authHandler.ServeHTTP(w, r)
}
