package model

import "github.com/google/uuid"

/*
	Request Struct
*/

// -

/*
	Response Struct
*/

type ResponseForHomePage struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Url          string    `json:"url_category"`
	TotalProduct int       `json:"total_product"`
}
