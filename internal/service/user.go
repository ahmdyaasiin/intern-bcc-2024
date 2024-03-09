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
	"os"
	"strconv"
	"time"
)

type IUserService interface {
	Register(requests model.RequestForRegister) (*model.ResponseRegister, error)
	Verify(requests model.OtpParam) error
	Resend(requests model.RequestForResend) error
	Login(requests model.RequestForLogin) (model.ResponseForLogin, error)
	Renew(requests model.RequestForRenew) (model.ResponseForRenew, error)
	Reset(requests model.RequestForReset) error
	ResetGet(token string) error
	ResetPost(requests model.RequestForChangePassword, token string) error
}

type UserService struct {
	ur      repository.IUserRepository
	bcrypt  bcrypt.Interface
	jwtAuth jwt.Interface
}

func NewUserService(userRepository repository.IUserRepository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) IUserService {
	return &UserService{
		ur:      userRepository,
		bcrypt:  bcrypt,
		jwtAuth: jwtAuth,
	}
}

func (us *UserService) Register(requests model.RequestForRegister) (*model.ResponseRegister, error) {
	hashPassword, err := us.bcrypt.GenerateFromPassword(requests.Password)
	if err != nil {
		return nil, err
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

	user, err = us.ur.CreateUser(user)
	if err != nil {
		return nil, err
	}

	code := mail.GenerateVerificationCode()
	if err = mail.SendEmail(user.Email, "Verification Code", "Your Verification Code: "+code); err != nil {
		return nil, err
	}

	otp := &entity.OtpCode{
		ID:               uuid.New(),
		UserID:           user.ID,
		VerificationCode: code,
		ExpiredAt:        time.Now().Local().Add(1 * time.Hour).UnixMilli(),
	}
	err = us.ur.FirstOTP(otp)
	if err != nil {
		return nil, err
	}

	return &model.ResponseRegister{ID: user.ID}, nil
}

func (us *UserService) Verify(requests model.OtpParam) error {
	otp, err := us.ur.GetOTP(requests)
	if err != nil {
		return err
	}

	if otp.ExpiredAt < time.Now().Local().UnixMilli() {
		return errors.New("otp expired")
	}

	err = us.ur.VerifyUser(otp)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Resend(requests model.RequestForResend) error {
	user, err := us.ur.GetUser(model.UserParam{
		ID: requests.ID,
	})
	if err != nil {
		return err
	}

	if user.StatusAccount == "active" {
		return errors.New("account already verified")
	}

	otp, err := us.ur.GetOTP(model.OtpParam{
		UserID: requests.ID,
	})
	if err != nil {
		return err
	}
	otp.VerificationCode = mail.GenerateVerificationCode()
	otp.ExpiredAt = time.Now().Local().Add(1 * time.Hour).UnixMilli()
	if err = mail.SendEmail(user.Email, "Verification Code", "Your Verification Code: "+otp.VerificationCode); err != nil {
		return err
	}

	err = us.ur.UpdateOTP(otp)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Login(requests model.RequestForLogin) (model.ResponseForLogin, error) {
	res := model.ResponseForLogin{}

	user, err := us.ur.GetUser(model.UserParam{
		Email: requests.Email,
	})
	if err != nil {
		return res, err
	}

	if user.StatusAccount == "inactive" {
		return res, errors.New("mohon verifikasi akun anda terlebih dahulu")
	}

	err = us.bcrypt.CompareAndHashPassword(user.Password, requests.Password)
	if err != nil {
		return res, err
	}

	accessToken, err := us.jwtAuth.CreateAccessToken(user.ID)
	if err != nil {
		return res, err
	}

	checkRefreshToken, err := us.ur.GetSession(model.SessionParam{
		UserID: user.ID,
	})
	if err == nil {
		err = us.ur.ClearSession(&checkRefreshToken)
		if err != nil {
			return res, err
		}
	}
	refreshToken, err := us.jwtAuth.CreateRefreshToken(user.ID)
	if err != nil {
		return res, err
	}

	expiredAt, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TIME"))
	rToken := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiredAt: time.Now().Local().Add(time.Duration(expiredAt) * time.Hour).UnixMilli(),
	}
	if err = us.ur.SetSession(rToken); err != nil {
		return res, err
	}

	res.AccessToken = accessToken
	res.RefreshToken = refreshToken

	return res, nil
}

func (us *UserService) Renew(requests model.RequestForRenew) (model.ResponseForRenew, error) {
	res := model.ResponseForRenew{}

	rToken, err := us.ur.GetSession(model.SessionParam{
		Token: requests.RefreshToken,
	})
	if err != nil {
		return res, err
	}

	if rToken.ExpiredAt < time.Now().Local().UnixMilli() {
		return res, errors.New("refresh token kadaluwarsa")
	}

	accessToken, err := us.jwtAuth.CreateAccessToken(rToken.UserID)
	if err != nil {
		return res, err
	}

	res.AccessToken = accessToken

	return res, nil
}

func (us *UserService) Reset(requests model.RequestForReset) error {
	user, err := us.ur.GetUser(model.UserParam{
		Email: requests.Email,
	})
	if err != nil {
		return err
	}

	token := mail.GenerateRandomString(30)
	if err = mail.SendEmail(user.Email, "Link Reset Password", "Your link: "+os.Getenv("LINK_FRONTEND")+token); err != nil {
		return err
	}

	uToken, err := us.ur.GetToken(model.TokenParam{
		UserID: user.ID,
	})
	if err == nil {
		err = us.ur.DeleteToken(&uToken)
		if err != nil {
			return err
		}
	}

	nToken := &entity.ResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: time.Now().Local().Add(time.Hour * 1).UnixMilli(),
	}
	if err = us.ur.SetToken(nToken); err != nil {
		return err
	}

	return nil
}

func (us *UserService) ResetGet(token string) error {
	_, err := us.ur.GetToken(model.TokenParam{
		Token: token,
	})
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) ResetPost(requests model.RequestForChangePassword, token string) error {
	rToken, err := us.ur.GetToken(model.TokenParam{
		Token: token,
	})
	if err != nil {
		return err
	}

	user, err := us.ur.GetUser(model.UserParam{
		ID: rToken.UserID,
	})
	if err != nil {
		return err
	}

	user.Password, err = us.bcrypt.GenerateFromPassword(requests.Password)
	if err != nil {
		return err
	}

	if err = us.ur.UpdatePassword(&user); err != nil {
		return err
	}

	// delete reset tokens
	// delete session

	return nil
}
