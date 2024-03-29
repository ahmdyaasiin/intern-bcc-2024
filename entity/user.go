package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `json:"id" gorm:"type:varchar(36);not null;primaryKey"`
	Name            string    `json:"name" gorm:"type:varchar(255);not null"`
	Email           string    `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password        string    `json:"-" gorm:"type:varchar(255);not null"`
	Address         string    `json:"address" gorm:"type:varchar(255);not null"`
	Latitude        float64   `json:"latitude" gorm:"type:float(10,6);not null"`
	Longitude       float64   `json:"longitude" gorm:"type:float(10,6);not null"`
	StatusAccount   string    `json:"-" gorm:"type:enum('blocked', 'inactive', 'active');not null"`
	AccountNumber   string    `json:"account_number" gorm:""`
	AccountNumberID uuid.UUID `json:"account_number_id" gorm:""`
	UrlPhotoProfile string    `json:"url_photo_profile" gorm:"type:varchar(255);not null"`
	CreatedAt       int64     `json:"-" gorm:"autoCreateTime:milli;not null"`
	UpdatedAt       int64     `json:"-" gorm:"autoCreateTime:milli;autoUpdateTime:milli;not null"`

	OtpCode      []OtpCode     `json:"-" gorm:"foreignKey:user_id;references:id"`
	RefreshToken []Session     `json:"-" gorm:"foreignKey:user_id;references:id"`
	ResetToken   []ResetToken  `json:"-" gorm:"foreignKey:user_id;references:id"`
	Product      []Product     `json:"-" gorm:"foreignKey:user_id;references:id"`
	Transaction  []Transaction `json:"-" gorm:"foreignKey:user_id;references:id"`
}
