package db

import "ecommerce/models"

func GetProduct(item Cart) models.Product {
	query := "SELECT name,price,quantity from products where name='" + item.ProductName + "';"
	product := models.Product{}
	err = Db.QueryRow(query).Scan(&product.Name, &product.Price, &product.Quantity)
	if err != nil {
		return models.Product{}
	}
	return product
}
func GiveMeTotal(id string, totalchan chan string) {
	query := "SELECT sum(price*quantity) from cart where user_id=" + id + ";"
	var total string
	Db.QueryRow(query).Scan(&total)

	totalchan <- total
}
