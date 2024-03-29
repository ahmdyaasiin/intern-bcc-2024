package repository

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/pkg/response"
)

type IMediaRepository interface {
	GetMedia(tx *gorm.DB, medias *[]entity.Media, id uuid.UUID) response.Details
	DeleteAllMedia(tx *gorm.DB, productID uuid.UUID) response.Details
	Create(tx *gorm.DB, media *entity.Media) response.Details
	FindWithout(tx *gorm.DB, media *[]entity.Media, productID uuid.UUID, without string) response.Details
	Delete(tx *gorm.DB, media *entity.Media) response.Details
}

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) IMediaRepository {
	return &MediaRepository{db}
}

func (mr *MediaRepository) GetMedia(tx *gorm.DB, medias *[]entity.Media, id uuid.UUID) response.Details {
	if err := tx.Debug().Where("product_id = ?", id).Find(&medias).Error; err != nil {
		return response.Details{Code: 500, Message: "Media produk gagal ditemukan", Error: err}
	}

	return response.Details{Code: 200, Message: "Media produk berhasil ditemukan", Error: nil}
}

func (mr *MediaRepository) DeleteAllMedia(tx *gorm.DB, productID uuid.UUID) response.Details {
	if err := tx.Debug().Where("product_id = ?", productID).Delete(entity.Media{}).Error; err != nil {
		return response.Details{Code: 500, Message: "Media produk gagal dihapus", Error: err}
	}

	return response.Details{Code: 200, Message: "Media produk berhasil dihapus", Error: nil}
}

func (mr *MediaRepository) Create(tx *gorm.DB, media *entity.Media) response.Details {
	if err := tx.Debug().Create(media).Error; err != nil {
		return response.Details{Code: 500, Message: "Media gagal dibuat", Error: err}
	}

	return response.Details{Code: 200, Message: "Media berhasil dibuat", Error: nil}
}

func (mr *MediaRepository) FindWithout(tx *gorm.DB, media *[]entity.Media, productID uuid.UUID, without string) response.Details {
	query := fmt.Sprintf("SELECT * FROM media WHERE product_id = '%s' AND url NOT IN %s", productID, without)
	if err := tx.Debug().Raw(query).Scan(media).Error; err != nil {
		return response.Details{Code: 500, Message: "Media not found", Error: err}
	}

	return response.Details{Code: 200, Message: "Media found", Error: nil}
}

func (mr *MediaRepository) Delete(tx *gorm.DB, media *entity.Media) response.Details {
	if err := tx.Debug().Delete(media).Error; err != nil {
		return response.Details{Code: 500, Message: "Media gagal dihapus", Error: err}
	}

	return response.Details{Code: 200, Message: "Media berhasil dihapus", Error: nil}
}
