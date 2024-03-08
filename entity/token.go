package entity

import "github.com/google/uuid"

type ResetToken struct {
	ID        uuid.UUID `gorm:"type:varchar(36);not null;primary_key"`
	UserID    uuid.UUID `gorm:"type:varchar(36);not null"`
	Token     string    `gorm:"type:varchar(6);not null"`
	ExpiredAt int64     `gorm:"-"`
	CreatedAt int64     `gorm:"autoCreateTime:milli;not null"`
}
