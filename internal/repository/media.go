package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/pkg/response"
)

type IMediaRepository interface {
	GetMedia(id uuid.UUID) ([]*entity.Media, response.Details)
}

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) IMediaRepository {
	return &MediaRepository{db}
}

func (mr *MediaRepository) GetMedia(id uuid.UUID) ([]*entity.Media, response.Details) {
	var medias []*entity.Media
	if err := mr.db.Debug().Where("product_id = ?", id).Find(&medias).Error; err != nil {
		return medias, response.Details{Code: 500, Message: "Failed to get all media of product", Error: err}
	}

	return medias, response.Details{Code: 200, Message: "Success to get all media of product", Error: nil}
}
