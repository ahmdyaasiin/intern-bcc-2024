package service

import (
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/pkg/bcrypt"
	"intern-bcc-2024/pkg/jwt"
)

type Service struct {
	UserService  IUserService
	OtpService   IOtpService
	TokenService ITokenService
}

type InitParam struct {
	Repository *repository.Repository
	Bcrypt     bcrypt.Interface
	JwtAuth    jwt.Interface
}

func NewService(param InitParam) *Service {
	userService := NewUserService(param.Repository.UserRepository, param.Repository.OtpRepository, param.Repository.TokenRepository, param.Repository.SessionRepository, param.Bcrypt, param.JwtAuth)
	otpService := NewOtpService(param.Repository.OtpRepository)
	tokenService := NewTokenService(param.Repository.TokenRepository)

	return &Service{
		UserService:  userService,
		OtpService:   otpService,
		TokenService: tokenService,
	}
}
