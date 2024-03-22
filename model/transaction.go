package model

import "github.com/google/uuid"

/*
	Request Struct
*/

type RequestForRefuseTransaction struct {
	CancelCode string `json:"cancel_code" binding:"required"`
}

type RequestForWithdrawTransaction struct {
	WithdrawCode string `json:"withdraw_code" binding:"required"`
}

type RequestForPaymentNotificationHandler struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
}

/*
	Response Struct
*/

type ResponseForBuyProduct struct {
	PaymentID string `json:"payment_id"`
}

type ResponseForActiveTransactions struct {
	TransactionID  uuid.UUID `json:"transaction_id"`
	ProductID      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	ProductPrice   string    `json:"product_price"`
	WithdrawalCode string    `json:"withdrawal_code"`
	UrlProduct     string    `json:"url_product"`
	OwnerName      string    `json:"owner_name"`
	OwnerID        string    `json:"owner_id"`
}
