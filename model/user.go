package model

import "github.com/google/uuid"

type ParamForFind struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Email  string
	Token  string
}

/*
	@ Request Struct
*/

type RequestForRegister struct {
	Name            string  `json:"name" binding:"required,min=3,max=18"`
	Email           string  `json:"email" binding:"required,email,secureDomain"`
	Password        string  `json:"password" binding:"required,eqfield=ConfirmPassword,min=8,max=64"`
	ConfirmPassword string  `json:"confirm_password" binding:"required,min=8,max=64"`
	Address         string  `json:"address" binding:"required"`
	Latitude        float64 `json:"latitude" binding:"required"`
	Longitude       float64 `json:"longitude" binding:"required"`
}

type RequestForVerify struct {
	UserID           uuid.UUID `json:"id" binding:"required"`
	VerificationCode string    `json:"verification_code" binding:"required"`
}

type RequestForResend struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type RequestForLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RequestForRenewAccessToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RequestForReset struct {
	Email string `json:"email" binding:"required,email"`
}

type RequestForChangePassword struct {
	Password        string `json:"password" binding:"required,eqfield=ConfirmPassword,min=8,max=64"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=64"`
}

/*
	@ Response Struct
*/

type ResponseForRegister struct {
	ID uuid.UUID `json:"id"`
}

type ResponseForLogin struct {
	UserID       uuid.UUID `json:"-"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

type ResponseForRenew struct {
	AccessToken string `json:"access_token"`
}
