package db

import (
	"ecommerce/models"
	"errors"
	"fmt"
	"sync"
)

func Create(name, email, pass string) error {

	dbpass := GetPass(email)
	if dbpass == "" {
		err := INSERT(name, email, pass)
		return err
	}
	return fmt.Errorf("user already exists")
}
func Login(email string, pass string) error {

	dbpass := GetPass(email)
	if dbpass == pass {
		return nil
	}
	return errors.New("failed ")
}

func GetPass(email string) string {
	query := "SELECT PASSWORD from users where email ='" + email + "';"

	var password string
	Db.QueryRow(query).Scan(&password)

	return password

}
func INSERT(name, email, pass string) error {

	query := "INSERT into users (name,email,password) VALUES ('" + name + "','" + email + "', '" + pass + "');"
	_, err := Db.Exec(query)
	return err
}
func GetUser(email string, user *models.User, wg *sync.WaitGroup) {

	defer wg.Done()
	query := "SELECT id, email, name FROM users WHERE email = '" + email + "';"

	err = Db.QueryRow(query).Scan(&user.Id, &user.Email, &user.Name)
	if err != nil {
		*user = models.User{} // reset user if there's an error
	}

}
