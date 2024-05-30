package db

import (
	"database/sql"
	"ecommerce/logger"
	"errors"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
)

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name" validate:"required,min=5,max=20,alpha"`
	Email string `json:"email" validate:"required,email"`
}
type UserTypeRepo struct {
	table string
}

var userTypeRepo *UserTypeRepo

func initUserTypeRepo() {
	userTypeRepo = &UserTypeRepo{
		table: "users",
	}
}

func GetUserTypeRepo() *UserTypeRepo {
	return userTypeRepo
}

func (r *UserTypeRepo) Create(name, email, pass string) error {

	dbpass := r.GetPass(email)
	if dbpass != "" {
		return fmt.Errorf("User Exists")
	}

	columns := map[string]interface{}{
		"name":     name,
		"email":    email,
		"password": pass,
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
			"Failed to create New User",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": query,
				"args":  args,
			}),
		)
		return err
	}
	GetWriteDB().QueryRow(query, args...).Scan()
	return nil

}
func (r *UserTypeRepo) Login(email string, pass string) error {

	dbpass := r.GetPass(email)
	if dbpass == pass {
		return nil
	}
	return errors.New("failed ")
}

func (r *UserTypeRepo) GetPass(email string) string {
	//query := "SELECT PASSWORD from users where email ='" + email + "';"

	var password string
	// Build the query
	queryString, args, err := GetQueryBuilder().
		Select("PASSWORD").
		From(r.table).
		Where(sq.Eq{"email": email}).
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
		return password
	}

	err = GetReadDB().Get(&password, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return password
		}
		slog.Error(
			"Failed to get the content",
			logger.Extra(map[string]any{
				"error": err.Error(),
			}),
		)
		return password
	}
	return password

}
func (r *UserTypeRepo) GetUser(email string) (User, error) {

	//query := "SELECT id, email, name FROM users WHERE email = '" + email + "';"
	var user User

	queryString, args, err := GetQueryBuilder().
		Select("id", "name", "email").
		From(r.table).
		Where(sq.Eq{"email": email}).
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
		return user, err
	}

	err = GetReadDB().Get(&user, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		slog.Error(
			"Failed to get the content",
			logger.Extra(map[string]any{
				"error": err.Error(),
			}),
		)
		return user, err
	}
	return user, nil

}
