package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/database/mysql"
	"intern-bcc-2024/pkg/jwt"
	"intern-bcc-2024/pkg/response"
	"log"
	"time"
)

type ISessionService interface {
	RenewSession(requests model.RequestForRenewAccessToken) (*model.ResponseForRenew, response.Details)
	Logout(ctx *gin.Context) response.Details
}

type SessionService struct {
	db      *gorm.DB
	sr      repository.ISessionRepository
	jwtAuth jwt.Interface
}

func NewSessionService(sessionRepository repository.ISessionRepository, jwtAuth jwt.Interface) ISessionService {
	return &SessionService{
		db:      mysql.Connection,
		sr:      sessionRepository,
		jwtAuth: jwtAuth,
	}
}

func (ss *SessionService) RenewSession(requests model.RequestForRenewAccessToken) (*model.ResponseForRenew, response.Details) {
	session := new(entity.Session)
	res := new(model.ResponseForRenew)

	tx := ss.db.Begin()
	defer tx.Rollback()

	respDetails := ss.sr.Find(tx, session, model.ParamForFind{
		Token: requests.RefreshToken,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return res, respDetails
	}

	if session.ExpiredAt < time.Now().Local().UnixMilli() {
		log.Println("Refresh token is expired")

		return res, response.Details{Code: 403, Message: "Refresh token is expired", Error: errors.New("refresh token expired")}
	}

	accessToken, err := ss.jwtAuth.CreateAccessToken(session.UserID)
	if err != nil {
		log.Println(err)

		return res, response.Details{Code: 500, Message: "Failed to create access token", Error: err}
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return res, response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	res.AccessToken = accessToken
	return res, response.Details{Code: 200, Message: "Success renew access token", Error: nil}
}

func (ss *SessionService) Logout(ctx *gin.Context) response.Details {
	session := new(entity.Session)

	tx := ss.db.Begin()
	defer tx.Rollback()

	user, err := ss.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	respDetails := ss.sr.Find(tx, session, model.ParamForFind{
		UserID: user.ID,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if err = ss.sr.Delete(tx, session).Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to delete the session", Error: err}
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: " Success logout", Error: nil}
}
