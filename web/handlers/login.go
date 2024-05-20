package handlers

import (
	"ecommerce/auth"
	"ecommerce/web/utils"
	"encoding/json"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	id := CheckUser(user)
	if id == -1 {
		utils.SendData(w, "Invalid username or password")
		return
	}
	token, err := auth.GenerateToken(user.Username, id)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	utils.SendData(w, token)

}

func CheckUser(user User) int {

	for i := 0; i < len(Users); i++ {
		if Users[i].Username == user.Username && Users[i].Password == user.Password {
			return i
		}
	}
	return -1

}
