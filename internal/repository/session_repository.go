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
		return response.Details{Code: 500, Message: "Session gagal ditemukan", Error: err}
	}

	return response.Details{Code: 200, Message: "Session berhasil ditemukan", Error: nil}
}

func (sr *SessionRepository) Create(tx *gorm.DB, session *entity.Session) response.Details {
	if err := tx.Debug().Create(session).Error; err != nil {
		return response.Details{Code: 500, Message: "Session gagal dibuat", Error: err}
	}

	return response.Details{Code: 200, Message: "Session berhasil dibuat", Error: nil}
}

func (sr *SessionRepository) Delete(tx *gorm.DB, session *entity.Session) response.Details {
	if err := tx.Debug().Delete(session).Error; err != nil {
		return response.Details{Code: 500, Message: "Session gagal dihapus", Error: err}
	}

	return response.Details{Code: 200, Message: "Session berhasil dihapus", Error: nil}
}
