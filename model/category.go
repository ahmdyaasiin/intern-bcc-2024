package model

import "github.com/google/uuid"

type ResponseForHomePage struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	UrlCategory  string    `json:"url_category"`
	TotalProduct int       `json:"total_product"`
}