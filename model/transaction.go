package model

import "github.com/google/uuid"

type ResponseForBuyProduct struct {
	PaymentID uuid.UUID `json:"payment_id"`
}
