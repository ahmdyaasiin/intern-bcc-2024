package entity

import "github.com/google/uuid"

type Transaction struct {
	ID             uuid.UUID `json:"id" gorm:"type:varchar(36);not null;primary_key"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null"`
	ProductID      uuid.UUID `json:"product_id" gorm:"type:varchar(36);not null"`
	MidtransID     uuid.UUID `gorm:""`
	Amount         uint64    `json:"amount" gorm:"not null"`
	Status         string    `json:"status" gorm:"type:enum('unpaid', 'paid', 'completed')"`
	WithdrawalCode string    `json:""`
	CreatedAt      int64     `json:"-" gorm:"autoCreateTime:milli;not null"`
	UpdatedAt      int64     `json:"-" gorm:"autoCreateTime:milli;autoUpdateTime:milli;not null"`
}
