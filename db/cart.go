package db

import (
	"database/sql"
	"ecommerce/logger"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
)

type CartList struct {
	ProductName string `db:"product_name" json:"product_name"`
	Price       int    `db:"price" json:"price"`
	Quantity    int    `db:"quantity" json:"quantity"`
}

type Cart struct {
	ProductName string `json:"product_name" validate:"required,alpha"`
	Quantity    int    `json:"quantity" validate:"required"`
}
type CartTypeRepo struct {
	table string
}

var cartTypeRepo *CartTypeRepo

func initCartTypeRepo() {
	cartTypeRepo = &CartTypeRepo{
		table: "carts",
	}
}

func GetCartTypeRepo() *CartTypeRepo {
	return cartTypeRepo
}

func (r *CartTypeRepo) GetCart(id int, ch chan []CartList) {

	var AllProduct []CartList

	query, args, err := GetQueryBuilder().
		Select("product_name", "price", "quantity").
		From(r.table).
		Where(sq.Eq{"user_id": id}).
		OrderBy("quantity DESC").
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create resource types select query getcart",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": query,
				"args":  args,
			}),
		)
		ch <- []CartList{}

	}

	err = GetReadDB().Select(&AllProduct, query, args...)
	if err != nil {
		slog.Error(
			"Failed to get resource types getcart",
			logger.Extra(map[string]any{
				"error": err.Error(),
			}),
		)
		ch <- []CartList{}
	}
	ch <- AllProduct
}

func (r *CartTypeRepo) InsertToCart(item Cart, id int) error {
	product := GetProductTypeRepo().GetProduct(item)
	columns := map[string]interface{}{
		"user_id":      id,
		"product_name": product.Name,
		"price":        product.Price,
		"quantity":     item.Quantity,
	}
	var colNames []string
	var colValues []any

	for colName, colVal := range columns {
		colNames = append(colNames, colName)
		colValues = append(colValues, colVal)
	}

	query, args, err := GetQueryBuilder().
		Insert(r.table).
		Columns(colNames...).
		Values(colValues...).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create resource type insert query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": query,
				"args":  args,
			}),
		)
	}
	err = GetWriteDB().QueryRow(query, args...).Err()
	return err

}
func (r *CartTypeRepo) GiveMeTotal(id int, totalchan chan int) {
	var total int
	queryString, args, err := GetQueryBuilder().
		Select("sum(price*quantity)").
		From(r.table).
		Where(sq.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": queryString,
				"args":  args,
			}),
		)
		totalchan <- -1
	}

	err = GetReadDB().Get(&total, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			totalchan <- -1
		}
		slog.Error(
			"Failed to get the content",
			logger.Extra(map[string]any{
				"error": err.Error(),
			}),
		)
		totalchan <- -1
	}
	totalchan <- total
}
