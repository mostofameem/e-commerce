package db

import (
	"database/sql"
	"ecommerce/logger"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
)

type Product struct {
	Name     string `json:"name" db:"name"`
	Price    int    `json:"price" db:"price"`
	Quantity int    `json:"quantity" db:"quantity"`
}

type ProductTypeRepo struct {
	table string
}

var productTypeRepo *ProductTypeRepo

func initProductTypeRepo() {
	productTypeRepo = &ProductTypeRepo{
		table: "products",
	}
}

func GetProductTypeRepo() *ProductTypeRepo {
	return productTypeRepo
}

func (r *ProductTypeRepo) GetProduct(item Cart) Product {
	//query := "SELECT name,price,quantity from products where name='" + item.ProductName + "';"
	product := Product{}

	// Build the query
	queryString, args, err := GetQueryBuilder().
		Select("name", "price", "quantity").
		From(r.table).
		Where(sq.Eq{"name": item.ProductName}).
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
		return product
	}

	err = GetReadDB().Get(&product, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return product
		}
		slog.Error(
			"Failed to get the content",
			logger.Extra(map[string]any{
				"error": err.Error(),
			}),
		)
		return product
	}
	return product
}
