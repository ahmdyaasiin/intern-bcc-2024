package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
	"intern-bcc-2024/pkg/validation"
	"log"
	"strings"
)

func (r *Rest) RenewSession(ctx *gin.Context) {
	var requests model.RequestForRenewAccessToken

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

	token, respDetails := r.service.SessionService.RenewSession(requests)
	if respDetails.Code == 401 {
		if strings.Contains(respDetails.Message, "invalid") {
			response.WithExpired(ctx, respDetails.Code, respDetails.Message, false)
		} else {
			response.WithExpired(ctx, respDetails.Code, respDetails.Message, true)
		}

		return
	}

	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "Successfully updated access token", token)
}

func (r *Rest) Logout(ctx *gin.Context) {
	respDetails := r.service.SessionService.Logout(ctx)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "User logout successfully")
}
