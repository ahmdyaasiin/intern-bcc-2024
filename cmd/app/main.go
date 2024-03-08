package main

import (
	"intern-bcc-2024/pkg/config"
	"intern-bcc-2024/pkg/database/mysql"
)

func main() {
	config.LoadEnv()

	db := mysql.ConnectDatabase()
	mysql.Migrate(db)

}
