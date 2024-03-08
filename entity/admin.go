package entity

import "github.com/google/uuid"

type Admin struct {
	ID        uuid.UUID `gorm:"type:varchar(36);not null;primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Username  string    `gorm:"type:varchar(18);not null"`
	Password  string    `gorm:"type:varchar(64);not null"`
	CreatedAt int64     `gorm:"autoCreateTime:milli;not null"`
	UpdatedAt int64     `gorm:"autoCreateTime:milli;autoUpdateTime:milli;not null"`
}
