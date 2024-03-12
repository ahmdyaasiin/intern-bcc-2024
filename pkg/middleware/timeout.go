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
	timeLimit, _ := strconv.Atoi(os.Getenv("TIME_OUT_LIMIT"))

	return timeout.New(
		timeout.WithTimeout(time.Duration(timeLimit)*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse),
	)
}

func timeoutResponse(c *gin.Context) {
	response.MessageOnly(c, 408, "Please try again. The request take to much time")
}
