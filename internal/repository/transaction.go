package repository

import (
	"fmt"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/pkg/response"
)

type ITransactionRepository interface {
	CreateTransaction(transaction *entity.Transaction) response.Details
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &TransactionRepository{db}
}

func (tr *TransactionRepository) CreateTransaction(transaction *entity.Transaction) response.Details {
	if err := tr.db.Debug().Create(transaction).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to create transaction", Error: err}
	}

	fmt.Println("hoho")

	return response.Details{Code: 200, Message: "Success to create transaction", Error: nil}
}
