package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/jwt"
	md "intern-bcc-2024/pkg/midtrans"
	"intern-bcc-2024/pkg/response"
)

type ITransactionService interface {
	BuyProduct(ctx *gin.Context, id uuid.UUID) (model.ResponseForBuyProduct, response.Details)
}

type TransactionService struct {
	pr      repository.IProductRepository
	trr     repository.ITransactionRepository
	jwtAuth jwt.Interface
}

func NewTransactionService(productRepository repository.IProductRepository, transactionRepository repository.ITransactionRepository, jwtAuth jwt.Interface) ITransactionService {
	return &TransactionService{
		pr:      productRepository,
		trr:     transactionRepository,
		jwtAuth: jwtAuth,
	}
}

func (ts *TransactionService) BuyProduct(ctx *gin.Context, id uuid.UUID) (model.ResponseForBuyProduct, response.Details) {
	user, err := ts.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		return model.ResponseForBuyProduct{}, response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	product, respDetails := ts.pr.Find(model.ParamForFind{
		ID: id,
	})
	if respDetails.Error != nil {
		return model.ResponseForBuyProduct{}, respDetails
	}

	transaction := entity.Transaction{
		ID:        uuid.New(),
		UserID:    user.ID,
		ProductID: product.ID,
		Amount:    product.Price,
		Status:    "on progress",
	}

	respDetails = ts.trr.CreateTransaction(&transaction)
	if respDetails.Error != nil {
		return model.ResponseForBuyProduct{}, respDetails
	}

	paymentID, respDetails := md.CreateToken(&product)
	if respDetails.Error != nil {
		return model.ResponseForBuyProduct{}, respDetails
	}

	return model.ResponseForBuyProduct{
		PaymentID: paymentID, // seharusnya snap id midtrans
	}, response.Details{Code: 200, Message: "Success create transaction", Error: nil}
}
