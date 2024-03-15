package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"intern-bcc-2024/pkg/config"
	"log"
)

var Connection *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open(config.LoadDataSourceName()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal(err)
	}

	Connection = db
}
