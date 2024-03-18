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

	product, respDetails := r.service.ProductService.GetProduct(id, ctx)
	if respDetails.Error != nil {

		response.MessageOnly(ctx, respDetails.Code, respDetails.Message)
		return
	}

	response.WithData(ctx, 200, "Success get product", product)
}

func (r *Rest) DetailProductOwner(ctx *gin.Context) {

}

func (r *Rest) UpdateProduct(ctx *gin.Context) {

}

func (r *Rest) DeleteProduct(ctx *gin.Context) {

}

func (r *Rest) AllMyProduct(ctx *gin.Context) {
	
}
