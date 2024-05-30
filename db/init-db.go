package db

func InitDB() {
	ConnectDB()
	InitQueryBuilder()
	//MigrateDB()
}
