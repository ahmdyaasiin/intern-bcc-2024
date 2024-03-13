package repository

import (
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/pkg/response"
)

type IProductRepository interface {
	ProductHomePage(page int) ([]*entity.Product, response.Details)
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{db}
}

func (pr *ProductRepository) ProductHomePage(page int) ([]*entity.Product, response.Details) {
	//
	return []*entity.Product{}, response.Details{}
}
