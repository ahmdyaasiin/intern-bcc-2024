package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type IUserRepository interface {
	Find(tx *gorm.DB, user *entity.User, param model.ParamForFind) response.Details
	Create(tx *gorm.DB, user *entity.User) response.Details
	Update(tx *gorm.DB, user *entity.User) response.Details
	Delete(user *entity.User) response.Details
	Verify(user entity.User, otp entity.OtpCode) response.Details
	Change(user entity.User, token entity.ResetToken) response.Details
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) Find(tx *gorm.DB, user *entity.User, param model.ParamForFind) response.Details {
	if err := tx.Debug().Where(&param).First(user).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to find user", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to find user", Error: nil}
}

func (ur *UserRepository) Create(tx *gorm.DB, user *entity.User) response.Details {
	if err := tx.Debug().Create(user).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to create user", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to create user", Error: nil}
}

func (ur *UserRepository) Update(tx *gorm.DB, user *entity.User) response.Details {
	if err := tx.Debug().Updates(user).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to update user", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to update user", Error: nil}
}

func (ur *UserRepository) Delete(user *entity.User) response.Details {
	if err := ur.db.Debug().Delete(user).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to delete user", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to delete user", Error: nil}
}

func (ur *UserRepository) Verify(user entity.User, otp entity.OtpCode) response.Details {
	tx := ur.db.Begin()

	if err := tx.Debug().Updates(user).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to update user status", Error: err}
	}

	if err := tx.Debug().Delete(otp).Error; err != nil {
		tx.Rollback()
		return response.Details{Code: 500, Message: "Failed to delete OTP", Error: err}
	}

	tx.Commit()

	return response.Details{Code: 200, Message: "Success to verify user", Error: nil}
}

func (ur *UserRepository) Change(user entity.User, token entity.ResetToken) response.Details {
	tx := ur.db.Begin()

	if err := tx.Debug().Updates(user).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to update user status", Error: err}
	}

	if err := tx.Debug().Delete(token).Error; err != nil {
		tx.Rollback()
		return response.Details{Code: 500, Message: "Failed to delete Token", Error: err}
	}

	tx.Commit()

	return response.Details{Code: 200, Message: "Success to change password", Error: nil}
}
