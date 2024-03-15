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
	Search(request model.RequestForSearch, limit, offset int) (model.ResponseSearch, response.Details)
	Find(param model.ParamForFind) (entity.Product, response.Details)
	GetByID(productID uuid.UUID, user entity.User) (*model.ResponseForGetProductByID, response.Details)
}

type ProductRepository struct {
	db *gorm.DB
	cr ICategoryRepository
}

func NewProductRepository(db *gorm.DB, cr ICategoryRepository) IProductRepository {
	return &ProductRepository{db, cr}
}

func (pr *ProductRepository) Search(request model.RequestForSearch, limit, offset int) (model.ResponseSearch, response.Details) {
	var res []*model.ResponseForSearch
	query := fmt.Sprintf("SELECT users.id AS owner_id, users.name AS owner_name, products.id AS product_id, products.name AS product_name, products.price AS product_price, (SELECT media.url_media FROM media WHERE products.id = media.product_id LIMIT 1) AS url_photo_product, ROUND(ACOS(SIN(RADIANS(%.6f)) * SIN(RADIANS(latitude)) + COS(RADIANS(%.6f)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(%.6f))) * 6371, 1) AS owner_distance FROM products INNER JOIN users ON products.user_id = users.id ", request.Latitude, request.Latitude, request.Longitude)

	queryCategory := ""
	if request.Category != "" {
		categoryDetails, respDetails := pr.cr.Find(model.ParamForFind{
			Name: request.Category,
		})
		if respDetails.Error != nil {
			return model.ResponseSearch{}, respDetails
		}
		queryCategory = fmt.Sprintf("AND products.category_id = '%s' ", categoryDetails.ID)
	}

	if request.Sort == "default" {
		request.Sort = "products.created_at DESC"
	} else if request.Sort == "distance" {
		request.Sort = "owner_distance"
	}

	query += fmt.Sprintf("WHERE products.user_id != '%s' AND products.name LIKE '%%%s%%' %sORDER BY %s LIMIT %d OFFSET %d", request.UserID, request.Query, queryCategory, request.Sort, limit, offset)
	if err := pr.db.Debug().Raw(query).Scan(&res).Error; err != nil {
		return model.ResponseSearch{}, response.Details{Code: 500, Message: "Failed to find products", Error: err}
	}

	allCategories, respDetails := pr.cr.GetAllCategories()
	if respDetails.Error != nil {
		return model.ResponseSearch{}, respDetails
	}

	fmt.Println(query)
	return model.ResponseSearch{
		Product:  res,
		Category: allCategories,
	}, response.Details{Code: 200, Message: "Success to find products", Error: nil}
}

func (pr *ProductRepository) Find(param model.ParamForFind) (entity.Product, response.Details) {
	product := entity.Product{}
	if err := pr.db.Debug().Where(&param).First(&product).Error; err != nil {
		return product, response.Details{Code: 500, Message: "Failed to find product", Error: err}
	}

	return product, response.Details{Code: 200, Message: "Success to find product", Error: nil}
}

func (pr *ProductRepository) GetByID(productID uuid.UUID, user entity.User) (*model.ResponseForGetProductByID, response.Details) {
	var res model.ResponseForGetProductByID

	query := fmt.Sprintf("SELECT users.id AS owner_id, users.name AS owner_name, users.url_photo_profile AS owner_photo_profile, products.id AS product_id, products.name AS product_name, products.description AS product_description, products.price AS product_price, ROUND(ACOS(SIN(RADIANS(%.6f)) * SIN(RADIANS(latitude)) + COS(RADIANS(%.6f)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(%.6f))) * 6371,1) AS owner_distance FROM products INNER JOIN users ON products.user_id = users.id WHERE products.id = '%s'", user.Latitude, user.Latitude, user.Longitude, productID)
	if err := pr.db.Debug().Raw(query).Scan(&res).Error; err != nil {
		return nil, response.Details{Code: 500, Message: "Failed to find products", Error: err}
	}

	return &res, response.Details{Code: 200, Message: "Success to find product", Error: nil}
}
