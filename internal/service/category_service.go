package service

import (
	"gorm.io/gorm"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/database/mysql"
	"intern-bcc-2024/pkg/response"
	"log"
)

type ICategoryService interface {
	GetHomePage() (*[]model.ResponseForHomePage, response.Details)
}

type CategoryService struct {
	db *gorm.DB
	cr repository.ICategoryRepository
}

func NewCategoryService(categoryRepository repository.ICategoryRepository) ICategoryService {
	return &CategoryService{
		db: mysql.Connection,
		cr: categoryRepository,
	}
}

func (cs *CategoryService) GetHomePage() (*[]model.ResponseForHomePage, response.Details) {
	categories := new([]model.ResponseForHomePage)

	tx := cs.db.Begin()
	defer tx.Rollback()

	respDetails := cs.cr.GetAllCategories(tx, categories)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return categories, respDetails
	}

	if err := tx.Commit().Error; err != nil {
		log.Println(err)

		return categories, response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return categories, response.Details{Code: 200, Message: "Success get all categories", Error: nil}
}
