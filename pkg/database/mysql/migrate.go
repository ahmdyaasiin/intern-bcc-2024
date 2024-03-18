package mysql

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"log"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&entity.AccountNumberType{},
		&entity.User{},
		&entity.OtpCode{},
		&entity.Session{},
		&entity.ResetToken{},
		&entity.Category{},
		&entity.Product{},
		&entity.Media{},
		&entity.Transaction{},
	); err != nil {
		log.Fatal(err)
	}

}
