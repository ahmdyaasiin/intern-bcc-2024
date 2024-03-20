package model

import "github.com/google/uuid"

type RequestUpdateAccountNumber struct {
	ID            uuid.UUID `json:"id" binding:"required"`
	AccountNumber string    `json:"account_number" binding:"required"`
}
