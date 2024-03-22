package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/database/mysql"
	"intern-bcc-2024/pkg/jwt"
	"intern-bcc-2024/pkg/mail"
	"intern-bcc-2024/pkg/response"
	"intern-bcc-2024/pkg/supabase"
	"log"
	"os"
	"strconv"
)

type IProductService interface {
	SearchProducts(requests model.RequestForSearch) (*model.ResponseSearch, response.Details)
	DetailProduct(ctx *gin.Context, id uuid.UUID) (*model.ResponseForGetProductByID, response.Details)
	DetailProductOwner(ctx *gin.Context, id uuid.UUID) (*model.ResponseForGetProductByIDOwner, response.Details)
	Find(requests model.ParamForFind) (*entity.Product, response.Details)
	FindActiveProducts(ctx *gin.Context) (*[]model.ResponseForActiveProducts, response.Details)
	DeleteProduct(ctx *gin.Context, id uuid.UUID) response.Details
	AddProduct(ctx *gin.Context, requests model.RequestForAddProduct) response.Details
	UpdateProduct(ctx *gin.Context, requests model.RequestForEditProduct, id uuid.UUID) response.Details
}

type ProductService struct {
	db       *gorm.DB
	cr       repository.ICategoryRepository
	pr       repository.IProductRepository
	mr       repository.IMediaRepository
	ur       repository.IUserRepository
	jwtAuth  jwt.Interface
	supabase supabase.Interface
}

func NewProductService(categoryService repository.ICategoryRepository, productRepository repository.IProductRepository, mediaRepository repository.IMediaRepository, userRepository repository.IUserRepository, jwtAuth jwt.Interface, supabase supabase.Interface) IProductService {
	return &ProductService{
		db:       mysql.Connection,
		cr:       categoryService,
		pr:       productRepository,
		mr:       mediaRepository,
		ur:       userRepository,
		jwtAuth:  jwtAuth,
		supabase: supabase,
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

func (ps *ProductService) DetailProduct(ctx *gin.Context, id uuid.UUID) (*model.ResponseForGetProductByID, response.Details) {
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

func (ps *ProductService) DetailProductOwner(ctx *gin.Context, id uuid.UUID) (*model.ResponseForGetProductByIDOwner, response.Details) {
	res := new(model.ResponseForGetProductByIDOwner)
	product := new(entity.Product)
	medias := new([]entity.Media)
	category := new(entity.Category)
	categories := new([]model.ResponseForHomePage)

	tx := ps.db.Begin()
	defer tx.Rollback()

	user, err := ps.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return res, response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	respDetails := ps.pr.Find(tx, product, model.ParamForFind{
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

	respDetails = ps.cr.GetAllCategories(tx, categories)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return res, respDetails
	}

	respDetails = ps.cr.Find(tx, category, model.ParamForFind{
		ID: product.CategoryID,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return res, respDetails
	}

	res = &model.ResponseForGetProductByIDOwner{
		Product: model.ResponseForProductForIDOwner{
			ID:           product.ID,
			Name:         product.Name,
			Description:  product.Description,
			Price:        product.Price,
			CategoryID:   product.CategoryID,
			CategoryName: category.Name,
		},
		Categories: categories,
	}
	res.Product.Media = m
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

		return response.Details{Code: 403, Message: "Anda bukan pemilik produk ini"}
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

	return response.Details{Code: 200, Message: "Berhasil menghapus produk", Error: nil}
}

func (ps *ProductService) AddProduct(ctx *gin.Context, requests model.RequestForAddProduct) response.Details {
	product := new(entity.Product)
	media := new(entity.Media)
	category := new(entity.Category)

	tx := ps.db.Begin()
	defer tx.Rollback()

	user, err := ps.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	respDetails := ps.cr.Find(tx, category, model.ParamForFind{
		Name: requests.Category,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	price, err := strconv.Atoi(requests.Price)
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to convert string into int", Error: err}
	}

	product = &entity.Product{
		ID:          uuid.New(),
		UserID:      user.ID,
		CategoryID:  category.ID,
		Name:        requests.Name,
		Description: requests.Description,
		Price:       uint64(price),
		CancelCode:  mail.GenerateSixCode(),
	}
	respDetails = ps.pr.Create(tx, product)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	for i, photo := range requests.Photo {

		photo.Filename = fmt.Sprintf("%s-%s", mail.GenerateRandomString(30), photo.Filename)
		link, err := ps.supabase.Upload(photo)
		if err != nil {
			log.Println(err)

			return response.Details{Code: 500, Message: "Failed upload photo - " + strconv.Itoa(i+1), Error: err}
		}

		media = &entity.Media{
			ID:        uuid.New(),
			ProductID: product.ID,
			Url:       link,
		}
		respDetails = ps.mr.Create(tx, media)
		if respDetails.Error != nil {
			log.Println(respDetails.Error)

			return respDetails
		}
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 201, Message: "Success create product", Error: nil}
}

func (ps *ProductService) UpdateProduct(ctx *gin.Context, requests model.RequestForEditProduct, id uuid.UUID) response.Details {
	product := new(entity.Product)
	media := new([]entity.Media)
	category := new(entity.Category)

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
		log.Println("bukan pemilik produk")

		return response.Details{Code: 403, Message: "Anda bukan pemiliki produk", Error: errors.New("bukan pemilik produk")}
	}

	oldPhoto := "("
	for i, op := range requests.OldPhoto {
		oldPhoto += fmt.Sprintf("'%s'", op)

		if i != len(requests.OldPhoto)-1 {
			oldPhoto += ","
		}
	}

	oldPhoto += ")"
	if respDetails = ps.mr.FindWithout(tx, media, product.ID, oldPhoto); respDetails.Error != nil {
		log.Println("failed get media without old_photo")

		return respDetails
	}

	for _, m := range *media {
		if respDetails = ps.mr.Delete(tx, &m); respDetails.Error != nil {
			log.Println(respDetails.Error)

			return respDetails
		}
	}

	price, err := strconv.Atoi(requests.Price)
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed convert to int", Error: err}
	}

	for i, photo := range requests.Photo {

		photo.Filename = fmt.Sprintf("%s-%s", mail.GenerateRandomString(30), photo.Filename)
		link, err := ps.supabase.Upload(photo)
		if err != nil {
			log.Println(err)

			return response.Details{Code: 500, Message: "Failed upload photo - " + strconv.Itoa(i+1), Error: err}
		}

		med := &entity.Media{
			ID:        uuid.New(),
			ProductID: product.ID,
			Url:       link,
		}
		respDetails = ps.mr.Create(tx, med)
		if respDetails.Error != nil {
			log.Println(respDetails.Error)

			return respDetails
		}
	}

	product.Name = requests.Name
	product.Description = requests.Description
	product.Price = uint64(price)
	product.CategoryID = category.ID

	if respDetails = ps.pr.Update(tx, product); respDetails.Error != nil {
		return respDetails
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success update product", Error: nil}

}
