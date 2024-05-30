package handlers

import (
	"ecommerce/db"
	"ecommerce/logger"
	"ecommerce/web/utils"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

type NewUser struct {
	Name     string `json:"name" validate:"required,min=3,max=20,alpha"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user NewUser
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

	err = db.GetUserTypeRepo().Create(user.Name, user.Email, user.Password)
	if err != nil {
		log.Println(err)
		slog.Error("Failed to insert user db ", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": user,
		}))
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	utils.SendBothData(w, user, "Register successful ")
}
