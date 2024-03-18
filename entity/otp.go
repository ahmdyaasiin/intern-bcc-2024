package entity

import "github.com/google/uuid"

type OtpCode struct {
	ID        uuid.UUID `gorm:"type:varchar(36);not null;primaryKey"`
	UserID    uuid.UUID `gorm:"type:varchar(36);not null"`
	Code      string    `gorm:"type:varchar(6);not null;unique"`
	CreatedAt int64     `gorm:"autoCreateTime:milli;not null"`
	UpdatedAt int64     `gorm:"autoCreateTime:milli;autoUpdateTime:milli;not null"`
}
