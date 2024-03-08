package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `gorm:"type:varchar(36);not null;primary_key"`
	Name            string    `gorm:"type:varchar(18);not null"`
	Email           string    `gorm:"type:varchar(255);not null;unique"`
	Password        string    `gorm:"type:varchar(255);not null"`
	Address         string    `gorm:"type:varchar(255);not null"`
	Latitude        float64   `gorm:"type:float(10,6);not null"`
	Longitude       float64   `gorm:"type:float(10,6);not null"`
	StatusAccount   string    `gorm:"type:enum('blocked', 'inactive', 'active');not null"`
	UrlPhotoProfile string    `gorm:"type:varchar(255);not null"`
	CreatedAt       int64     `gorm:"autoCreateTime:milli;not null"`
	UpdatedAt       int64     `gorm:"autoCreateTime:milli;autoUpdateTime:milli;not null"`

	Wallet       []Wallet       `gorm:"foreignKey:user_id;references:id"`
	OtpCode      []OtpCode      `gorm:"foreignKey:user_id;references:id"`
	RefreshToken []RefreshToken `gorm:"foreignKey:user_id;references:id"`
	ResetToken   []ResetToken   `gorm:"foreignKey:user_id;references:id"`
	Product      []Product      `gorm:"foreignKey:user_id;references:id"`
}
