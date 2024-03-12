package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type ISessionRepository interface {
	Find(param model.ParamForFind) (entity.RefreshToken, response.Details)
	Create(session *entity.RefreshToken) response.Details
	Update(session *entity.RefreshToken) response.Details
	Delete(session *entity.RefreshToken) response.Details
}

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) ISessionRepository {
	return &SessionRepository{db}
}

func (sr *SessionRepository) Find(param model.ParamForFind) (entity.RefreshToken, response.Details) {
	refreshToken := entity.RefreshToken{}
	if err := sr.db.Debug().Where(&param).First(&refreshToken).Error; err != nil {
		return refreshToken, response.Details{Code: 500, Message: "Failed to find session", Error: err}
	}

	return refreshToken, response.Details{Code: 200, Message: "Success to find session", Error: nil}
}

func (sr *SessionRepository) Create(session *entity.RefreshToken) response.Details {
	if err := sr.db.Debug().Create(session).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to create session", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to create session", Error: nil}
}

func (sr *SessionRepository) Update(session *entity.RefreshToken) response.Details {
	if err := sr.db.Debug().Updates(session).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to update session", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to update session", Error: nil}
}

func (sr *SessionRepository) Delete(session *entity.RefreshToken) response.Details {
	if err := sr.db.Debug().Delete(session).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to delete session", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to delete session", Error: nil}
}
