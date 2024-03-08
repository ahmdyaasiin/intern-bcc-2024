package entity

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `gorm:"type:varchar(36);not null;primary_key"`
	UserID      uuid.UUID `gorm:"type:varchar(36);not null"`
	CategoryID  uuid.UUID `gorm:"type:varchar(36);not null"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text;not null"`
	Price       uint64    `gorm:"not null"`
	CreatedAt   int64     `gorm:"autoCreateTime:milli;not null"`
	UpdatedAt   int64     `gorm:"autoCreateTime:milli;autoUpdateTime:milli;not null"`

	Media []Media `gorm:"foreignKey:product_id;references:id"`
}
