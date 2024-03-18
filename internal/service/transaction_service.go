package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/database/mysql"
	"intern-bcc-2024/pkg/jwt"
	md "intern-bcc-2024/pkg/midtrans"
	"intern-bcc-2024/pkg/response"
	"log"
)

type ITransactionService interface {
	BuyProduct(ctx *gin.Context, id uuid.UUID) (model.ResponseForBuyProduct, response.Details)
	VerifyPayment(idTransaction uuid.UUID) response.Details
}

type TransactionService struct {
	db      *gorm.DB
	pr      repository.IProductRepository
	trr     repository.ITransactionRepository
	jwtAuth jwt.Interface
}

func NewTransactionService(productRepository repository.IProductRepository, transactionRepository repository.ITransactionRepository, jwtAuth jwt.Interface) ITransactionService {
	return &TransactionService{
		db:      mysql.Connection,
		pr:      productRepository,
		trr:     transactionRepository,
		jwtAuth: jwtAuth,
	}
}

func (ts *TransactionService) BuyProduct(ctx *gin.Context, id uuid.UUID) (model.ResponseForBuyProduct, response.Details) {
	product := new(entity.Product)

	tx := ts.db.Begin()
	defer tx.Rollback()

	user, err := ts.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return model.ResponseForBuyProduct{}, response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	// TODO check if there's a transaction is on progress
	// TODO return err if buyer_id == owner_id

	respDetails := ts.pr.Find(tx, product, model.ParamForFind{
		ID: id,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return model.ResponseForBuyProduct{}, respDetails
	}

	if user.ID == product.UserID {
		log.Println("you're the owner of the product")

		return model.ResponseForBuyProduct{}, response.Details{Code: 403, Message: "You're the owner of the product", Error: errors.New("you're the owner")}
	}

	idTransaction := uuid.New()
	transaction := &entity.Transaction{
		ID:        idTransaction,
		UserID:    user.ID,
		ProductID: product.ID,
		Amount:    product.Price,
		Status:    "on progress",
	}

	paymentID, respDetails := md.CreateToken(idTransaction, product)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return model.ResponseForBuyProduct{}, respDetails
	}

	respDetails = ts.trr.CreateTransaction(tx, transaction)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return model.ResponseForBuyProduct{}, respDetails
	}

	return model.ResponseForBuyProduct{
		PaymentID: paymentID, // seharusnya snap id midtrans
	}, response.Details{Code: 200, Message: "Success create transaction", Error: nil}
}

func (ts *TransactionService) VerifyPayment(idTransaction uuid.UUID) response.Details {
	status, err := md.VerifyPayment(idTransaction)
	if err != nil {
		return response.Details{Code: 500, Message: "Failed to verify payment", Error: err}
	}

	if status == false {
		return response.Details{Code: 403, Message: "Transaction haven't paid yet", Error: errors.New("unpaid")}
	}

	// TODO set transaction status
	// TODO Generate Token for COD

	return response.Details{Code: 200, Message: "Success set the transaction status", Error: nil}
}
