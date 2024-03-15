package midtrans

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/pkg/response"
)

var s snap.Client

func CreateToken(product *entity.Product) (string, response.Details) {
	serverKey := "SB-Mid-server-TQt21w4SNEN7CYjLvOt6vQs1"
	s.New(serverKey, midtrans.Sandbox)

	resp, err := s.CreateTransactionToken(&snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  product.ID.String(),
			GrossAmt: int64(product.Price),
		},
		EnabledPayments: snap.AllSnapPaymentType,
	})
	if err != nil {
		return "", response.Details{Code: 500, Message: "Failed to create token", Error: err}
	}

	return resp, response.Details{Code: 201, Message: "Success to create token", Error: nil}
}
