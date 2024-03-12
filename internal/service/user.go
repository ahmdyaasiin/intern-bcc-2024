package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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
	Register(requests model.RequestForRegister) (*entity.User, response.Details)
	Verify(requests model.RequestForVerify) response.Details
	Resend(requests model.RequestForResend) response.Details
	ResetPassword(requests model.RequestForReset) response.Details
	CheckToken(token string) response.Details
	ChangePassword(token string, request model.RequestForChangePassword) response.Details
	Login(requests model.RequestForLogin) (model.ResponseForLogin, response.Details)
	Renew(requests model.RequestForRenewAccessToken) (model.ResponseForRenew, response.Details)
	Find(requets model.ParamForFind) (entity.User, response.Details)
	Logout(ctx *gin.Context) response.Details
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

func (us *UserService) Register(requests model.RequestForRegister) (*entity.User, response.Details) {
	_, respDetails := us.ur.Find(model.ParamForFind{
		Email: requests.Email,
	})
	if respDetails.Error == nil {
		return nil, response.Details{Code: 209, Message: "Email has been registered", Error: errors.New("email registered")}
	}

	hashPassword, err := us.bcrypt.GenerateFromPassword(requests.Password)
	if err != nil {
		return nil, response.Details{Code: 500, Message: "Failed to generate password", Error: err}
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
		return nil, respDetails
	}

	code := mail.GenerateVerificationCode()
	expiredTime, err := strconv.Atoi(os.Getenv("EXPIRED_OTP"))
	if err != nil {
		return nil, response.Details{Code: 500, Message: "Failed to convert expired time from .env", Error: err}
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
		return nil, respDetails
	}

	// fake scenario
	//return nil, response.Details{Code: 500, Message: "Failed to send verification code to user", Error: errors.New("fake scenario")}

	if err = mail.SendEmail(user.Email, "Verification Code", "Your Verification Code: "+code); err != nil {
		return nil, response.Details{Code: 500, Message: "Failed to send verification code to user, please login and resend your code", Error: err}
	}

	return user, response.Details{Code: 201, Message: "Success register", Error: nil}
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

func (us *UserService) Resend(requests model.RequestForResend) response.Details {
	user, respDetails := us.ur.Find(model.ParamForFind{
		ID: requests.ID,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	if user.StatusAccount == "active" {
		return response.Details{Code: 403, Message: "User already verified", Error: errors.New("user already verified")}
	}

	otp, respDetails := us.or.Find(model.ParamForFind{
		UserID: requests.ID,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	expiredTime, err := strconv.Atoi(os.Getenv("EXPIRED_OTP"))
	if err != nil {
		return response.Details{Code: 500, Message: "Failed to convert expired time from .env", Error: err}
	}

	if otp.LastSentAt > time.Now().Local().Add(-5*time.Minute).UnixMilli() {
		return response.Details{Code: 403, Message: "Please wait 5 minutes again", Error: errors.New("limit to send otp")}
	}

	otp.VerificationCode = mail.GenerateVerificationCode()
	otp.ExpiredAt = time.Now().Local().Add(time.Duration(expiredTime) * time.Hour).UnixMilli()

	if err = mail.SendEmail(user.Email, "Verification Code", "Your Verification Code: "+otp.VerificationCode); err != nil {
		fmt.Println(err)
		return response.Details{Code: 500, Message: "Failed to resend verification code to user", Error: err}
	}

	respDetails = us.or.Update(&otp)
	if respDetails.Error != nil {
		return respDetails
	}

	return response.Details{Code: 200, Message: "Success to resend verification code to user", Error: nil}
}

func (us *UserService) ResetPassword(requests model.RequestForReset) response.Details {
	user, respDetails := us.ur.Find(model.ParamForFind{
		Email: requests.Email,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	token, respDetails := us.tr.Find(model.ParamForFind{
		UserID: user.ID,
	})
	if respDetails.Error == nil {
		if token.LastSentAt > time.Now().Local().Add(-5*time.Minute).UnixMilli() {
			return response.Details{Code: 403, Message: "Please wait 5 minutes again", Error: errors.New("limit to send link reset")} // add exact time
		}

		respDetails = us.tr.Delete(&token)
		if respDetails.Error != nil {
			return respDetails
		}
	}

	link := mail.GenerateRandomString(30)
	if err := mail.SendEmail(user.Email, "Link Reset Password", "Your link: "+os.Getenv("LINK_FRONTEND")+link); err != nil {
		return response.Details{Code: 500, Message: "Failed to send link reset password", Error: err}
	}

	token = entity.ResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     link,
		ExpiredAt: time.Now().Local().Add(time.Hour * 1).UnixMilli(),
	}

	respDetails = us.tr.Create(&token)
	if respDetails.Error != nil {
		return respDetails
	}

	return response.Details{Code: 200, Message: "Success send link reset password", Error: nil}
}

func (us *UserService) CheckToken(token string) response.Details {
	if _, respDetails := us.tr.Find(model.ParamForFind{
		Token: token,
	}); respDetails.Error != nil {
		return respDetails
	}

	return response.Details{Code: 200, Message: "Token is valid", Error: nil}
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

	if user.StatusAccount == "inactive" {
		res.UserID = user.ID
		return res, response.Details{Code: 403, Message: "Please verify your account first", Error: errors.New("unverified")}
	}

	err := us.bcrypt.CompareAndHashPassword(user.Password, requests.Password)
	if err != nil {
		return res, response.Details{Code: 500, Message: "Failed to generate password", Error: err}
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

func (us *UserService) Renew(requests model.RequestForRenewAccessToken) (model.ResponseForRenew, response.Details) {
	res := model.ResponseForRenew{}

	refreshToken, respDetails := us.sr.Find(model.ParamForFind{
		Token: requests.RefreshToken,
	})
	if respDetails.Error != nil {
		return res, respDetails
	}

	if refreshToken.ExpiredAt < time.Now().Local().UnixMilli() {
		return res, response.Details{Code: 403, Message: "Refresh token is expired", Error: errors.New("refresh token expired")}
	}

	accessToken, err := us.jwtAuth.CreateAccessToken(refreshToken.UserID)
	if err != nil {
		return res, response.Details{Code: 500, Message: "Failed to create access token", Error: err}
	}

	res.AccessToken = accessToken
	return res, response.Details{Code: 200, Message: "Success renew access token", Error: nil}
}

func (us *UserService) Find(requests model.ParamForFind) (entity.User, response.Details) {
	return us.ur.Find(requests)
}

func (us *UserService) Logout(ctx *gin.Context) response.Details {
	user, err := us.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		return response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	sToken, respDetails := us.sr.Find(model.ParamForFind{
		UserID: user.ID,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	if err = us.sr.Delete(&sToken).Error; err != nil {
		return response.Details{Code: 500, Message: "Failed to delete the session", Error: err}
	}

	return response.Details{Code: 200, Message: " Success logout", Error: nil}
}

//

//func (us *UserService) Login(requests model.RequestForLogin) (model.ResponseForLogin, error) {
//	res := model.ResponseForLogin{}
//
//	user, err := us.ur.GetUser(model.UserParam{
//		Email: requests.Email,
//	})
//	if err != nil {
//		return res, err
//	}
//
//	if user.StatusAccount == "inactive" {
//		return res, errors.New("mohon verifikasi akun anda terlebih dahulu")
//	}
//
//	err = us.bcrypt.CompareAndHashPassword(user.Password, requests.Password)
//	if err != nil {
//		return res, err
//	}
//
//	accessToken, err := us.jwtAuth.CreateAccessToken(user.ID)
//	if err != nil {
//		return res, err
//	}
//
//	checkRefreshToken, err := us.ur.GetSession(model.SessionParam{
//		UserID: user.ID,
//	})
//	if err == nil {
//		err = us.ur.ClearSession(&checkRefreshToken)
//		if err != nil {
//			return res, err
//		}
//	}
//	refreshToken, err := us.jwtAuth.CreateRefreshToken(user.ID)
//	if err != nil {
//		return res, err
//	}
//
//	expiredAt, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TIME"))
//	rToken := &entity.RefreshToken{
//		ID:        uuid.New(),
//		UserID:    user.ID,
//		Token:     refreshToken,
//		ExpiredAt: time.Now().Local().Add(time.Duration(expiredAt) * time.Hour).UnixMilli(),
//	}
//	if err = us.ur.SetSession(rToken); err != nil {
//		return res, err
//	}
//
//	res.AccessToken = accessToken
//	res.RefreshToken = refreshToken
//
//	return res, nil
//}
//
//func (us *UserService) Renew(requests model.RequestForRenew) (model.ResponseForRenew, error) {
//	res := model.ResponseForRenew{}
//
//	rToken, err := us.ur.GetSession(model.SessionParam{
//		Token: requests.RefreshToken,
//	})
//	if err != nil {
//		return res, err
//	}
//
//	if rToken.ExpiredAt < time.Now().Local().UnixMilli() {
//		return res, errors.New("refresh token kadaluwarsa")
//	}
//
//	accessToken, err := us.jwtAuth.CreateAccessToken(rToken.UserID)
//	if err != nil {
//		return res, err
//	}
//
//	res.AccessToken = accessToken
//
//	return res, nil
//}
//
//func (us *UserService) ResetGet(token string) error {
//	_, err := us.ur.GetToken(model.TokenParam{
//		Token: token,
//	})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (us *UserService) ResetPost(requests model.RequestForChangePassword, token string) error {
//	rToken, err := us.ur.GetToken(model.TokenParam{
//		Token: token,
//	})
//	if err != nil {
//		return err
//	}
//
//	user, err := us.ur.GetUser(model.UserParam{
//		ID: rToken.UserID,
//	})
//	if err != nil {
//		return err
//	}
//
//	user.Password, err = us.bcrypt.GenerateFromPassword(requests.Password)
//	if err != nil {
//		return err
//	}
//
//	if err = us.ur.UpdatePassword(&user); err != nil {
//		return err
//	}
//
//	// delete reset tokens
//	// delete session
//
//	return nil
//}
