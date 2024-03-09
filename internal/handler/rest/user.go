package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
	"intern-bcc-2024/pkg/validation"
)

func (r *Rest) Register(ctx *gin.Context) {
	var requests model.RequestForRegister
	if err := ctx.ShouldBindJSON(&requests); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[validation.GetField(fe)] = validation.GetErrorMsg(fe)
			}

			response.ValidationError(ctx, 422, "Validasi error", out)
			return
		}

		response.Message(ctx, 422, "Failed to bind input")
		return
	}

	user, err := r.service.UserService.Register(requests)
	if err != nil {
		response.Error(ctx, 500, "An unexpected error occurred", err)
		return
	}

	response.Success(ctx, 201, "User has been successfully registered", user)

}

func (r *Rest) Verify(ctx *gin.Context) {
	var requests model.OtpParam
	if err := ctx.ShouldBindJSON(&requests); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[validation.GetField(fe)] = validation.GetErrorMsg(fe)
			}

			response.ValidationError(ctx, 422, "Validasi error", out)
			return
		}

		response.Message(ctx, 422, "Failed to bind input")
		return
	}

	if err := r.service.UserService.Verify(requests); err != nil {
		response.Error(ctx, 500, "An unexpected error occurred", err)
		return
	}

	response.Message(ctx, 200, "User has been successfully verified")
}

func (r *Rest) Resend(ctx *gin.Context) {
	var requests model.RequestForResend
	if err := ctx.ShouldBindJSON(&requests); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[validation.GetField(fe)] = validation.GetErrorMsg(fe)
			}

			response.ValidationError(ctx, 422, "Validasi error", out)
			return
		}

		response.Message(ctx, 422, "Failed to bind input")
		return
	}

	if err := r.service.UserService.Resend(requests); err != nil {
		response.Error(ctx, 500, "An unexpected error occurred", err)
		return
	}

	response.Message(ctx, 200, "The OTP code was sent successfully")
}

func (r *Rest) Login(ctx *gin.Context) {
	var requests model.RequestForLogin
	if err := ctx.ShouldBindJSON(&requests); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[validation.GetField(fe)] = validation.GetErrorMsg(fe)
			}

			response.ValidationError(ctx, 422, "Validasi error", out)
			return
		}

		response.Message(ctx, 422, "Failed to bind input")
		return
	}

	tokens, err := r.service.UserService.Login(requests)
	if err != nil {
		response.Error(ctx, 500, "An unexpected error occurred", err)
		return
	}

	response.Success(ctx, 200, "User has successfully logged in", tokens)
}

func (r *Rest) Renew(ctx *gin.Context) {
	var requests model.RequestForRenew
	if err := ctx.ShouldBindJSON(&requests); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[validation.GetField(fe)] = validation.GetErrorMsg(fe)
			}

			response.ValidationError(ctx, 422, "Validasi error", out)
			return
		}

		response.Message(ctx, 422, "Failed to bind input")
		return
	}

	resp, err := r.service.UserService.Renew(requests)
	if err != nil {
		response.Error(ctx, 500, "An unexpected error occurred", err)
		return
	}

	response.Success(ctx, 200, "Successfully updated access token", resp)
}

func (r *Rest) Reset(ctx *gin.Context) {
	var requests model.RequestForReset
	if err := ctx.ShouldBindJSON(&requests); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[validation.GetField(fe)] = validation.GetErrorMsg(fe)
			}

			response.ValidationError(ctx, 422, "Validasi error", out)
			return
		}

		response.Message(ctx, 422, "Failed to bind input")
		return
	}

	if err := r.service.UserService.Reset(requests); err != nil {

		response.Error(ctx, 500, "An unexpected error occurred", err)
		return
	}

	response.Message(ctx, 200, "Successfully sent the password reset link")
}

func (r *Rest) ResetGet(ctx *gin.Context) {
	token := ctx.Param("token")
	if err := r.service.UserService.ResetGet(token); err != nil {

		response.Error(ctx, 500, "An unexpected error occurred", err)
		return
	}

	response.Message(ctx, 200, "Token is valid")
}

func (r *Rest) ResetPost(ctx *gin.Context) {
	var requests model.RequestForChangePassword
	if err := ctx.ShouldBindJSON(&requests); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[validation.GetField(fe)] = validation.GetErrorMsg(fe)
			}

			response.ValidationError(ctx, 422, "Validasi error", out)
			return
		}

		response.Message(ctx, 422, "Failed to bind input")
		return
	}

	token := ctx.Param("token")
	if err := r.service.UserService.ResetPost(requests, token); err != nil {
		response.Message(ctx, 500, "failed to change password")
		return
	}

	response.Message(ctx, 200, "Token is valid")
}
