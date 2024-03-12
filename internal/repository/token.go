package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type ITokenRepository interface {
	Find(param model.ParamForFind) (entity.ResetToken, response.Details)
	Create(token *entity.ResetToken) response.Details
	Update(token *entity.ResetToken) response.Details
	Delete(token *entity.ResetToken) response.Details
}

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) ITokenRepository {
	return &TokenRepository{db}
}

func (tr *TokenRepository) Find(param model.ParamForFind) (entity.ResetToken, response.Details) {
	resetToken := entity.ResetToken{}
	if err := tr.db.Debug().Where(&param).First(&resetToken).Error; err != nil {
		return resetToken, response.Details{Code: 500, Message: "Failed to find token", Error: err}
	}

	return resetToken, response.Details{Code: 200, Message: "Success to find token", Error: nil}
}

func (tr *TokenRepository) Create(token *entity.ResetToken) response.Details {
	if err := tr.db.Debug().Create(token).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to create token", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to create token", Error: nil}
}

func (tr *TokenRepository) Update(token *entity.ResetToken) response.Details {
	if err := tr.db.Debug().Updates(token).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to update token", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to update token", Error: nil}
}

func (tr *TokenRepository) Delete(token *entity.ResetToken) response.Details {
	if err := tr.db.Debug().Delete(token).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to delete token", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to delete token", Error: nil}
}
