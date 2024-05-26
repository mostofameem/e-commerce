package handlers

import (
	"ecommerce/db"
	"ecommerce/web/utils"
	"encoding/json"
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
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = utils.Validate(user)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = db.Create(user.Name, user.Email, user.Password)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}
	utils.SendBothData(w, user, "Register successful ")
}
