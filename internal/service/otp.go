package service

import (
	"intern-bcc-2024/internal/repository"
)

type IOtpService interface {
	//
}

type OtpService struct {
	or repository.IOtpRepository
}

func NewOtpService(otpRepository repository.IOtpRepository) IOtpService {
	return &OtpService{
		or: otpRepository,
	}
}
