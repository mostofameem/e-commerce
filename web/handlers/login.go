package handlers

import (
	"ecommerce/db"
	"ecommerce/logger"
	"ecommerce/web/middlewares"
	"ecommerce/web/utils"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("Failed to get user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": user,
		}))
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = utils.Validate(user)
	if err != nil {
		slog.Error("Failed to validate user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": user,
		}))
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = db.GetUserTypeRepo().Login(user.Email, user.Password)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, fmt.Errorf("wrong username / password "))
		return
	}

	userinfo, err := db.GetUserTypeRepo().GetUser(user.Email)

	if err != nil {
		slog.Error(
			"Failed to get user Info",
			logger.Extra(map[string]any{
				"error": err.Error(),
			}),
		)
		return
	}

	accessToken, refreshToken, err := middlewares.GenerateToken(userinfo)
	if err != nil {
		slog.Error("Error Generating Tokens", logger.Extra(map[string]any{
			"error":        err.Error(),
			"payload":      user,
			"Accesstoken":  accessToken,
			"RefreshToken": refreshToken,
		}))
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}
	utils.SendBothData(w, refreshToken, accessToken)
}
