package middleware

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"intern-bcc-2024/pkg/response"
)

func (m *middleware) Timeout() gin.HandlerFunc {
	timeLimit, err := strconv.Atoi(os.Getenv("TIME_OUT_LIMIT"))
	if err != nil {
		return func(c *gin.Context) {
			response.MessageOnly(c, 500, "Failed convert TIME_OUT_LIMIT value")
			c.Abort()
			return
		}
	}

	return timeout.New(
		timeout.WithTimeout(time.Duration(timeLimit)*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse),
	)
}

func timeoutResponse(c *gin.Context) {
	response.MessageOnly(c, 408, "Please try again. The request takes to much time")
}
