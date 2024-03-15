package service

import (
	"errors"
	"github.com/google/uuid"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/bcrypt"
	"intern-bcc-2024/pkg/jwt"
	"intern-bcc-2024/pkg/mail"
	"intern-bcc-2024/pkg/response"
	"os"
	"strconv"
	"time"
)

type IUserService interface {
	Register(requests model.RequestForRegister) (model.ResponseForRegister, response.Details)
	Verify(requests model.RequestForVerify) response.Details
	ChangePassword(token string, request model.RequestForChangePassword) response.Details
	Login(requests model.RequestForLogin) (model.ResponseForLogin, response.Details)
	Find(requests model.ParamForFind) (entity.User, response.Details)
}

type UserService struct {
	ur      repository.IUserRepository
	or      repository.IOtpRepository
	tr      repository.ITokenRepository
	sr      repository.ISessionRepository
	bcrypt  bcrypt.Interface
	jwtAuth jwt.Interface
}

func NewUserService(userRepository repository.IUserRepository, otpRepository repository.IOtpRepository, tokenRepository repository.ITokenRepository, sessionRepository repository.ISessionRepository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) IUserService {
	return &UserService{
		ur:      userRepository,
		or:      otpRepository,
		tr:      tokenRepository,
		sr:      sessionRepository,
		bcrypt:  bcrypt,
		jwtAuth: jwtAuth,
	}
}

func (us *UserService) Register(requests model.RequestForRegister) (model.ResponseForRegister, response.Details) {
	_, respDetails := us.ur.Find(model.ParamForFind{
		Email: requests.Email,
	})
	if respDetails.Error == nil {
		return model.ResponseForRegister{}, response.Details{Code: 409, Message: "Email has been registered", Error: errors.New("email registered")}
	}

	hashPassword, err := us.bcrypt.GenerateFromPassword(requests.Password)
	if err != nil {
		return model.ResponseForRegister{}, response.Details{Code: 500, Message: "Failed to generate password", Error: err}
	}

	user := &entity.User{
		ID:              uuid.New(),
		Name:            requests.Name,
		Email:           requests.Email,
		Password:        hashPassword,
		Address:         requests.Address,
		Latitude:        requests.Latitude,
		Longitude:       requests.Longitude,
		StatusAccount:   "inactive",
		UrlPhotoProfile: "default.jpg",
	}

	respDetails = us.ur.Create(user)
	if respDetails.Error != nil {
		return model.ResponseForRegister{}, respDetails
	}

	code := mail.GenerateVerificationCode()
	expiredTime, err := strconv.Atoi(os.Getenv("EXPIRED_OTP"))
	if err != nil {
		return model.ResponseForRegister{}, response.Details{Code: 500, Message: "Failed to convert expired time from .env", Error: err}
	}

	otp := &entity.OtpCode{
		ID:               uuid.New(),
		UserID:           user.ID,
		VerificationCode: code,
		LastSentAt:       time.Now().Local().Add(-5 * time.Minute).UnixMilli(),
		ExpiredAt:        time.Now().Local().Add(time.Duration(expiredTime) * time.Minute).UnixMilli(),
	}

	respDetails = us.or.Create(otp)
	if respDetails.Error != nil {
		return model.ResponseForRegister{}, respDetails
	}

	if err = mail.SendEmail(user.Email, "Verification Code", "Your Verification Code: "+code); err != nil {
		return model.ResponseForRegister{}, response.Details{Code: 500, Message: "Failed to send verification code to user, please login and resend your code", Error: err}
	}

	return model.ResponseForRegister{
		ID: user.ID,
	}, response.Details{Code: 201, Message: "Success register", Error: nil}
}

func (us *UserService) Verify(requests model.RequestForVerify) response.Details {
	otp, respDetails := us.or.Find(model.ParamForFind{
		UserID: requests.UserID,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	if otp.VerificationCode != requests.VerificationCode {
		return response.Details{Code: 401, Message: "Verification code is wrong", Error: errors.New("verification code is wrong")}
	}

	if otp.ExpiredAt < time.Now().Local().UnixMilli() {
		return response.Details{Code: 401, Message: "Verification code is expired", Error: errors.New("verification code is expired")}
	}

	user, respDetails := us.ur.Find(model.ParamForFind{
		ID: requests.UserID,
	})

	if user.StatusAccount == "active" {
		return response.Details{Code: 403, Message: "User already verified", Error: errors.New("user already verified")}
	}

	user.StatusAccount = "active"
	respDetails = us.ur.Verify(user, otp)
	if respDetails.Error != nil {
		return respDetails
	}

	return respDetails
}

func (us *UserService) ChangePassword(token string, request model.RequestForChangePassword) response.Details {
	tokenDetails, respDetails := us.tr.Find(model.ParamForFind{
		Token: token, // token must be unique
	})
	if respDetails.Error != nil {
		return respDetails
	}

	if tokenDetails.ExpiredAt < time.Now().Local().UnixMilli() {
		return response.Details{Code: 401, Message: "Token is expired", Error: errors.New("token is expired")}
	}

	user, respDetails := us.ur.Find(model.ParamForFind{
		ID: tokenDetails.UserID,
	})
	hashPassword, err := us.bcrypt.GenerateFromPassword(request.Password)
	if err != nil {
		return response.Details{Code: 500, Message: "Failed to generate password", Error: err}
	}

	user.Password = hashPassword
	if respDetails = us.ur.Change(user, tokenDetails); respDetails.Error != nil {
		return respDetails
	}

	return response.Details{Code: 200, Message: "Success change password", Error: nil}
}

func (us *UserService) Login(requests model.RequestForLogin) (model.ResponseForLogin, response.Details) {
	res := model.ResponseForLogin{}

	user, respDetails := us.ur.Find(model.ParamForFind{
		Email: requests.Email,
	})
	if respDetails.Error != nil {
		return res, respDetails
	}

	if user.StatusAccount == "blocked" {
		return res, response.Details{Code: 403, Message: "Your account is blocked", Error: errors.New("blocked")}
	}

	if user.StatusAccount == "inactive" {
		res.UserID = user.ID
		return res, response.Details{Code: 403, Message: "Please verify your account first", Error: errors.New("unverified")}
	}

	err := us.bcrypt.CompareAndHashPassword(user.Password, requests.Password)
	if err != nil {
		return res, response.Details{Code: 500, Message: "Wrong password", Error: err}
	}

	accessToken, err := us.jwtAuth.CreateAccessToken(user.ID)
	if err != nil {
		return res, response.Details{Code: 500, Message: "Failed to create access token", Error: err}
	}

	rToken, respDetails := us.sr.Find(model.ParamForFind{
		UserID: user.ID,
	})
	if respDetails.Error == nil {
		if respDetails = us.sr.Delete(&rToken); respDetails.Error != nil {
			return res, respDetails
		}
	}

	refreshToken, err := us.jwtAuth.CreateRefreshToken(user.ID)
	if err != nil {
		return res, response.Details{Code: 500, Message: "Failed to create refresh token", Error: err}
	}

	expiredAt, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TIME"))
	rToken = entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiredAt: time.Now().Local().Add(time.Duration(expiredAt) * time.Hour).UnixMilli(),
	}
	if respDetails = us.sr.Create(&rToken); respDetails.Error != nil {
		return res, respDetails
	}

	res.AccessToken = accessToken
	res.RefreshToken = refreshToken

	return res, response.Details{Code: 200, Message: "Success login", Error: nil}
}

func (us *UserService) Find(requests model.ParamForFind) (entity.User, response.Details) {
	return us.ur.Find(requests)
}
