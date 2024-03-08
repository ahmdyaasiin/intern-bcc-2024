package entity

import "github.com/google/uuid"

type Category struct {
	ID          uuid.UUID `gorm:"type:varchar(36);not null;primary_key"`
	Name        string    `gorm:"type:varchar(255);not null"`
	UrlCategory string    `gorm:"type:varchar(255);not null"`
	CreatedAt   int64     `gorm:"autoCreateTime:milli;not null"`

	Product []Product `gorm:"foreignKey:category_id;references:id"`
}
