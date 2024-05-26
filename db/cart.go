package db

import (
	"ecommerce/models"
	"fmt"
)

type Cart struct {
	ProductName string `json:"product_name" validate:"required,alpha"`
	Quantity    string `json:"quantity" validate:"required"`
}

func GetCart(id string, ch chan []models.CartList) {
	query := "SELECT product_name,quantity,price from cart where user_id=" + id + ";"
	rows, err := Db.Query(query)
	if err != nil {
		ch <- []models.CartList{}
	}

	var AllProduct []models.CartList

	for rows.Next() {
		var Product models.CartList
		err := rows.Scan(&Product.ProductName, &Product.Quantity, &Product.Price)
		if err != nil {
			ch <- []models.CartList{}
		}
		AllProduct = append(AllProduct, Product)
	}
	ch <- AllProduct
}

func InsertToCart(item Cart, id string) error {

	product := GetProduct(item)
	query := "INSERT INTO cart values('" + id + "','" + product.Name + "','" + product.Price + "','" + item.Quantity + "');"
	_, err := Db.Exec(query)
	return err

}
func ShowCart(id string) ([]models.CartList, string) {
	listch := make(chan []models.CartList)
	totalch := make(chan string)

	go GiveMeTotal(id, totalch)
	go GetCart(id, listch)

	list := <-listch
	total := <-totalch
	total = fmt.Sprintf("Total =%s", total)
	return list, total

}

func StringToInt(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		n = (n*10 + int(s[i]-'0'))
	}
	return n
}
