package service

import (
	"errors"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/mail"
	"intern-bcc-2024/pkg/response"
	"os"
	"strconv"
	"time"
)

type IOtpService interface {
	Resend(requests model.RequestForResend) response.Details
}

type OtpService struct {
	or repository.IOtpRepository
	ur repository.IUserRepository
}

func NewOtpService(otpRepository repository.IOtpRepository, userRepository repository.IUserRepository) IOtpService {
	return &OtpService{
		or: otpRepository,
		ur: userRepository,
	}
}

func (cs *OtpService) Resend(requests model.RequestForResend) response.Details {
	user, respDetails := cs.ur.Find(model.ParamForFind{
		ID: requests.ID,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	if user.StatusAccount == "active" {
		return response.Details{Code: 403, Message: "User already verified", Error: errors.New("user already verified")}
	}

	otp, respDetails := cs.or.Find(model.ParamForFind{
		UserID: requests.ID,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	expiredTime, err := strconv.Atoi(os.Getenv("EXPIRED_OTP"))
	if err != nil {
		return response.Details{Code: 500, Message: "Failed to convert expired time from .env", Error: err}
	}

	if otp.LastSentAt > time.Now().Local().Add(-5*time.Minute).UnixMilli() {
		return response.Details{Code: 403, Message: "Please wait 5 minutes again", Error: errors.New("limit to send otp")}
	}

	otp.VerificationCode = mail.GenerateVerificationCode()
	otp.ExpiredAt = time.Now().Local().Add(time.Duration(expiredTime) * time.Hour).UnixMilli()

	if err = mail.SendEmail(user.Email, "Verification Code", "Your Verification Code: "+otp.VerificationCode); err != nil {
		return response.Details{Code: 500, Message: "Failed to resend verification code to user", Error: err}
	}

	respDetails = cs.or.Update(&otp)
	if respDetails.Error != nil {
		return respDetails
	}

	return response.Details{Code: 200, Message: "Success to resend verification code to user", Error: nil}
}
