package service

import (
	"intern-bcc-2024/internal/repository"
)

type ITokenService interface {
	//
}

type TokenService struct {
	tr repository.ITokenRepository
}

func NewTokenService(tokenRepository repository.ITokenRepository) ITokenService {
	return &TokenService{
		tr: tokenRepository,
	}
}
