package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type ISessionRepository interface {
	Find(tx *gorm.DB, token *entity.Session, param model.ParamForFind) response.Details
	Create(tx *gorm.DB, session *entity.Session) response.Details
	Update(session *entity.Session) response.Details
	Delete(tx *gorm.DB, session *entity.Session) response.Details
}

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) ISessionRepository {
	return &SessionRepository{db}
}

func (sr *SessionRepository) Find(tx *gorm.DB, token *entity.Session, param model.ParamForFind) response.Details {
	if err := tx.Debug().Where(&param).First(token).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to find session", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to find session", Error: nil}
}

func (sr *SessionRepository) Create(tx *gorm.DB, session *entity.Session) response.Details {
	if err := tx.Debug().Create(session).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to create session", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to create session", Error: nil}
}

func (sr *SessionRepository) Update(session *entity.Session) response.Details {
	if err := sr.db.Debug().Updates(session).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to update session", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to update session", Error: nil}
}

func (sr *SessionRepository) Delete(tx *gorm.DB, session *entity.Session) response.Details {
	if err := tx.Debug().Delete(session).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to delete session", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to delete session", Error: nil}
}
