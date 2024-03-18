package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
	"intern-bcc-2024/pkg/validation"
	"log"
)

func (r *Rest) ResetPassword(ctx *gin.Context) {
	var requests model.RequestForReset

	if err := ctx.ShouldBindJSON(&requests); err != nil {
		var ve validator.ValidationErrors
		errorList := validation.GetError(err, ve)
		if errorList != nil {
			log.Println("Failed to validate user requests")

			response.WithErrors(ctx, 422, "Failed to validate user requests", errorList)
			return
		}

		log.Println("Failed to bind requests")

		response.MessageOnly(ctx, 422, "Failed to bind requests")
		return
	}

	respDetails := r.service.TokenService.ResetPassword(requests)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "Successfully sent the password reset link")
}

func (r *Rest) CheckResetToken(ctx *gin.Context) {
	token := ctx.Param("token")
	if respDetails := r.service.TokenService.CheckToken(token); respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "Token is valid")
}
