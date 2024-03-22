package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type ICategoryRepository interface {
	GetAllCategories(tx *gorm.DB, categories *[]model.ResponseForHomePage) response.Details
	Find(tx *gorm.DB, category *entity.Category, param model.ParamForFind) response.Details
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{db}
}

func (cr *CategoryRepository) Find(tx *gorm.DB, category *entity.Category, param model.ParamForFind) response.Details {
	if err := tx.Debug().Where(&param).First(category).Error; err != nil {
		return response.Details{Code: 500, Message: "Kategori gagal ditemukan", Error: err}
	}

	return response.Details{Code: 200, Message: "Kategori berhasil ditemukan", Error: nil}
}

func (cr *CategoryRepository) GetAllCategories(tx *gorm.DB, categories *[]model.ResponseForHomePage) response.Details {
	if err := tx.Debug().Raw("SELECT categories.*, COUNT(products.id) AS total_product FROM categories LEFT JOIN products ON categories.id = products.category_id GROUP BY categories.id").Scan(categories).Error; err != nil {
		return response.Details{Code: 500, Message: "Kategori gagal ditemukan", Error: err}
	}

	return response.Details{Code: 200, Message: "Kategori berhasil ditemukan", Error: nil}
}
