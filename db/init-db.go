package db

func InitDB() {
	ConnectDB()
	InitQueryBuilder()
	InitRedis()
}
