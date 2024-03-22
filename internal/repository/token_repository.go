package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type ITokenRepository interface {
	Find(tx *gorm.DB, token *entity.ResetToken, param model.ParamForFind) response.Details
	Create(tx *gorm.DB, token *entity.ResetToken) response.Details
	Delete(tx *gorm.DB, token *entity.ResetToken) response.Details
}

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) ITokenRepository {
	return &TokenRepository{db}
}

func (tr *TokenRepository) Find(tx *gorm.DB, token *entity.ResetToken, param model.ParamForFind) response.Details {
	if err := tx.Debug().Where(&param).First(&token).Error; err != nil {
		return response.Details{Code: 500, Message: "Token gagal ditemukan", Error: err}
	}

	return response.Details{Code: 200, Message: "Token berhasil ditemukan", Error: nil}
}

func (tr *TokenRepository) Create(tx *gorm.DB, token *entity.ResetToken) response.Details {
	if err := tx.Debug().Create(token).Error; err != nil {
		return response.Details{Code: 500, Message: "Token gagal dibuat", Error: err}
	}

	return response.Details{Code: 200, Message: "Token berhasil dibuat", Error: nil}
}

func (tr *TokenRepository) Delete(tx *gorm.DB, token *entity.ResetToken) response.Details {
	if err := tx.Debug().Delete(token).Error; err != nil {
		return response.Details{Code: 500, Message: "Token gagal dihapus", Error: err}
	}

	return response.Details{Code: 200, Message: "Token berhasil dihapus", Error: nil}
}
