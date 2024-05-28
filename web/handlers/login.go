package handlers

import (
	"ecommerce/auth"
	"ecommerce/db"
	"ecommerce/models"
	"ecommerce/web/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	var user LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = utils.Validate(user)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = db.Login(user.Email, user.Password)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, fmt.Errorf("wrong username / password "))
		return
	}

	var wg sync.WaitGroup
	var usr models.User

	wg.Add(1)
	go db.GetUser(user.Email, &usr, &wg)
	wg.Wait()

	accessToken, refreshToken, err := auth.GenerateToken(usr)

	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	log.Println(refreshToken)
	utils.SendBothData(w, accessToken, usr)
}
