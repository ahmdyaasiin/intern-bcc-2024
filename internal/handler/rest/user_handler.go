package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
	"intern-bcc-2024/pkg/validation"
	"log"
	"strings"
)

func (r *Rest) Register(ctx *gin.Context) {
	var requests model.RequestForRegister
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

	userID, respDetails := r.service.UserService.Register(requests)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 201, "User has been successfully registered", userID)

}

func (r *Rest) VerifyAfterRegister(ctx *gin.Context) {
	var requests model.RequestForVerify
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

	respDetails := r.service.UserService.VerifyAfterRegister(requests)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "User has been successfully verified")
}

func (r *Rest) ChangePasswordFromReset(ctx *gin.Context) {
	var requests model.RequestForChangePassword

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

	token := ctx.Param("token")
	if respDetails := r.service.UserService.ChangePasswordFromReset(token, requests); respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "Success change password")
}

func (r *Rest) Login(ctx *gin.Context) {
	var requests model.RequestForLogin

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

	tokens, respDetails := r.service.UserService.Login(requests)
	if strings.Contains(respDetails.Message, "verify") {
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

	response.WithData(ctx, 200, "Success to get data", user.(entity.User))
}

func (r *Rest) UpdateAccountNumber(ctx *gin.Context) {
	var requests model.RequestUpdateAccountNumber
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

	respDetails := r.service.UserService.UpdateAccountNumber(ctx, requests)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 201, "Success update")
}

func (r *Rest) GetName(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.MessageOnly(ctx, 422, "Failed convert id")
		return
	}

	user, respDetails := r.service.UserService.Find(model.ParamForFind{
		ID: id,
	})
	if respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "Success get product", model.ResponseForChatName{
		Name: user.Name,
	})
}
