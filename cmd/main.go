package main

import (
	"ecommerce/db"
	"ecommerce/web"
	"log"
	"net/http"
)

func main() {

	if err := db.InitDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Db.Close()

	mux := web.StartServer()
	log.Printf("Server Running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", mux))

}
