package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/jwt"
	"intern-bcc-2024/pkg/response"
	"os"
	"strconv"
)

type IProductService interface {
	SearchProducts(requests model.RequestForSearch) (model.ResponseSearch, response.Details)
	GetProduct(id uuid.UUID, ctx *gin.Context) (*model.ResponseForGetProductByID, response.Details)
}

type ProductService struct {
	pr      repository.IProductRepository
	mr      repository.IMediaRepository
	ur      repository.IUserRepository
	jwtAuth jwt.Interface
}

func NewProductService(productRepository repository.IProductRepository, mediaRepository repository.IMediaRepository, userRepository repository.IUserRepository, jwtAuth jwt.Interface) IProductService {
	return &ProductService{
		pr:      productRepository,
		mr:      mediaRepository,
		ur:      userRepository,
		jwtAuth: jwtAuth,
	}
}

func (ps *ProductService) SearchProducts(requests model.RequestForSearch) (model.ResponseSearch, response.Details) {
	limit, err := strconv.Atoi(os.Getenv("LIMIT_PRODUCTS"))
	offset := (requests.Page - 1) * limit
	if err != nil {
		return model.ResponseSearch{}, response.Details{Code: 500, Message: "Failed convert .env key (limit products)"}
	}

	products, respDetails := ps.pr.Search(requests, limit, offset)
	if respDetails.Error != nil {
		return model.ResponseSearch{}, respDetails
	}

	return products, response.Details{Code: 200, Message: "Success get searched products", Error: nil}
}

func (ps *ProductService) GetProduct(id uuid.UUID, ctx *gin.Context) (*model.ResponseForGetProductByID, response.Details) {
	user, err := ps.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		return &model.ResponseForGetProductByID{}, response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	product, respDetails := ps.pr.GetByID(id, user)
	if respDetails.Error != nil {
		return &model.ResponseForGetProductByID{}, respDetails
	}

	medias, respDetails := ps.mr.GetMedia(id)
	if respDetails.Error != nil {
		return &model.ResponseForGetProductByID{}, respDetails
	}

	m := make([]string, len(medias))
	for i, media := range medias {
		m[i] = media.UrlMedia
	}

	product.Media = m

	return product, response.Details{Code: 200, Message: "Success get product details", Error: nil}
}
