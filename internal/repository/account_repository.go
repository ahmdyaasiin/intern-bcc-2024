package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type IAccountRepository interface {
	GetAllAccountTypes(tx *gorm.DB, accounts *[]entity.AccountNumberType) response.Details
	Find(tx *gorm.DB, account *entity.AccountNumberType, param model.ParamForFind) response.Details
}

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) IAccountRepository {
	return &AccountRepository{db}
}

func (ar *AccountRepository) Find(tx *gorm.DB, account *entity.AccountNumberType, param model.ParamForFind) response.Details {
	if err := tx.Debug().Where(&param).First(account).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to find account number type", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to find account number type", Error: nil}
}

func (ar *AccountRepository) GetAllAccountTypes(tx *gorm.DB, accounts *[]entity.AccountNumberType) response.Details {
	if err := tx.Debug().Where("id != '00000000-0000-0000-0000-000000000000'").Find(accounts).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to find account types", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to find account types", Error: nil}
}
