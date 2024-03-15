package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"intern-bcc-2024/pkg/response"
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
