package mysql

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"log"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Wallet{},
		&entity.OtpCode{},
		&entity.RefreshToken{},
		&entity.ResetToken{},
		&entity.Category{},
		&entity.Product{},
		&entity.Media{},
		&entity.Admin{},
	); err != nil {
		log.Fatal(err)
	}

}
