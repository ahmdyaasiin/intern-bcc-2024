package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/database/mysql"
	"intern-bcc-2024/pkg/jwt"
	"intern-bcc-2024/pkg/response"
	"log"
	"os"
	"strconv"
)

type IProductService interface {
	SearchProducts(requests model.RequestForSearch) (*model.ResponseSearch, response.Details)
	DetailProduct(id uuid.UUID, ctx *gin.Context) (*model.ResponseForGetProductByID, response.Details)
	DetailProductOwner(id uuid.UUID, ctx *gin.Context) (*model.ResponseForGetProductByIDOwner, response.Details)
	Find(requests model.ParamForFind) (*entity.Product, response.Details)
	FindActiveProducts(ctx *gin.Context) (*[]model.ResponseForActiveProducts, response.Details)
	DeleteProduct(ctx *gin.Context, id uuid.UUID) response.Details
}

type ProductService struct {
	db      *gorm.DB
	cr      repository.ICategoryRepository
	pr      repository.IProductRepository
	mr      repository.IMediaRepository
	ur      repository.IUserRepository
	jwtAuth jwt.Interface
}

func NewProductService(categoryService repository.ICategoryRepository, productRepository repository.IProductRepository, mediaRepository repository.IMediaRepository, userRepository repository.IUserRepository, jwtAuth jwt.Interface) IProductService {
	return &ProductService{
		db:      mysql.Connection,
		cr:      categoryService,
		pr:      productRepository,
		mr:      mediaRepository,
		ur:      userRepository,
		jwtAuth: jwtAuth,
	}
}

func (ps *ProductService) SearchProducts(requests model.RequestForSearch) (*model.ResponseSearch, response.Details) {
	products := new([]model.ResponseForSearch)
	category := new(entity.Category)
	allCategories := new([]model.ResponseForHomePage)

	res := new(model.ResponseSearch)

	tx := ps.db.Begin()
	defer tx.Rollback()

	limit, err := strconv.Atoi(os.Getenv("LIMIT_PRODUCTS"))
	offset := (requests.Page - 1) * limit
	if err != nil {
		log.Println(err)

		return res, response.Details{Code: 500, Message: "Failed convert .env key (limit products)"}
	}

	requests.Limit = limit
	requests.Offset = offset
	respDetails := ps.pr.Search(tx, category, products, requests)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return res, respDetails
	}

	respDetails = ps.cr.GetAllCategories(tx, allCategories)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return res, respDetails
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return res, response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	res.Product = products
	res.Category = allCategories
	return res, response.Details{Code: 200, Message: "Success get searched products", Error: nil}
}

func (ps *ProductService) DetailProduct(id uuid.UUID, ctx *gin.Context) (*model.ResponseForGetProductByID, response.Details) {
	product := new(model.ResponseForGetProductByID)
	medias := new([]entity.Media)

	tx := ps.db.Begin()
	defer tx.Rollback()

	user, err := ps.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return product, response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	respDetails := ps.pr.GetByID(tx, product, id, user)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return product, respDetails
	}

	respDetails = ps.mr.GetMedia(tx, medias, id)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return product, respDetails
	}

	m := make([]string, len(*medias))
	for i, media := range *medias {
		m[i] = media.Url
	}

	product.Media = m
	return product, response.Details{Code: 200, Message: "Success get product details", Error: nil}
}

func (ps *ProductService) DetailProductOwner(id uuid.UUID, ctx *gin.Context) (*model.ResponseForGetProductByIDOwner, response.Details) {
	res := new(model.ResponseForGetProductByIDOwner)
	product := new(entity.Product)
	medias := new([]entity.Media)

	tx := ps.db.Begin()
	defer tx.Rollback()

	user, err := ps.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return res, response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	respDetails := ps.pr.FindProductOwner(tx, product, model.ParamForFind{
		ID:     id,
		UserID: user.ID,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return res, respDetails
	}

	respDetails = ps.mr.GetMedia(tx, medias, id)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return res, respDetails
	}

	m := make([]string, len(*medias))
	for i, media := range *medias {
		m[i] = media.Url
	}

	res = &model.ResponseForGetProductByIDOwner{
		ProductID:          product.ID,
		ProductName:        product.Name,
		ProductDescription: product.Description,
		ProductPrice:       product.Price,
	}
	res.Media = m
	return res, response.Details{Code: 200, Message: "Success get product details", Error: nil}
}

func (ps *ProductService) Find(requests model.ParamForFind) (*entity.Product, response.Details) {
	product := new(entity.Product)

	tx := ps.db.Begin()
	defer tx.Rollback()

	respDetails := ps.pr.Find(tx, product, requests)
	if respDetails.Error != nil {
		return product, respDetails
	}

	if err := tx.Commit().Error; err != nil {
		return product, response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return product, response.Details{Code: 200, Message: "Success get product", Error: nil}
}

func (ps *ProductService) FindActiveProducts(ctx *gin.Context) (*[]model.ResponseForActiveProducts, response.Details) {
	product := new([]model.ResponseForActiveProducts)

	tx := ps.db.Begin()
	defer tx.Rollback()

	user, err := ps.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return product, response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	respDetails := ps.pr.FindActiveProducts(tx, product, user)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return product, respDetails
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return product, response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return product, response.Details{Code: 200, Message: "Success get all active products", Error: nil}
}

func (ps *ProductService) DeleteProduct(ctx *gin.Context, id uuid.UUID) response.Details {
	product := new(entity.Product)

	tx := ps.db.Begin()
	defer tx.Rollback()

	user, err := ps.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	respDetails := ps.pr.Find(tx, product, model.ParamForFind{
		ID: id,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if user.ID != product.UserID {
		log.Println("the user is not the owner of product")

		return response.Details{Code: 403, Message: "Your'e not the owner of the product"}
	}

	if respDetails = ps.mr.DeleteAllMedia(tx, id); respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if respDetails = ps.pr.Delete(tx, product); respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success delete product", Error: nil}
}
