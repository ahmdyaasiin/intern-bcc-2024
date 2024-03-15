package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type ICategoryRepository interface {
	GetAllCategories() ([]*model.ResponseForHomePage, response.Details)
	Find(param model.ParamForFind) (entity.Category, response.Details)
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{db}
}

func (cr *CategoryRepository) Find(param model.ParamForFind) (entity.Category, response.Details) {
	category := entity.Category{}
	if err := cr.db.Debug().Where(&param).First(&category).Error; err != nil {
		return category, response.Details{Code: 500, Message: "Failed to find category", Error: err}
	}

	return category, response.Details{Code: 200, Message: "Success to find category", Error: nil}
}

func (cr *CategoryRepository) GetAllCategories() ([]*model.ResponseForHomePage, response.Details) {
	var categories []*model.ResponseForHomePage

	if err := cr.db.Debug().Raw("SELECT categories.*, COUNT(products.id) AS total_product FROM categories LEFT JOIN products ON categories.id = products.category_id GROUP BY categories.id").Scan(&categories).Error; err != nil {
		return nil, response.Details{Code: 500, Message: "Failed to get categories", Error: err}
	}

	return categories, response.Details{Code: 200, Message: "Success to get categories", Error: nil}
}
