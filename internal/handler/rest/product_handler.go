package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
	"intern-bcc-2024/pkg/validation"
	"log"
	"strconv"
)

func (r *Rest) SearchProducts(ctx *gin.Context) {
	query := ctx.Query("query")
	category := ctx.Query("category")
	sort := ctx.Query("sort")
	pageParam := ctx.Query("page")
	page, err := strconv.Atoi(pageParam)

	if err != nil {
		response.MessageOnly(ctx, 422, "Validation error (page parameter)")
		return
	}

	if sort != "distance" && sort != "price" {
		sort = "default"
	}

	if page <= 0 {
		page = 1
	}

	user, ok := ctx.Get("user")
	if !ok {

		response.MessageOnly(ctx, 500, "Failed to get latitude and longitude")
		return
	}

	products, respDetails := r.service.ProductService.SearchProducts(model.RequestForSearch{
		Query: query, Category: category, Sort: sort, Page: page, Latitude: user.(entity.User).Latitude, Longitude: user.(entity.User).Longitude, UserID: user.(entity.User).ID,
	})
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "Success get the searched products", products)
}

func (r *Rest) DetailProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.MessageOnly(ctx, 422, "Failed convert id")
		return
	}

	product, respDetails := r.service.ProductService.DetailProduct(ctx, id)
	if respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "Success get product", product)
}

func (r *Rest) AddProduct(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(1 << 20)
	if err != nil {
		ctx.JSON(422, gin.H{"message": "Gagal memproses formulir multipart"})
		return
	}

	form := ctx.Request.MultipartForm
	if form == nil {
		ctx.JSON(422, gin.H{"message": "Formulir multipart tidak ditemukan"})
		return
	}

	var requests model.RequestForAddProduct
	if err = ctx.ShouldBind(&requests); err != nil {
		var ve validator.ValidationErrors
		errorList := validation.GetError(err, ve)
		if errorList != nil {
			log.Println("Failed to validate user requests")

			response.WithErrors(ctx, 422, "Failed to validate user requests", errorList)
			return
		}

		log.Println("Failed to bind requests")

		response.MessageOnly(ctx, 422, "Failed to bind requests")
		return
	}

	requests.Photo = form.File["photo"]
	respDetails := r.service.ProductService.AddProduct(ctx, requests)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "User has been successfully verified")
}

func (r *Rest) DetailProductOwner(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.MessageOnly(ctx, 422, "Failed convert id")
		return
	}

	product, respDetails := r.service.ProductService.DetailProductOwner(ctx, id)
	if respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "Berhasil mendapatkan data produk", product)
}

func (r *Rest) UpdateProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.MessageOnly(ctx, 422, "Failed convert id")
		return
	}

	var requests model.RequestForEditProduct
	if err := ctx.ShouldBind(&requests); err != nil {
		var ve validator.ValidationErrors
		errorList := validation.GetError(err, ve)
		if errorList != nil {
			log.Println("Failed to validate user requests")

			response.WithErrors(ctx, 422, "Failed to validate user requests", errorList)
			return
		}

		log.Println("Failed to bind requests")

		response.MessageOnly(ctx, 422, "Failed to bind requests")
		return
	}

	respDetails := r.service.ProductService.UpdateProduct(ctx, requests, id)
	if respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "Berhasil memperbarui produk")
}

func (r *Rest) DeleteProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.MessageOnly(ctx, 422, "Failed convert id")
		return
	}

	respDetails := r.service.ProductService.DeleteProduct(ctx, id)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.MessageOnly(ctx, 200, "Success delete product")
}

func (r *Rest) FindActiveProducts(ctx *gin.Context) {
	products, respDetails := r.service.ProductService.FindActiveProducts(ctx)
	if respDetails.Error != nil {
		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "Success get active products", products)

}
