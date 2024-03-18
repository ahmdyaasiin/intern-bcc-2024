package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
	"intern-bcc-2024/pkg/validation"
)

func (r *Rest) ResendOtp(ctx *gin.Context) {
	var requests model.RequestForResend
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

	respDetails := r.service.OtpService.Resend(requests)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "The OTP code was sent successfully")
}
