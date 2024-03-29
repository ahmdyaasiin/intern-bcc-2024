package middleware

import (
	"github.com/gin-gonic/gin"
	"intern-bcc-2024/internal/service"
	"intern-bcc-2024/pkg/jwt"
)

type Interface interface {
	Auth(ctx *gin.Context)
	Cors() gin.HandlerFunc
	Timeout() gin.HandlerFunc
}

type middleware struct {
	jwtAuth jwt.Interface
	service *service.Service
}

func Init(jwtAuth jwt.Interface, service *service.Service) Interface {
	return &middleware{
		jwtAuth: jwtAuth,
		service: service,
	}
}
