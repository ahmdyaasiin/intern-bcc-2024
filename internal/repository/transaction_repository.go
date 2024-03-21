package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type ITransactionRepository interface {
	CreateTransaction(tx *gorm.DB, transaction *entity.Transaction) response.Details
	Find(tx *gorm.DB, transaction *entity.Transaction, param model.ParamForFind) response.Details
	Update(tx *gorm.DB, transaction *entity.Transaction) response.Details
	Delete(tx *gorm.DB, transaction *entity.Transaction) response.Details
	FindActiveTransactions(tx *gorm.DB, transaction *[]model.ResponseForActiveTransactions, user entity.User) response.Details
	BulkDelete(tx *gorm.DB, status string, createdAt int64) (int64, response.Details)
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

func (tr *TransactionRepository) Find(tx *gorm.DB, transaction *entity.Transaction, param model.ParamForFind) response.Details {
	if err := tx.Debug().Where(&param).First(transaction).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to find transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to find transaction", Error: nil}
}

func (tr *TransactionRepository) Update(tx *gorm.DB, transaction *entity.Transaction) response.Details {
	if err := tx.Debug().Updates(transaction).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to update transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to update transaction", Error: nil}
}

func (tr *TransactionRepository) Delete(tx *gorm.DB, transaction *entity.Transaction) response.Details {
	if err := tx.Debug().Delete(transaction).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to delete transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to delete transaction", Error: nil}
}

func (tr *TransactionRepository) FindActiveTransactions(tx *gorm.DB, transaction *[]model.ResponseForActiveTransactions, user entity.User) response.Details {
	if err := tx.Debug().Raw("SELECT transactions.id AS transaction_id, products.id as product_id, products.name AS product_name, products.price AS product_price, transactions.withdrawal_code AS withdrawal_code, transactions.user_id AS owner_id, (SELECT name FROM users WHERE users.id = owner_id) AS owner_name, (SELECT url FROM media WHERE media.product_id = products.id LIMIT 1) AS url_product FROM transactions INNER JOIN products ON transactions.product_id = products.id WHERE transactions.user_id = '" + user.ID.String() + "' AND transactions.status = 'paid'").Scan(transaction).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to get active transactions"}
	}

	return response.Details{Code: 200, Message: "Success to get active transactions"}
}

func (tr *TransactionRepository) BulkDelete(tx *gorm.DB, status string, createdAt int64) (int64, response.Details) {
	result := tx.Where("status = ? AND ? > created_at", status, createdAt).Delete(entity.Transaction{})

	if result.Error != nil {
		return 0, response.Details{Code: 500, Message: "Failed to delete expired transaction", Error: result.Error}
	}
	return result.RowsAffected, response.Details{Code: 200, Message: "Success to delete expired transaction", Error: nil}
}
