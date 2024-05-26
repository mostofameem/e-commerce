package main

import (
	"ecommerce/config"
	"ecommerce/db"
	"ecommerce/web"
	"log"
	"net/http"
)

func main() {
	config.ReadConfig()

	if err := db.InitDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Db.Close()

	mux := web.StartServer()
	log.Printf("Server Running on port 3000")
	log.Fatal(http.ListenAndServe(":"+config.GetConfig().S_port, mux))

}
