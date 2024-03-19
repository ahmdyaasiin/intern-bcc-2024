package model

import "github.com/google/uuid"

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

type ResponseForSearch struct {
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductPrice    int64     `json:"product_price"`
	UrlPhotoProduct string    `json:"url_photo_product"`
	OwnerID         uuid.UUID `json:"owner_id"`
	OwnerName       string    `json:"owner_name"`
	OwnerDistance   string    `json:"owner_distance"`
}

type ResponseSearch struct {
	Product  *[]ResponseForSearch   `json:"products"`
	Category *[]ResponseForHomePage `json:"categories"`
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

type ResponseForGetProductByIDOwner struct {
	ProductID          uuid.UUID `json:"product_id"`
	ProductName        string    `json:"product_name"`
	ProductDescription string    `json:"product_description"`
	ProductPrice       uint64    `json:"product_price"`
	Media              []string  `json:"media"`
}

type ResponseForActiveProducts struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	ProductID     string    `json:"product_id"`
	ProductName   string    `json:"product_name"`
	ProductPrice  string    `json:"product_price"`
	CancelCode    string    `json:"cancel_code"`
	UrlProduct    string    `json:"url_product"`
	OwnerName     string    `json:"owner_name"`
}
