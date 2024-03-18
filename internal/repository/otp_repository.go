package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type IOtpRepository interface {
	Find(tx *gorm.DB, otp *entity.OtpCode, param model.ParamForFind) response.Details
	Create(tx *gorm.DB, otp *entity.OtpCode) response.Details
	Update(tx *gorm.DB, otp *entity.OtpCode) response.Details
	Delete(tx *gorm.DB, otp *entity.OtpCode) response.Details
}

type OtpRepository struct {
	db *gorm.DB
}

func NewOtpRepository(db *gorm.DB) IOtpRepository {
	return &OtpRepository{db}
}

func (or *OtpRepository) Find(tx *gorm.DB, otp *entity.OtpCode, param model.ParamForFind) response.Details {
	if err := tx.Debug().Where(&param).First(&otp).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to find otp", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to find otp", Error: nil}
}

func (or *OtpRepository) Create(tx *gorm.DB, otp *entity.OtpCode) response.Details {
	if err := tx.Debug().Create(otp).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to create otp", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to create user", Error: nil}
}

func (or *OtpRepository) Update(tx *gorm.DB, otp *entity.OtpCode) response.Details {
	if err := tx.Debug().Updates(otp).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to update otp", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to update otp", Error: nil}
}

func (or *OtpRepository) Delete(tx *gorm.DB, otp *entity.OtpCode) response.Details {
	if err := tx.Debug().Delete(otp).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to delete otp", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to delete otp", Error: nil}
}
