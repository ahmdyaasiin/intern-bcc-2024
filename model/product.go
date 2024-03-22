package model

import (
	"github.com/google/uuid"
	"mime/multipart"
)

/*
	Request Struct
*/

type RequestForSearch struct {
	Query     string
	Category  string
	Sort      string
	Page      int
	Latitude  float64
	Longitude float64
	UserID    uuid.UUID
	Limit     int
	Offset    int
}

type RequestForAddProduct struct {
	Name        string                  `form:"name" binding:"required"`
	Description string                  `form:"description" binding:"required"`
	Price       string                  `form:"price" binding:"required"`
	Category    string                  `form:"category" binding:"required"`
	Photo       []*multipart.FileHeader `form:"photo" binding:"required,max=3"`
}

type RequestForEditProduct struct {
	Name        string                  `form:"name" binding:"required"`
	Description string                  `form:"description" binding:"required"`
	Price       string                  `form:"price" binding:"required"`
	Category    string                  `form:"category" binding:"required"`
	OldPhoto    []string                `form:"old_photo" binding:"max=3"`
	Photo       []*multipart.FileHeader `form:"photo" binding:"max=3"`
}

/*
	Response Struct
*/

type ResponseSearch struct {
	Product  *[]ResponseForSearch   `json:"products"`
	Category *[]ResponseForHomePage `json:"categories"`
}

type ResponseForActiveProducts struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	ProductID     string    `json:"product_id"`
	ProductName   string    `json:"product_name"`
	ProductPrice  string    `json:"product_price"`
	CancelCode    string    `json:"cancel_code"`
	UrlProduct    string    `json:"url_product"`
	BuyerName     string    `json:"buyer_name"`
}

type ResponseForGetProductByIDOwner struct {
	Product    ResponseForProductForIDOwner `json:"product"`
	Categories *[]ResponseForHomePage       `json:"categories"`
}

type ResponseForSearch struct {
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductPrice    int64     `json:"product_price"`
	UrlPhotoProduct string    `json:"url_photo_product"`
	OwnerID         uuid.UUID `json:"owner_id"`
	OwnerName       string    `json:"owner_name"`
	OwnerDistance   string    `json:"owner_distance"`
}

type ResponseForGetProductByID struct {
	ProductID          uuid.UUID `json:"product_id"`
	ProductName        string    `json:"product_name"`
	ProductDescription string    `json:"product_description"`
	ProductPrice       uint64    `json:"product_price"`
	Media              []string  `json:"media"`
	OwnerID            uuid.UUID `json:"owner_id"`
	OwnerName          string    `json:"owner_name"`
	OwnerDistance      string    `json:"owner_distance"`
	OwnerPhotoProfile  string    `json:"owner_photo_profile"`
}

type ResponseForProductForIDOwner struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        uint64    `json:"price"`
	Media        []string  `json:"media"`
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
}
