package service

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/pkg/database/mysql"
	"intern-bcc-2024/pkg/response"
	"log"
)

type IAccountService interface {
	AllAccountNumber() (*[]entity.AccountNumberType, response.Details)
}

type AccountService struct {
	db *gorm.DB
	ar repository.IAccountRepository
}

func NewAccountService(accountRepository repository.IAccountRepository) IAccountService {
	return &AccountService{
		db: mysql.Connection,
		ar: accountRepository,
	}
}

func (as *AccountService) AllAccountNumber() (*[]entity.AccountNumberType, response.Details) {
	accounts := new([]entity.AccountNumberType)

	tx := as.db.Begin()
	defer tx.Rollback()

	respDetails := as.ar.GetAllAccountTypes(tx, accounts)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return accounts, respDetails
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err)

		return accounts, response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return accounts, response.Details{Code: 200, Message: "Success get all account types", Error: nil}
}
