package models

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name" validate:"required,min=5,max=20,alpha"`
	Email string `json:"email" validate:"required,email"`
}

type Product struct {
	Name     string `json:"name" db:"name"`
	Price    string `json:"price" db:"price"`
	Quantity string `json:"quantity" db:"quantity"`
}

type CartList struct {
	ProductName string `json:"product_name"`
	Quantity    string `json:"quantity"`
	Price       string `json:"price"`
}
