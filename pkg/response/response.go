package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Message string `json:"message"`
	Error   any    `json:"error"`
}

type ResponseValidationError struct {
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}

func Success(ctx *gin.Context, httpStatusCode int, message string, data any) {
	ctx.JSON(httpStatusCode, Response{
		Message: message,
		Data:    data,
	})
}

func Message(ctx *gin.Context, httpStatusCode int, message string) {
	ctx.JSON(httpStatusCode, ResponseMessage{
		Message: message,
	})
}

func Error(ctx *gin.Context, httpStatusCode int, message string, err error) {
	ctx.JSON(httpStatusCode, ResponseError{
		Message: message,
		Error:   err.Error(),
	})
}

func ValidationError(ctx *gin.Context, httpStatusCode int, message string, errors any) {
	ctx.JSON(httpStatusCode, ResponseValidationError{
		Message: message,
		Errors:  errors,
	})
}
