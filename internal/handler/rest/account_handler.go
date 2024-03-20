package rest

import (
	"github.com/gin-gonic/gin"
	"intern-bcc-2024/pkg/response"
)

func (r *Rest) AllAccountNumber(ctx *gin.Context) {
	accountType, respDetails := r.service.AccountService.AllAccountNumber()
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "Success get products for homepage", accountType)
}
