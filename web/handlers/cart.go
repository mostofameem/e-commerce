package handlers

import (
	"ecommerce/auth"
	"ecommerce/db"
	"ecommerce/web/utils"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func BuyProduct(w http.ResponseWriter, r *http.Request) {
	item, err := UrlOperation(r.URL.String()) //get item details from url
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = utils.Validate(item) //item validate
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := GetIdFromHeader(r.Header.Get("Authorization")) //get user id from token
	if err != nil {
		utils.SendError(w, 404, err)
		return
	}

	err = db.InsertToCart(item, userID) //add to cart
	if err != nil {
		utils.SendError(w, 404, err)
		return
	}
	utils.SendData(w, "Added To cart successful")
}
func UrlOperation(r string) (db.Cart, error) {

	var item db.Cart
	parsedUrl, err := url.Parse(r)
	if err != nil {
		return item, err
	}
	queryParams := parsedUrl.Query()
	item.ProductName = queryParams.Get("product_name")
	item.Quantity = queryParams.Get("quantity")

	err = utils.Validate(item)
	if err != nil {
		return db.Cart{}, err
	}
	return item, nil
}
func GetIdFromHeader(r string) (string, error) {

	authHeader := r
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := auth.ParseToken(tokenString)
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

func ShowCart(w http.ResponseWriter, r *http.Request) {

	id, err := GetIdFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		utils.SendError(w, 404, err)
		return
	}
	list, total := db.ShowCart(id)
	utils.SendBothData(w, total, list)
}
