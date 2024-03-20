package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/database/mysql"
	"intern-bcc-2024/pkg/mail"
	"intern-bcc-2024/pkg/response"
	"log"
	"os"
	"time"
)

type ITokenService interface {
	ResetPassword(requests model.RequestForReset) response.Details
	CheckResetToken(token string) response.Details
}

type TokenService struct {
	db *gorm.DB
	tr repository.ITokenRepository
	ur repository.IUserRepository
}

func NewTokenService(tokenRepository repository.ITokenRepository, userRepository repository.IUserRepository) ITokenService {
	return &TokenService{
		db: mysql.Connection,
		tr: tokenRepository,
		ur: userRepository,
	}
}

func (ts *TokenService) ResetPassword(requests model.RequestForReset) response.Details {
	user := new(entity.User)
	token := new(entity.ResetToken)

	tx := ts.db.Begin()
	defer tx.Rollback()

	respDetails := ts.ur.Find(tx, user, model.ParamForFind{
		Email: requests.Email,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	respDetails = ts.tr.Find(tx, token, model.ParamForFind{
		UserID: user.ID,
	})
	if respDetails.Error == nil {
		if time.Now().Local().Add(-1*5*time.Minute).UnixMilli() <= token.CreatedAt {
			log.Println(respDetails.Error)

			return response.Details{Code: 403, Message: "Please wait 5 minutes again", Error: errors.New("limit to send link reset")}
		}

		respDetails = ts.tr.Delete(tx, token)
		if respDetails.Error != nil {
			log.Println(respDetails.Error)

			return respDetails
		}
	}

	link := mail.GenerateRandomString(30)
	if err := mail.SendEmail(user.Email, "Link Reset Password", fmt.Sprintf("Your link: %s/resetconfirm/%s", os.Getenv("LINK_FRONTEND"), link)); err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to send link reset password", Error: err}
	}

	token = &entity.ResetToken{
		ID:     uuid.New(),
		UserID: user.ID,
		Token:  link,
	}

	respDetails = ts.tr.Create(tx, token)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success send link reset password", Error: nil}
}

func (ts *TokenService) CheckResetToken(tokenRequest string) response.Details {
	token := new(entity.ResetToken)

	tx := ts.db.Begin()
	defer tx.Rollback()

	if respDetails := ts.tr.Find(tx, token, model.ParamForFind{
		Token: tokenRequest,
	}); respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Token is valid", Error: nil}
}
