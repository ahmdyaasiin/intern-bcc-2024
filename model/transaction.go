package model

import "github.com/google/uuid"

type ResponseForBuyProduct struct {
	PaymentID string `json:"payment_id"`
}

type PaymentNotificationHandler struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
}

type RequestForCancelTransaction struct {
	TransactionID uuid.UUID `json:"transaction_id" binding:"required"`
}

type RequestForRefuseTransaction struct {
	CancelCode string `json:"cancel_code" binding:"required"`
}

type RequestForWithdrawTransaction struct {
	WithdrawCode string `json:"withdraw_code" binding:"required"`
}

type ResponseForActiveTransactions struct {
	TransactionID  uuid.UUID `json:"transaction_id"`
	ProductID      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	ProductPrice   string    `json:"product_price"`
	WithdrawalCode string    `json:"withdrawal_code"`
	UrlProduct     string    `json:"url_product"`
	OwnerName      string    `json:"owner_name"`
}
