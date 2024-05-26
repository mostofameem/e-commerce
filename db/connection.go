package db

import (
	"database/sql"
	"ecommerce/config"
	"fmt"

	_ "github.com/lib/pq"
)

var Db *sql.DB
var err error

func InitDB() error {

	config := config.GetConfig()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Dbname, config.Sslmode)

	Db, err = sql.Open("postgres", connStr)

	return err
}
func Close() {
	Db.Close()
}
