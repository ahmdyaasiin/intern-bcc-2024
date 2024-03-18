package entity

import "github.com/google/uuid"

type Media struct {
	ID        uuid.UUID `gorm:"type:varchar(36);not null;primaryKey"`
	ProductID uuid.UUID `gorm:"type:varchar(36);not null"`
	Url       string    `gorm:"type:varchar(255);not null"`
	CreatedAt int64     `gorm:"autoCreateTime:milli;not null"`
	UpdatedAt int64     `gorm:"autoCreateTime:milli;autoUpdateTime:milli;not null"`
}
