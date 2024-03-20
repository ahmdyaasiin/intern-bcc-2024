package entity

import "github.com/google/uuid"

type AccountNumberType struct {
	ID        uuid.UUID `json:"id" gorm:"type:varchar(36);not null;primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	CreatedAt int64     `json:"-" gorm:"autoCreateTime:milli;not null"`
	UpdatedAt int64     `json:"-" gorm:"autoCreateTime:milli;autoUpdateTime:milli;not null"`

	User []User `json:"-" gorm:"foreignKey:account_number_id;references:id"`
}
