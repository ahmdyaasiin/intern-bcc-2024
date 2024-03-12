package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository    IUserRepository
	OtpRepository     IOtpRepository
	TokenRepository   ITokenRepository
	SessionRepository ISessionRepository
}

func NewRepository(db *gorm.DB) *Repository {
	or := NewOtpRepository(db)
	ur := NewUserRepository(db, or)
	tr := NewTokenRepository(db)
	sr := NewSessionRepository(db)

	return &Repository{
		OtpRepository:     or,
		UserRepository:    ur,
		TokenRepository:   tr,
		SessionRepository: sr,
	}
}
