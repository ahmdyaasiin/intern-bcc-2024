package service

import (
	"intern-bcc-2024/internal/repository"
)

type IProductService interface {
	//
}

type ProductService struct {
	pr repository.IProductRepository
}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &ProductService{
		pr: productRepository,
	}
}
