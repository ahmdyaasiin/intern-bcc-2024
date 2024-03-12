package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/jwt"
	"intern-bcc-2024/pkg/response"
	"time"
)

type ISessionService interface {
	Renew(requests model.RequestForRenewAccessToken) (model.ResponseForRenew, response.Details)
	Logout(ctx *gin.Context) response.Details
}

type SessionService struct {
	sr      repository.ISessionRepository
	jwtAuth jwt.Interface
}

func NewSessionService(sessionRepository repository.ISessionRepository, jwtAuth jwt.Interface) ISessionService {
	return &SessionService{
		sr:      sessionRepository,
		jwtAuth: jwtAuth,
	}
}

func (ss *SessionService) Renew(requests model.RequestForRenewAccessToken) (model.ResponseForRenew, response.Details) {
	res := model.ResponseForRenew{}

	refreshToken, respDetails := ss.sr.Find(model.ParamForFind{
		Token: requests.RefreshToken,
	})
	if respDetails.Error != nil {
		return res, respDetails
	}

	if refreshToken.ExpiredAt < time.Now().Local().UnixMilli() {
		return res, response.Details{Code: 403, Message: "Refresh token is expired", Error: errors.New("refresh token expired")}
	}

	accessToken, err := ss.jwtAuth.CreateAccessToken(refreshToken.UserID)
	if err != nil {
		return res, response.Details{Code: 500, Message: "Failed to create access token", Error: err}
	}

	res.AccessToken = accessToken
	return res, response.Details{Code: 200, Message: "Success renew access token", Error: nil}
}

func (ss *SessionService) Logout(ctx *gin.Context) response.Details {
	user, err := ss.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		return response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	sToken, respDetails := ss.sr.Find(model.ParamForFind{
		UserID: user.ID,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	if err = ss.sr.Delete(&sToken).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to delete the session", Error: err}
	}

	return response.Details{Code: 200, Message: " Success logout", Error: nil}
}
