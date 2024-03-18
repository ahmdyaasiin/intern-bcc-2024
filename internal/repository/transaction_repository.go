package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/pkg/response"
)

type ITransactionRepository interface {
	CreateTransaction(tx *gorm.DB, transaction *entity.Transaction) response.Details
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &TransactionRepository{db}
}

func (tr *TransactionRepository) CreateTransaction(tx *gorm.DB, transaction *entity.Transaction) response.Details {
	if err := tx.Debug().Create(transaction).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to create transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to create transaction", Error: nil}
}
