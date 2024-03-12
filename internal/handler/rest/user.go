package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
	"intern-bcc-2024/pkg/validation"
)

func (r *Rest) RegisterAccount(ctx *gin.Context) {
	var requests model.RequestForRegister
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

	user, respDetails := r.service.UserService.Register(requests)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 201, "User has been successfully registered", model.ResponseForRegister{
		ID: user.ID,
	})

}

func (r *Rest) VerifyAccount(ctx *gin.Context) {
	var requests model.RequestForVerify
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

	respDetails := r.service.UserService.Verify(requests)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "User has been successfully verified")
}

func (r *Rest) ChangePassword(ctx *gin.Context) {
	var requests model.RequestForChangePassword
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

	token := ctx.Param("token")
	if respDetails := r.service.UserService.ChangePassword(token, requests); respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "Success change password")
}

func (r *Rest) LoginAccount(ctx *gin.Context) {
	var requests model.RequestForLogin
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

	tokens, respDetails := r.service.UserService.Login(requests)
	if respDetails.Code == 403 {

		response.WithData(ctx, respDetails.Code, respDetails.Message, model.ResponseForRegister{
			ID: tokens.UserID,
		})
		return
	}

	if respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "User has successfully logged in", tokens)
}

func (r *Rest) MyData(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {

		response.MessageOnly(ctx, 500, "Failed to get data")
		return
	}

	response.WithData(ctx, 200, "Success get data", user.(entity.User))
	return
}