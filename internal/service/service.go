package service

import (
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/pkg/bcrypt"
	"intern-bcc-2024/pkg/jwt"
)

type Service struct {
	UserService        IUserService
	OtpService         IOtpService
	TokenService       ITokenService
	SessionService     ISessionService
	ProductService     IProductService
	CategoryService    ICategoryService
	TransactionService ITransactionService
}

type InitParam struct {
	Repository *repository.Repository
	Bcrypt     bcrypt.Interface
	JwtAuth    jwt.Interface
}

func NewService(param InitParam) *Service {
	userService := NewUserService(param.Repository.UserRepository, param.Repository.OtpRepository, param.Repository.TokenRepository, param.Repository.SessionRepository, param.Bcrypt, param.JwtAuth)
	otpService := NewOtpService(param.Repository.OtpRepository, param.Repository.UserRepository)
	sessionService := NewSessionService(param.Repository.SessionRepository, param.JwtAuth)
	tokenService := NewTokenService(param.Repository.TokenRepository, param.Repository.UserRepository)
	productService := NewProductService(param.Repository.ProductRepository, param.Repository.MediaRepository, param.Repository.UserRepository, param.JwtAuth)
	categoryService := NewCategoryService(param.Repository.CategoryRepository)
	transactionService := NewTransactionService(param.Repository.ProductRepository, param.Repository.TransactionRepository, param.JwtAuth)

	return &Service{
		UserService:        userService,
		OtpService:         otpService,
		TokenService:       tokenService,
		SessionService:     sessionService,
		ProductService:     productService,
		CategoryService:    categoryService,
		TransactionService: transactionService,
	}
}
