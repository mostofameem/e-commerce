package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var Db *sql.DB
var err error

func InitDB() error {

	connStr := "host=localhost port=5432 user=root  password=admin dbname=Db sslmode=disable"

	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return nil
}
func Close() {
	Db.Close()
}
