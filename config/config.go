package config

type Config struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Dbname    string `json:"dbname"`
	Sslmode   string `json:"sslmode"`
	Http_port string `json:"http_port"`
	JwtSecret string `json:"jwt_secret"`
}

var config *Config

func init() {
	config = &Config{}
}

func GetConfig() Config {
	return *config
}
