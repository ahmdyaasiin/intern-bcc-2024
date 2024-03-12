package service

import (
	"intern-bcc-2024/internal/repository"
)

type ISessionService interface {
	//
}

type SessionService struct {
	tr repository.ISessionRepository
}

func NewSessionService(userRepository repository.ISessionRepository) ISessionService {
	return &SessionService{
		tr: userRepository,
	}
}
