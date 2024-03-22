package repository

import (
	"fmt"
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
		return response.Details{Code: 500, Message: "Transaksi gagal dibuat", Error: err}
	}

	return response.Details{Code: 200, Message: "Transaksi berhasil dibuat", Error: nil}
}

func (tr *TransactionRepository) Find(tx *gorm.DB, transaction *entity.Transaction, param model.ParamForFind) response.Details {
	if err := tx.Debug().Where(&param).First(transaction).Error; err != nil {
		return response.Details{Code: 500, Message: "Transaksi gagal ditemukan", Error: err}
	}

	return response.Details{Code: 200, Message: "Transaksi berhasil ditemukan", Error: nil}
}

func (tr *TransactionRepository) Update(tx *gorm.DB, transaction *entity.Transaction) response.Details {
	if err := tx.Debug().Updates(transaction).Error; err != nil {
		return response.Details{Code: 500, Message: "Transaksi gagal diperbarui", Error: err}
	}

	return response.Details{Code: 200, Message: "Transaksi berhasil diperbarui", Error: nil}
}

func (tr *TransactionRepository) Delete(tx *gorm.DB, transaction *entity.Transaction) response.Details {
	if err := tx.Debug().Delete(transaction).Error; err != nil {
		return response.Details{Code: 500, Message: "Transaksi gagal dihapus", Error: err}
	}

	return response.Details{Code: 200, Message: "Transaksi berhasil dihapus", Error: nil}
}

func (tr *TransactionRepository) FindActiveTransactions(tx *gorm.DB, transaction *[]model.ResponseForActiveTransactions, user entity.User) response.Details {
	query := fmt.Sprintf("SELECT transactions.id AS transaction_id, products.id AS product_id, products.name AS product_name, products.price AS product_price, transactions.withdrawal_code AS withdrawal_code, products.user_id AS owner_id, users.name AS owner_name, ( SELECT url FROM media WHERE product_id = products.id LIMIT 1 ) AS url_product FROM transactions INNER JOIN products ON transactions.product_id = products.id LEFT JOIN users ON products.user_id = users.id WHERE transactions.user_id = '%s' AND transactions.status = 'paid'", user.ID.String())
	if err := tx.Debug().Raw(query).Scan(transaction).Error; err != nil {
		return response.Details{Code: 500, Message: "Transaksi aktif gagal didapatkan"}
	}

	return response.Details{Code: 200, Message: "Transaksi aktif berhasil didapatkan"}
}

func (tr *TransactionRepository) BulkDelete(tx *gorm.DB, status string, createdAt int64) (int64, response.Details) {
	result := tx.Where("status = ? AND ? > created_at", status, createdAt).Delete(entity.Transaction{})
	if result.Error != nil {
		return 0, response.Details{Code: 500, Message: "Transaksi kadaluwarsa gagal dihapus", Error: result.Error}
	}

	return result.RowsAffected, response.Details{Code: 200, Message: "Transaksi kadaluwarsa berhasil dihapus", Error: nil}
}
