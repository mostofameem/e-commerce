package handlers

import (
	"ecommerce/db"
	"ecommerce/logger"
	"ecommerce/web/middlewares"
	"ecommerce/web/utils"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

func BuyProduct(w http.ResponseWriter, r *http.Request) {
	item, err := UrlOperation(r.URL.String()) //get item details from url
	if err != nil {
		slog.Error("Error item info ", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": item,
		}))
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	err = utils.Validate(item) //item validate
	if err != nil {
		slog.Error("Failed to validate item data", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": item,
		}))
		utils.SendError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := middlewares.GetIdFromHeader(r.Header.Get("Authorization")) //get user id from token
	if err != nil {
		slog.Error("Failed to get UserId", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": userID,
		}))
		utils.SendError(w, 404, err)
		return
	}

	err = db.GetCartTypeRepo().InsertToCart(item, userID) //add to cart
	if err != nil {
		slog.Error("Failed to insert into cart ", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": item,
			"User Id": userID,
		}))
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

func ShowCart(w http.ResponseWriter, r *http.Request) {

	id, err := middlewares.GetIdFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		utils.SendError(w, 404, err)
		return
	}
	listch := make(chan []db.CartList)
	totalch := make(chan string)

	go db.GetCartTypeRepo().GetCart(id, listch)
	go db.GetCartTypeRepo().GiveMeTotal(id, totalch)

	list := <-listch
	total := <-totalch
	utils.SendBothData(w, fmt.Sprintf("Total =%s", total), list)
}
