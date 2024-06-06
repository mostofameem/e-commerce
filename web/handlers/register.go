package handlers

import (
	"context"
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
type Code struct {
	Code string `json:"code" validate:"required"`
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

	utils.SendBothData(w, user, "Success! OTP sent to your email")
}

func Verification(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/register/"):]
	email := db.GetUserTypeRepo().GetEmail(id)

	var usercode Code
	err := json.NewDecoder(r.Body).Decode(&usercode)
	if err != nil {
		slog.Error("Failed to get user data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": usercode,
		}))
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}
	ok := VerificationRedis(email, usercode.Code)
	if !ok {
		slog.Error("Invalid Code", logger.Extra(map[string]any{
			"payload": usercode,
		}))
		utils.SendError(w, http.StatusBadRequest, err)
	}
	db.GetUserTypeRepo().Verified(id)
	utils.SendData(w, "Account Verified Successful ")
}

func VerificationRedis(email, usercode string) bool {

	code := db.GetRedis().Get(context.Background(), email)
	return code.Val() == usercode
}
