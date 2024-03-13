package entity

import "github.com/google/uuid"

type Category struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);not null;primary_key"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	UrlCategory string    `json:"url_category" gorm:"type:varchar(255);not null"`
	CreatedAt   int64     `json:"-" gorm:"autoCreateTime:milli;not null"`

	Product []Product `json:"-" gorm:"foreignKey:category_id;references:id"`
}
