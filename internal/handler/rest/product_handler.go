package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/response"
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

	product, respDetails := r.service.ProductService.DetailProduct(id, ctx)
	if respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "Success get product", product)
}

func (r *Rest) AddProduct(ctx *gin.Context) {

}

func (r *Rest) DetailProductOwner(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.MessageOnly(ctx, 422, "Failed convert id")
		return
	}

	product, respDetails := r.service.ProductService.DetailProductOwner(id, ctx)
	if respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "Success get product", product)
}

func (r *Rest) UpdateProduct(ctx *gin.Context) {

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
