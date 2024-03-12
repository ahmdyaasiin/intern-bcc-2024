package response

import (
	"github.com/gin-gonic/gin"
)

type Details struct {
	Code    int
	Message string
	Error   error
}

type Success struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Message struct {
	Message string `json:"message"`
}

type Expired struct {
	Message   string `json:"message"`
	IsExpired bool   `json:"is_expired"`
}

type ValidationError struct {
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}

func WithExpired(ctx *gin.Context, httpStatusCode int, message string, expired bool) {
	ctx.JSON(httpStatusCode, Expired{
		Message:   message,
		IsExpired: expired,
	})
}

func WithData(ctx *gin.Context, httpStatusCode int, message string, data any) {
	ctx.JSON(httpStatusCode, Success{
		Message: message,
		Data:    data,
	})
}

func MessageOnly(ctx *gin.Context, httpStatusCode int, message string) {
	ctx.JSON(httpStatusCode, Message{
		Message: message,
	})
}

func WithErrors(ctx *gin.Context, httpStatusCode int, message string, errors any) {
	ctx.JSON(httpStatusCode, ValidationError{
		Message: message,
		Errors:  errors,
	})
}
