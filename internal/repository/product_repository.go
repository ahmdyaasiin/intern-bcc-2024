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
	FindActiveProducts(tx *gorm.DB, product *[]model.ResponseForActiveProducts, user entity.User) response.Details
	Update(tx *gorm.DB, product *entity.Product) response.Details
	Delete(tx *gorm.DB, product *entity.Product) response.Details
	Create(tx *gorm.DB, product *entity.Product) response.Details
}

type ProductRepository struct {
	db *gorm.DB
	cr ICategoryRepository
}

func NewProductRepository(db *gorm.DB, cr ICategoryRepository) IProductRepository {
	return &ProductRepository{db, cr}
}

func (pr *ProductRepository) Search(tx *gorm.DB, categoryDetails *entity.Category, products *[]model.ResponseForSearch, request model.RequestForSearch) response.Details {
	query := fmt.Sprintf("SELECT users.id AS owner_id, users.name AS owner_name, products.id AS product_id, products.name AS product_name, products.price AS product_price, (SELECT url FROM media WHERE products.id = media.product_id LIMIT 1) AS url_photo_product, ACOS(SIN(RADIANS(%.6f)) * SIN(RADIANS(users.latitude)) + COS(RADIANS(%.6f)) * COS(RADIANS(users.latitude)) * COS(RADIANS(users.longitude) - RADIANS(%.6f))) * 6371000 AS owner_distance FROM products INNER JOIN users ON products.user_id = users.id LEFT JOIN transactions ON products.id = transactions.product_id ", request.Latitude, request.Latitude, request.Longitude)

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

	query += fmt.Sprintf("WHERE products.user_id != '%s' AND products.name LIKE '%%%s%%' AND transactions.product_id IS NULL %sORDER BY %s LIMIT %d OFFSET %d", request.UserID, request.Query, queryCategory, request.Sort, request.Limit, request.Offset)
	if err := tx.Debug().Raw(query).Scan(products).Error; err != nil {
		return response.Details{Code: 500, Message: "Produk gagal ditemukan", Error: err}
	}

	return response.Details{Code: 200, Message: "Produk berhasil ditemukan", Error: nil}
}

func (pr *ProductRepository) Find(tx *gorm.DB, product *entity.Product, param model.ParamForFind) response.Details {
	if err := tx.Debug().Where(&param).First(&product).Error; err != nil {
		return response.Details{Code: 500, Message: "Produk gagal ditemukan", Error: err}
	}

	return response.Details{Code: 200, Message: "Produk berhasil ditemukan", Error: nil}
}

func (pr *ProductRepository) GetByID(tx *gorm.DB, product *model.ResponseForGetProductByID, productID uuid.UUID, user entity.User) response.Details {
	query := fmt.Sprintf("SELECT users.id AS owner_id, users.name AS owner_name, users.url_photo_profile AS owner_photo_profile, products.id AS product_id, products.name AS product_name, products.description AS product_description, products.price AS product_price, ACOS(SIN(RADIANS(%.6f)) * SIN(RADIANS(users.latitude)) + COS(RADIANS(%.6f)) * COS(RADIANS(users.latitude)) * COS(RADIANS(users.longitude) - RADIANS(%.6f))) * 6371000 AS owner_distance FROM products INNER JOIN users ON products.user_id = users.id WHERE products.id = '%s'", user.Latitude, user.Latitude, user.Longitude, productID)
	if err := tx.Debug().Raw(query).Scan(product).Error; err != nil {
		return response.Details{Code: 500, Message: "Produk gagal ditemukan", Error: err}
	}

	return response.Details{Code: 200, Message: "Produk berhasil ditemukan", Error: nil}
}

func (pr *ProductRepository) FindActiveProducts(tx *gorm.DB, product *[]model.ResponseForActiveProducts, user entity.User) response.Details {
	query := fmt.Sprintf("SELECT users.name AS buyer_name, transactions.id AS transaction_id, products.id AS product_id, products.name AS product_name, products.price AS product_price, products.cancel_code AS cancel_code, ( SELECT url FROM media WHERE product_id = products.id LIMIT 1 ) AS url_product FROM products LEFT JOIN transactions ON products.id = transactions.product_id LEFT JOIN users ON transactions.user_id = users.id WHERE products.user_id = '%s' AND transactions.status = 'paid'", user.ID.String())
	if err := tx.Debug().Raw(query).Scan(product).Error; err != nil {
		return response.Details{Code: 500, Message: "Produk yang aktif gagal ditemukan", Error: err}
	}

	return response.Details{Code: 200, Message: "Produk yang aktif berhasil ditemukan"}
}

func (pr *ProductRepository) Update(tx *gorm.DB, product *entity.Product) response.Details {
	if err := tx.Debug().Updates(product).Error; err != nil {
		return response.Details{Code: 500, Message: "Produk gagal diperbarui", Error: err}
	}

	return response.Details{Code: 200, Message: "Produk berhasil diperbarui", Error: nil}
}

func (pr *ProductRepository) Delete(tx *gorm.DB, product *entity.Product) response.Details {
	if err := tx.Debug().Delete(product).Error; err != nil {
		return response.Details{Code: 500, Message: "Produk gagal dihapus", Error: err}
	}

	return response.Details{Code: 200, Message: "Produk berhasil dihapus", Error: nil}
}

func (pr *ProductRepository) Create(tx *gorm.DB, product *entity.Product) response.Details {
	if err := tx.Debug().Create(product).Error; err != nil {
		return response.Details{Code: 500, Message: "Product gagal dibuat", Error: err}
	}

	return response.Details{Code: 200, Message: "Product berhasil dibuat", Error: nil}
}
