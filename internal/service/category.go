package service

import (
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type ICategoryService interface {
	GetHomePage() ([]*model.ResponseForHomePage, response.Details)
}

type CategoryService struct {
	cr repository.ICategoryRepository
}

func NewCategoryService(categoryRepository repository.ICategoryRepository) ICategoryService {
	return &CategoryService{
		cr: categoryRepository,
	}
}

func (cs *CategoryService) GetHomePage() ([]*model.ResponseForHomePage, response.Details) {
	categories, respDetails := cs.cr.GetAllCategories()
	if respDetails.Error != nil {
		return categories, respDetails
	}

	return categories, response.Details{Code: 200, Message: "Success get all categories", Error: nil}
}
