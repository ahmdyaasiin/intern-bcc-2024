package entity

import "github.com/google/uuid"

type Wallet struct {
	ID           uuid.UUID `gorm:"type:varchar(36);not null;primary_key"`
	UserID       uuid.UUID `gorm:"type:varchar(36);not null"`
	TotalBalance uint64    `gorm:"not null"`
	CreatedAt    int64     `gorm:"autoCreateTime:milli;not null"`
	UpdatedAt    int64     `gorm:"autoCreateTime:milli;autoUpdateTime:milli;not null"`
}
