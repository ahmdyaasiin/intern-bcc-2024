package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

func (m *middleware) Auth(ctx *gin.Context) {
	bearer := ctx.GetHeader("Authorization")
	if bearer == "" {
		response.MessageOnly(ctx, 401, "Authorization header is missing. Please provide valid authentication credentials")
		ctx.Abort()
		return
	}

	expired := false
	tokenSplit := strings.Split(bearer, " ")
	if len(tokenSplit) <= 1 {
		response.MessageOnly(ctx, 401, "Authorization header is missing. Please provide valid authentication credentials")
		ctx.Abort()
		return
	}

	token := tokenSplit[1]

	userId, err := m.jwtAuth.ValidateAccessToken(token)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			expired = true
		}

		response.WithExpired(ctx, 401, "Unauthorized", expired)
		ctx.Abort()
		return
	}

	user, respDetails := m.service.UserService.Find(model.ParamForFind{
		ID: userId,
	})
	if respDetails.Error != nil {
		response.MessageOnly(ctx, 401, "failed get user")
		ctx.Abort()
		return
	}

	ctx.Set("user", *user)
	ctx.Next()
}
