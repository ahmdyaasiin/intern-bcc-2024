package model

import "github.com/google/uuid"

type ParamForFind struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ProductID uuid.UUID
	Email     string
	Token     string
	Name      string
}
