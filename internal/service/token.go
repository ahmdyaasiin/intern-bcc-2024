package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/mail"
	"intern-bcc-2024/pkg/response"
	"os"
	"time"
)

type ITokenService interface {
	ResetPassword(requests model.RequestForReset) response.Details
	CheckToken(token string) response.Details
}

type TokenService struct {
	tr repository.ITokenRepository
	ur repository.IUserRepository
}

func NewTokenService(tokenRepository repository.ITokenRepository, userRepository repository.IUserRepository) ITokenService {
	return &TokenService{
		tr: tokenRepository,
		ur: userRepository,
	}
}

func (ts *TokenService) ResetPassword(requests model.RequestForReset) response.Details {
	user, respDetails := ts.ur.Find(model.ParamForFind{
		Email: requests.Email,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	token, respDetails := ts.tr.Find(model.ParamForFind{
		UserID: user.ID,
	})
	if respDetails.Error == nil {
		if token.LastSentAt > time.Now().Local().Add(-5*time.Minute).UnixMilli() {
			return response.Details{Code: 403, Message: "Please wait 5 minutes again", Error: errors.New("limit to send link reset")} // add exact time
		}

		respDetails = ts.tr.Delete(&token)
		if respDetails.Error != nil {
			return respDetails
		}
	}

	link := mail.GenerateRandomString(30)
	if err := mail.SendEmail(user.Email, "Link Reset Password", "Your link: "+os.Getenv("LINK_FRONTEND")+"/resetconfirm/"+link); err != nil {
		fmt.Println(err)
		return response.Details{Code: 500, Message: "Failed to send link reset password", Error: err}
	}

	token = entity.ResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     link,
		ExpiredAt: time.Now().Local().Add(time.Hour * 1).UnixMilli(),
	}

	respDetails = ts.tr.Create(&token)
	if respDetails.Error != nil {
		return respDetails
	}

	return response.Details{Code: 200, Message: "Success send link reset password", Error: nil}
}

func (ts *TokenService) CheckToken(token string) response.Details {
	if _, respDetails := ts.tr.Find(model.ParamForFind{
		Token: token,
	}); respDetails.Error != nil {
		return respDetails
	}

	return response.Details{Code: 200, Message: "Token is valid", Error: nil}
}
