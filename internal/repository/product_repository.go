package repository

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
)

type IProductRepository interface {
	Search(tx *gorm.DB, categoryDetails *entity.Category, products *[]model.ResponseForSearch, request model.RequestForSearch) response.Details
	Find(tx *gorm.DB, product *entity.Product, param model.ParamForFind) response.Details
	GetByID(tx *gorm.DB, product *model.ResponseForGetProductByID, productID uuid.UUID, user entity.User) response.Details
}

type ProductRepository struct {
	db *gorm.DB
	cr ICategoryRepository
}

func NewProductRepository(db *gorm.DB, cr ICategoryRepository) IProductRepository {
	return &ProductRepository{db, cr}
}

func (pr *ProductRepository) Search(tx *gorm.DB, categoryDetails *entity.Category, products *[]model.ResponseForSearch, request model.RequestForSearch) response.Details {
	query := fmt.Sprintf("SELECT users.id AS owner_id, users.name AS owner_name, products.id AS product_id, products.name AS product_name, products.price AS product_price, (SELECT media.url_media FROM media WHERE products.id = media.product_id LIMIT 1) AS url_photo_product, ROUND(ACOS(SIN(RADIANS(%.6f)) * SIN(RADIANS(latitude)) + COS(RADIANS(%.6f)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(%.6f))) * 6371, 1) AS owner_distance FROM products INNER JOIN users ON products.user_id = users.id ", request.Latitude, request.Latitude, request.Longitude)

	queryCategory := ""
	if request.Category != "" {
		respDetails := pr.cr.Find(tx, categoryDetails, model.ParamForFind{
			Name: request.Category,
		})
		if respDetails.Error != nil {
			return respDetails
		}
		queryCategory = fmt.Sprintf("AND products.category_id = '%s' ", categoryDetails.ID)
	}

	if request.Sort == "default" {
		request.Sort = "products.created_at DESC"
	} else if request.Sort == "distance" {
		request.Sort = "owner_distance"
	}

	query += fmt.Sprintf("WHERE products.user_id != '%s' AND products.name LIKE '%%%s%%' %sORDER BY %s LIMIT %d OFFSET %d", request.UserID, request.Query, queryCategory, request.Sort, request.Limit, request.Offset)
	if err := tx.Debug().Raw(query).Scan(products).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to find products", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to find products", Error: nil}
}

func (pr *ProductRepository) Find(tx *gorm.DB, product *entity.Product, param model.ParamForFind) response.Details {
	if err := tx.Debug().Where(&param).First(&product).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to find product", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to find product", Error: nil}
}

func (pr *ProductRepository) GetByID(tx *gorm.DB, product *model.ResponseForGetProductByID, productID uuid.UUID, user entity.User) response.Details {
	query := fmt.Sprintf("SELECT users.id AS owner_id, users.name AS owner_name, users.url_photo_profile AS owner_photo_profile, products.id AS product_id, products.name AS product_name, products.description AS product_description, products.price AS product_price, ROUND(ACOS(SIN(RADIANS(%.6f)) * SIN(RADIANS(latitude)) + COS(RADIANS(%.6f)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(%.6f))) * 6371,1) AS owner_distance FROM products INNER JOIN users ON products.user_id = users.id WHERE products.id = '%s'", user.Latitude, user.Latitude, user.Longitude, productID)
	if err := tx.Debug().Raw(query).Scan(product).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to find products", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to find product", Error: nil}
}
