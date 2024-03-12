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
	StatusAccount   string    `json:"-" gorm:"type:enum('blocked', 'inactive', 'active');not null"`
	UrlPhotoProfile string    `gorm:"type:varchar(255);not null"`
	CreatedAt       int64     `json:"-" gorm:"autoCreateTime:milli;not null"`
	UpdatedAt       int64     `json:"-" gorm:"autoCreateTime:milli;autoUpdateTime:milli;not null"`

	Wallet       []Wallet       `json:"-" gorm:"foreignKey:user_id;references:id"`
	OtpCode      []OtpCode      `json:"-" gorm:"foreignKey:user_id;references:id"`
	RefreshToken []RefreshToken `json:"-" gorm:"foreignKey:user_id;references:id"`
	ResetToken   []ResetToken   `json:"-" gorm:"foreignKey:user_id;references:id"`
	Product      []Product      `json:"-" gorm:"foreignKey:user_id;references:id"`
}
