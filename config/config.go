package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
	Sslmode  string `json:"sslmode"`
	S_port   string `json:"s_port"`
}

var config Config

func ReadConfig() {
	configFile, err := os.Open("../config.json")

	if err == nil {
		defer configFile.Close()

		jsonParser := json.NewDecoder(configFile)
		if err = jsonParser.Decode(&config); err != nil {
			log.Fatal("Error Loading Config File")
			return
		}
	}
}
func GetConfig() Config {
	return config
}
