package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/database/mysql"
	"intern-bcc-2024/pkg/jwt"
	"intern-bcc-2024/pkg/mail"
	md "intern-bcc-2024/pkg/midtrans"
	"intern-bcc-2024/pkg/response"
	"log"
)

type ITransactionService interface {
	BuyProduct(ctx *gin.Context, id uuid.UUID) (model.ResponseForBuyProduct, response.Details)
	CheckPayment(idTransaction uuid.UUID) response.Details
	FindActiveTransactions(ctx *gin.Context) (*[]model.ResponseForActiveTransactions, response.Details)
	CancelTransaction(ctx *gin.Context, idTransaction uuid.UUID, requestTransactionID uuid.UUID) response.Details
	RefuseTransaction(ctx *gin.Context, id uuid.UUID, requests model.RequestForRefuseTransaction) response.Details
	AcceptTransaction(ctx *gin.Context, id uuid.UUID, requests model.RequestForWithdrawTransaction) response.Details
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
	transaction := new(entity.Transaction)

	tx := ts.db.Begin()
	defer tx.Rollback()

	user, err := ts.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return model.ResponseForBuyProduct{}, response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	defaultUUID, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if err != nil {
		log.Println(err)

		return model.ResponseForBuyProduct{}, response.Details{Code: 500, Message: "Failed to convert default uuid", Error: err}
	}

	if user.AccountNumber == "0" || user.AccountNumberID == defaultUUID {
		log.Println("user haven't set account number and account number id")

		return model.ResponseForBuyProduct{}, response.Details{Code: 500, Message: "User haven't set account number and account number id", Error: errors.New("user haven't set account number and account number id")}
	}

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

	respDetails = ts.trr.Find(tx, transaction, model.ParamForFind{
		ProductID: product.ID,
	})
	if respDetails.Error == nil {
		if transaction.Status == "completed" {
			log.Println("the item was sold")

			return model.ResponseForBuyProduct{}, response.Details{Code: 403, Message: "the item was sold", Error: errors.New("the item was sold")}
		} else {
			log.Println("there's a transaction is on progress")

			return model.ResponseForBuyProduct{}, response.Details{Code: 403, Message: "please kindly wait, someone is on transaction for this item", Error: errors.New("please kindly wait, someone is on transaction for this item")}
		}
	}

	idTransaction := uuid.New()
	transaction = &entity.Transaction{
		ID:             idTransaction,
		UserID:         user.ID,
		ProductID:      product.ID,
		Amount:         product.Price,
		Status:         "unpaid",
		WithdrawalCode: mail.GenerateSixCode(),
	}

	paymentID, respDetails := md.CreateToken(idTransaction, product)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return model.ResponseForBuyProduct{}, respDetails
	}

	paymentIDParse, err := uuid.Parse(paymentID)
	if err != nil {
		log.Println("failed convert payment id to uuid")

		return model.ResponseForBuyProduct{}, response.Details{Code: 500, Message: "Failed to convert", Error: errors.New("failed to convert")}
	}

	transaction.MidtransID = paymentIDParse

	respDetails = ts.trr.CreateTransaction(tx, transaction)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return model.ResponseForBuyProduct{}, respDetails
	}

	product.CancelCode = mail.GenerateSixCode()
	if respDetails = ts.pr.Update(tx, product); respDetails.Error != nil {
		log.Println(err)

		return model.ResponseForBuyProduct{}, respDetails
	}

	if err = tx.Commit().Error; err != nil {
		return model.ResponseForBuyProduct{}, response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return model.ResponseForBuyProduct{
		PaymentID: paymentID, // seharusnya snap id midtrans
	}, response.Details{Code: 200, Message: "Success create transaction", Error: nil}
}

func (ts *TransactionService) CheckPayment(idTransaction uuid.UUID) response.Details {
	transaction := new(entity.Transaction)

	tx := ts.db.Begin()
	defer tx.Rollback()

	respDetails := ts.trr.Find(tx, transaction, model.ParamForFind{
		ID: idTransaction,
	})
	if respDetails.Error != nil {
		log.Println("transaction not found")

		return respDetails
	}

	resp, err := md.VerifyPayment(idTransaction)
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to verify payment", Error: err}
	}

	if (resp.TransactionStatus == "capture" && resp.FraudStatus == "accept") || resp.TransactionStatus == "settlement" {
		// success - set transaction status to success and generate withdrawal token
		transaction.Status = "paid"
		if err = ts.trr.Update(tx, transaction).Error; err != nil {
			return response.Details{Code: 500, Message: "Failed update transaction", Error: err}
		}

		if err = tx.Commit().Error; err != nil {
			return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
		}

		return response.Details{Code: 200, Message: "Transaction success", Error: nil}
	} else if resp.TransactionStatus == "cancel" || resp.TransactionStatus == "expire" {
		// expired - delete the transaction
		if err = ts.trr.Delete(tx, transaction).Error; err != nil {
			return response.Details{Code: 500, Message: "Failed to delete transaction", Error: err}
		}

		if err = tx.Commit().Error; err != nil {
			return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
		}

		return response.Details{Code: 200, Message: "Transaction cancel or expired", Error: nil}
	} else if resp.TransactionStatus == "pending" {
		return response.Details{Code: 200, Message: "Transaction pending", Error: nil}
	}

	// TODO set transaction status
	// TODO Generate Token for COD

	log.Println(resp.TransactionStatus)

	return response.Details{Code: 200, Message: fmt.Sprintf("Transaction %s", resp.TransactionStatus), Error: nil}
}

func (ts *TransactionService) FindActiveTransactions(ctx *gin.Context) (*[]model.ResponseForActiveTransactions, response.Details) {
	transaction := new([]model.ResponseForActiveTransactions)

	tx := ts.db.Begin()
	defer tx.Rollback()

	user, err := ts.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return transaction, response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	respDetails := ts.trr.FindActiveTransactions(tx, transaction, user)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return transaction, respDetails
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return transaction, response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return transaction, response.Details{Code: 200, Message: "Success get all active products", Error: nil}
}

func (ts *TransactionService) CancelTransaction(ctx *gin.Context, idTransaction uuid.UUID, requestTransactionID uuid.UUID) response.Details {
	transaction := new(entity.Transaction)

	tx := ts.db.Begin()
	defer tx.Rollback()

	if idTransaction != requestTransactionID {
		return response.Details{Code: 401, Message: "idTransaction param is not equals to idTransaction request"}
	}

	respDetails := ts.trr.Find(tx, transaction, model.ParamForFind{
		ID: idTransaction,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	user, err := ts.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	if transaction.UserID != user.ID {
		return response.Details{Code: 403, Message: "It's not your transaction", Error: errors.New("different owner of transaction")}
	}

	respDetails = ts.trr.Delete(tx, transaction)
	if respDetails.Error != nil {
		return respDetails
	}

	// TODO send the money back to buyer but charge the fee

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success cancel the transaction", Error: nil}
}

func (ts *TransactionService) RefuseTransaction(ctx *gin.Context, id uuid.UUID, requests model.RequestForRefuseTransaction) response.Details {
	transaction := new(entity.Transaction)
	product := new(entity.Product)

	tx := ts.db.Begin()
	defer tx.Rollback()

	user, err := ts.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	respDetails := ts.trr.Find(tx, transaction, model.ParamForFind{
		ID: id,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	if transaction.UserID != user.ID {
		return response.Details{Code: 401, Message: "It's not your transaction", Error: errors.New("it's not your transaction")}
	}

	respDetails = ts.pr.Find(tx, product, model.ParamForFind{
		ID: transaction.ProductID,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	if requests.CancelCode != product.CancelCode {
		return response.Details{Code: 401, Message: "Wrong cancel code", Error: errors.New("wrong cancel code")}
	}

	// TODO delete transaction
	respDetails = ts.trr.Delete(tx, transaction)
	if respDetails.Error != nil {
		return respDetails
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success refuse the transaction", Error: nil}
}

func (ts *TransactionService) AcceptTransaction(ctx *gin.Context, id uuid.UUID, requests model.RequestForWithdrawTransaction) response.Details {
	transaction := new(entity.Transaction)

	tx := ts.db.Begin()
	defer tx.Rollback()

	user, err := ts.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	respDetails := ts.trr.Find(tx, transaction, model.ParamForFind{
		ID: id,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	if transaction.UserID != user.ID {
		return response.Details{Code: 401, Message: "It's not your transaction", Error: errors.New("it's not your transaction")}
	}

	if requests.WithdrawCode != transaction.WithdrawalCode {
		return response.Details{Code: 401, Message: "Wrong withdrawal code", Error: errors.New("wrong withdrawal code")}
	}

	// TODO delete transaction
	transaction.Status = "completed"
	respDetails = ts.trr.Update(tx, transaction)
	if respDetails.Error != nil {
		return respDetails
	}

	// TODO withdraw to seller

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success accept the transaction", Error: nil}
}
