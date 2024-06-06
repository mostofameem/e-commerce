package db

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"ecommerce/logger"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type User struct {
	Id    int `json:"id"`
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

	// Hash the password
	hashedPass := hashPassword(pass)

	// Generate a secret code for email verification
	secretCode := generateSecretCode()

	// Send the secret code to the user's email
	err := sendVerificationEmail(email, secretCode)
	if err != nil {
		slog.Error(
			"Failed to send verification email",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"email": email,
			}),
		)
		return err
	}

	// Store the secret code in Redis with an expiration time (e.g., 5 minutes)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = GetRedis().Set(ctx, email, secretCode, 5*time.Minute).Err()
	if err != nil {
		slog.Error(
			"Failed to store secret code in Redis",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"email": email,
			}),
		)
		return err
	}

	// Store user data and the secret code in the database
	columns := map[string]interface{}{
		"name":     name,
		"email":    email,
		"password": string(hashedPass),
		//"secret_code": secretCode,
		"isactive": false, // initially not verified
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
			"Failed to create new user",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": query,
				"args":  args,
			}),
		)
		return err
	}

	// Execute the query
	if _, err := GetWriteDB().Exec(query, args...); err != nil {
		slog.Error(
			"Failed to execute Insert query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": query,
				"args":  args,
			}),
		)
		return err
	}

	return nil
}

func (r *UserTypeRepo) Login(email string, pass string) error {

	dbpass := r.GetPass(email)
	hashedPass := hashPassword(pass)

	if dbpass == string(hashedPass) {
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

func hashPassword(pass string) string {

	h := sha1.New()
	h.Write([]byte(pass))
	hashValue := h.Sum(nil)
	return hex.EncodeToString(hashValue)
}

func (r *UserTypeRepo) GetEmail(id string) string {
	var email string

	queryString, args, err := GetQueryBuilder().
		Select("email").
		From(r.table).
		Where(sq.Eq{"id": id}).
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
		return email
	}

	err = GetReadDB().Get(&email, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return email
		}
		slog.Error(
			"Failed to get the content",
			logger.Extra(map[string]any{
				"error": err.Error(),
			}),
		)
		return email
	}
	return email
}

func (r *UserTypeRepo) Verified(id string) error {

	query, args, err := GetQueryBuilder().
		Update(r.table).
		Set("isactive", true).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create new user",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": query,
				"args":  args,
			}),
		)
		return err
	}

	// Execute the query
	if _, err := GetWriteDB().Exec(query, args...); err != nil {
		slog.Error(
			"Failed to execute query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": query,
				"args":  args,
			}),
		)
		return err
	}

	return nil
}
