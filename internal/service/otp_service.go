package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/database/mysql"
	"intern-bcc-2024/pkg/mail"
	"intern-bcc-2024/pkg/response"
	"log"
	"time"
)

type IOtpService interface {
	Resend(requests model.RequestForResend) response.Details
}

type OtpService struct {
	db *gorm.DB
	or repository.IOtpRepository
	ur repository.IUserRepository
}

func NewOtpService(otpRepository repository.IOtpRepository, userRepository repository.IUserRepository) IOtpService {
	return &OtpService{
		db: mysql.Connection,
		or: otpRepository,
		ur: userRepository,
	}
}

func (cs *OtpService) Resend(requests model.RequestForResend) response.Details {
	user := new(entity.User)
	otp := new(entity.OtpCode)

	tx := cs.db.Begin()
	defer tx.Rollback()

	respDetails := cs.ur.Find(tx, user, model.ParamForFind{
		ID: requests.ID,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if user.StatusAccount == "active" {
		log.Println("User already verified")

		return response.Details{Code: 403, Message: "User already verified", Error: errors.New("user already verified")}
	}

	respDetails = cs.or.Find(tx, otp, model.ParamForFind{
		UserID: requests.ID,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if time.Now().Local().Add(-1*5*time.Minute).UnixMilli() <= otp.CreatedAt {
		log.Println("Limit send otp")

		return response.Details{Code: 403, Message: "Please wait 5 minutes again", Error: errors.New("limit to send link reset")}
	}

	otp.Code = mail.GenerateSixCode()
	if err := mail.SendEmail(user.Email, fmt.Sprintf("Verification Code", "Your Verification Code: %s"), otp.Code); err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to resend verification code to user", Error: err}
	}

	respDetails = cs.or.Update(tx, otp)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to resend verification code to user", Error: nil}
}
