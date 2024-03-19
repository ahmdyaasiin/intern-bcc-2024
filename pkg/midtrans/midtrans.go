package midtrans

import (
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/pkg/response"
	"os"
)

func CreateToken(idTransaction uuid.UUID, product *entity.Product) (string, response.Details) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  idTransaction.String(),
			GrossAmt: int64(product.Price),
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeGopay,
			snap.PaymentTypeShopeepay,
			snap.PaymentTypeBankTransfer,
		},
		Expiry: &snap.ExpiryDetails{
			Duration: 1,
			Unit:     "minute",
		},
	}

	var client snap.Client
	client.New(serverKey, midtrans.Sandbox)
	client.Options.SetPaymentOverrideNotification("https://seahorse-cool-kodiak.ngrok-free.app/api/v1/product/payment/callback")

	snapResp, err := client.CreateTransactionToken(req)
	if err != nil {
		return "", response.Details{Code: 500, Message: "Failed to create token", Error: err}
	}

	return snapResp, response.Details{Code: 201, Message: "Success to create token", Error: nil}
}

func VerifyPayment(idTransaction uuid.UUID) (*coreapi.TransactionStatusResponse, error) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")

	var client coreapi.Client
	client.New(serverKey, midtrans.Sandbox)

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := client.CheckTransaction(idTransaction.String())
	if e != nil {
		return nil, e
	} else {
		return transactionStatusResp, nil
	}
}
