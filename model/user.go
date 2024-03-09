package model

import "github.com/google/uuid"

type RequestForRegister struct {
	Name            string  `json:"name" binding:"required,min=3,max=18"`
	Email           string  `json:"email" binding:"required,email,secureDomain"`
	Password        string  `json:"password" binding:"required,eqfield=ConfirmPassword,min=8,max=64"`
	ConfirmPassword string  `json:"confirm_password" binding:"required,min=8,max=64"`
	Address         string  `json:"address" binding:"required"`
	Latitude        float64 `json:"latitude" binding:"required"`
	Longitude       float64 `json:"longitude" binding:"required"`
}

type ResponseRegister struct {
	ID uuid.UUID `json:"id"`
}

type OtpParam struct {
	//ID               uuid.UUID `json:"id" binding:"required"`
	UserID           uuid.UUID `json:"id" binding:"required"`
	VerificationCode string    `json:"verification_code"`
}

type UserParam struct {
	ID    uuid.UUID `json:"-" binding:"required_without=Email"`
	Email string    `json:"-"`
}

type RequestForResend struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type SessionParam struct {
	UserID uuid.UUID `json:"-" binding:"required_without=Token"`
	Token  string    `json:"-"`
}

type RequestForLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ResponseForLogin struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RequestForRenew struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ResponseForRenew struct {
	AccessToken string `json:"access_token"`
}

type RequestForReset struct {
	Email string `json:"email" binding:"required,email"`
}

type TokenParam struct {
	UserID uuid.UUID `json:"-"`
	Token  string    `json:"-"`
}

type RequestForChangePassword struct {
	Password        string `json:"password" binding:"required,eqfield=ConfirmPassword,min=8,max=64"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=64"`
}

type VerifyAccount struct {
	ID               uuid.UUID `json:"id" binding:"required"`
	VerificationCode string    `json:"verification_code" binding:"required"`
}

type ResendAndRenew struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetPassword struct {
	Email string `json:"email" binding:"required"`
}

type ChangePassword struct {
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
