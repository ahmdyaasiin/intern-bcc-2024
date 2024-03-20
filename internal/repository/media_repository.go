package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/pkg/response"
)

type IMediaRepository interface {
	GetMedia(tx *gorm.DB, medias *[]entity.Media, id uuid.UUID) response.Details
	DeleteAllMedia(tx *gorm.DB, productID uuid.UUID) response.Details
}

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) IMediaRepository {
	return &MediaRepository{db}
}

func (mr *MediaRepository) GetMedia(tx *gorm.DB, medias *[]entity.Media, id uuid.UUID) response.Details {
	if err := tx.Debug().Where("product_id = ?", id).Find(&medias).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to get all media of product", Error: err}
	}

	return response.Details{Code: 200, Message: "Success to get all media of product", Error: nil}
}

func (mr *MediaRepository) DeleteAllMedia(tx *gorm.DB, productID uuid.UUID) response.Details {
	if err := tx.Debug().Where("product_id = ?", productID).Delete(entity.Media{}).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed delete all media product", Error: err}
	}

	return response.Details{Code: 200, Message: "Success delete all media product", Error: nil}
}
