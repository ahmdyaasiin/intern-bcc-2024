package entity

import "github.com/google/uuid"

type AccountNumberType struct {
	ID        uuid.UUID `gorm:"type:varchar(36);not null;primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	CreatedAt int64     `gorm:"autoCreateTime:milli;not null"`
	UpdatedAt int64     `gorm:"autoCreateTime:milli;autoUpdateTime:milli;not null"`

	User []User `json:"-" gorm:"foreignKey:account_number_id;references:id"`
}
