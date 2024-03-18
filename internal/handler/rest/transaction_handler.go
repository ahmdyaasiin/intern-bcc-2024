package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
	"intern-bcc-2024/pkg/validation"
)

func (r *Rest) BuyProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.MessageOnly(ctx, 422, "Failed convert id")
		return
	}

	transaction, respDetails := r.service.TransactionService.BuyProduct(ctx, id)
	if respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "Success create transaction", transaction)
}

func (r *Rest) CheckPayment(ctx *gin.Context) {
	var requests model.PaymentNotificationHandler
	if err := ctx.ShouldBindJSON(&requests); err != nil {
		var ve validator.ValidationErrors
		errorList := validation.GetError(err, ve)
		if errorList != nil {
			response.WithErrors(ctx, 422, "Failed to validate user requests", errorList)
			return
		}

		response.MessageOnly(ctx, 422, "Failed to bind requests")
		return
	}

	respDetails := r.service.TransactionService.VerifyPayment(requests.OrderID)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "Success verify the payment")
}

func (r *Rest) AllMyTransaction(ctx *gin.Context) {

}

func (r *Rest) CancelTransaction(ctx *gin.Context) {

}

func (r *Rest) RefuseTransaction(ctx *gin.Context) {

}

func (r *Rest) AcceptTransaction(ctx *gin.Context) {

}
