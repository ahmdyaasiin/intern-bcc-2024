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
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) Find(tx *gorm.DB, user *entity.User, param model.ParamForFind) response.Details {
	if err := tx.Debug().Where(param).First(user).Error; err != nil {
		return response.Details{Code: 500, Message: "User gagal ditemukan", Error: err}
	}

	return response.Details{Code: 200, Message: "User berhasil ditemukan", Error: nil}
}

func (ur *UserRepository) Create(tx *gorm.DB, user *entity.User) response.Details {
	if err := tx.Debug().Create(user).Error; err != nil {
		return response.Details{Code: 500, Message: "User gagal dibuat", Error: err}
	}

	return response.Details{Code: 200, Message: "User berhasil dibuat", Error: nil}
}

func (ur *UserRepository) Update(tx *gorm.DB, user *entity.User) response.Details {
	if err := tx.Debug().Where("id = ?", user.ID).Updates(&user).Error; err != nil {
		return response.Details{Code: 500, Message: "Data pengguna gagal diperbarui", Error: err}
	}

	return response.Details{Code: 200, Message: "Data pengguna berhasil diperbarui", Error: nil}
}
