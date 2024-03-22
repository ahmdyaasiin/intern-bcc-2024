package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/model"
	"intern-bcc-2024/pkg/bcrypt"
	"intern-bcc-2024/pkg/database/mysql"
	"intern-bcc-2024/pkg/jwt"
	"intern-bcc-2024/pkg/mail"
	"intern-bcc-2024/pkg/response"
	"log"
	"os"
	"strconv"
	"time"
)

type IUserService interface {
	Register(requests model.RequestForRegister) (*model.ResponseForRegister, response.Details)
	VerifyAfterRegister(requests model.RequestForVerify) response.Details
	ChangePasswordFromReset(token string, request model.RequestForChangePassword) response.Details
	Login(requests model.RequestForLogin) (*model.ResponseForLogin, response.Details)
	Find(requests model.ParamForFind) (*entity.User, response.Details)
	UpdateAccountNumber(ctx *gin.Context, requests model.RequestUpdateAccountNumber) response.Details
}

type UserService struct {
	db      *gorm.DB
	ur      repository.IUserRepository
	or      repository.IOtpRepository
	tr      repository.ITokenRepository
	sr      repository.ISessionRepository
	ar      repository.IAccountRepository
	bcrypt  bcrypt.Interface
	jwtAuth jwt.Interface
}

func NewUserService(userRepository repository.IUserRepository, otpRepository repository.IOtpRepository, tokenRepository repository.ITokenRepository, sessionRepository repository.ISessionRepository, accountRepository repository.IAccountRepository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) IUserService {
	return &UserService{
		db:      mysql.Connection,
		ur:      userRepository,
		or:      otpRepository,
		tr:      tokenRepository,
		sr:      sessionRepository,
		ar:      accountRepository,
		bcrypt:  bcrypt,
		jwtAuth: jwtAuth,
	}
}

func (us *UserService) Register(requests model.RequestForRegister) (*model.ResponseForRegister, response.Details) {
	res := new(model.ResponseForRegister)
	user := new(entity.User)

	tx := us.db.Begin()
	defer tx.Rollback()

	respDetails := us.ur.Find(tx, user, model.ParamForFind{
		Email: requests.Email,
	})
	if respDetails.Error == nil {
		log.Println("Email has been registered")

		return res, response.Details{Code: 409, Message: "Email has been registered", Error: errors.New("email registered")}
	}

	hashPassword, err := us.bcrypt.GenerateFromPassword(requests.Password)
	if err != nil {
		log.Println(err)

		return res, response.Details{Code: 500, Message: "Failed to generate password", Error: err}
	}

	accountNumber, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if err != nil {
		log.Println(err)

		return res, response.Details{Code: 500, Message: "Failed to parse account number"}
	}

	user = &entity.User{
		ID:              uuid.New(),
		Name:            requests.Name,
		Email:           requests.Email,
		Password:        hashPassword,
		Address:         requests.Address,
		Latitude:        requests.Latitude,
		Longitude:       requests.Longitude,
		AccountNumber:   "0",
		AccountNumberID: accountNumber,
		StatusAccount:   "inactive",
		UrlPhotoProfile: "default.jpg",
	}

	respDetails = us.ur.Create(tx, user)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return res, respDetails
	}

	code := mail.GenerateSixCode()
	otp := &entity.OtpCode{
		ID:     uuid.New(),
		UserID: user.ID,
		Code:   code,
	}

	respDetails = us.or.Create(tx, otp)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return res, respDetails
	}

	if err = mail.SendEmail(user.Email, "Verification Code", "Your Verification Code: "+code); err != nil {
		log.Println(err)

		return res, response.Details{Code: 500, Message: "Failed to send verification code to user, please register again", Error: err}
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return res, response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	res.ID = user.ID
	return res, response.Details{Code: 201, Message: "Success Register", Error: nil}
}

func (us *UserService) VerifyAfterRegister(requests model.RequestForVerify) response.Details {
	user := new(entity.User)
	otp := new(entity.OtpCode)

	tx := us.db.Begin()
	defer tx.Rollback()

	respDetails := us.or.Find(tx, otp, model.ParamForFind{
		UserID: requests.UserID,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if otp.Code != requests.VerificationCode {
		log.Println("Verification code is wrong")

		return response.Details{Code: 401, Message: "Verification code is wrong", Error: errors.New("verification code is wrong")}
	}

	expiredTime, err := strconv.Atoi(os.Getenv("EXPIRED_OTP"))
	if err != nil {
		log.Println("Failed to convert expired time from .env")

		return response.Details{Code: 500, Message: "Failed to convert expired time from .env", Error: err}
	}

	if time.Now().Local().Add(-1*time.Duration(expiredTime)*time.Minute).UnixMilli() >= otp.CreatedAt {
		log.Println("Verification code is expired")

		return response.Details{Code: 401, Message: "Verification code is expired", Error: errors.New("verification code is expired")}
	}

	respDetails = us.ur.Find(tx, user, model.ParamForFind{
		ID: requests.UserID,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	user.StatusAccount = "active"
	respDetails = us.ur.Update(tx, user)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	respDetails = us.or.Delete(tx, otp)
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if err = tx.Commit().Error; err != nil {
		log.Println("failed to commit transaction")

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success verify account", Error: nil}
}

func (us *UserService) ChangePasswordFromReset(tokenRequest string, request model.RequestForChangePassword) response.Details {
	token := new(entity.ResetToken)
	user := new(entity.User)

	tx := us.db.Begin()
	defer tx.Rollback()

	respDetails := us.tr.Find(tx, token, model.ParamForFind{
		Token: tokenRequest, // token must be unique
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	expiredTime, err := strconv.Atoi(os.Getenv("EXPIRED_OTP"))
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to convert expired time from .env", Error: err}
	}

	if time.Now().Local().Add(-1*time.Duration(expiredTime)*time.Minute).UnixMilli() >= token.CreatedAt {
		log.Println("verification code is expired")

		return response.Details{Code: 401, Message: "Verification code is expired", Error: errors.New("verification code is expired")}
	}

	respDetails = us.ur.Find(tx, user, model.ParamForFind{
		ID: token.UserID,
	})
	hashPassword, err := us.bcrypt.GenerateFromPassword(request.Password)
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to generate password", Error: err}
	}

	user.Password = hashPassword
	if respDetails = us.ur.Update(tx, user); respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if respDetails = us.tr.Delete(tx, token); respDetails.Error != nil {
		return respDetails
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success change password", Error: nil}
}

func (us *UserService) Login(requests model.RequestForLogin) (*model.ResponseForLogin, response.Details) {
	user := new(entity.User)
	session := new(entity.Session)
	res := new(model.ResponseForLogin)

	tx := us.db.Begin()
	defer tx.Rollback()

	respDetails := us.ur.Find(tx, user, model.ParamForFind{
		Email: requests.Email,
	})
	if respDetails.Error != nil {
		log.Println(respDetails.Error)

		return res, respDetails
	}

	if user.StatusAccount == "blocked" {
		log.Println("Your account is blocked")

		return res, response.Details{Code: 403, Message: "Your account is blocked", Error: errors.New("blocked")}
	}

	if user.StatusAccount == "inactive" {
		log.Println("Your account is unverified")

		res.UserID = user.ID
		return res, response.Details{Code: 403, Message: "Please verify your account first", Error: errors.New("unverified")}
	}

	err := us.bcrypt.CompareAndHashPassword(user.Password, requests.Password)
	if err != nil {
		log.Println("Wrong password")

		return res, response.Details{Code: 401, Message: "Wrong password", Error: err}
	}

	accessToken, err := us.jwtAuth.CreateAccessToken(user.ID)
	if err != nil {
		log.Println("failed to create access token")

		return res, response.Details{Code: 500, Message: "Failed to create access token", Error: err}
	}

	respDetails = us.sr.Find(tx, session, model.ParamForFind{
		UserID: user.ID,
	})
	if respDetails.Error == nil {
		if respDetails = us.sr.Delete(tx, session); respDetails.Error != nil {
			log.Print("failed to delete session")

			return res, respDetails
		}
	}

	refreshToken, err := us.jwtAuth.CreateRefreshToken(user.ID)
	if err != nil {
		return res, response.Details{Code: 500, Message: "Failed to create refresh token", Error: err}
	}

	expiredAt, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TIME"))
	if err != nil {
		log.Println(err)

		return res, response.Details{Code: 500, Message: "Failed convert refresh token time", Error: err}
	}

	session = &entity.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiredAt: time.Now().Local().Add(time.Duration(expiredAt) * 24 * time.Hour).UnixMilli(),
	}
	if respDetails = us.sr.Create(tx, session); respDetails.Error != nil {
		return res, respDetails
	}

	if err = tx.Commit().Error; err != nil {
		return res, response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	res.AccessToken = accessToken
	res.RefreshToken = refreshToken

	return res, response.Details{Code: 200, Message: "Success login", Error: nil}
}

func (us *UserService) Find(requests model.ParamForFind) (*entity.User, response.Details) {
	user := new(entity.User)

	tx := us.db.Begin()
	defer tx.Rollback()

	respDetails := us.ur.Find(tx, user, requests)
	if respDetails.Error != nil {
		return user, respDetails
	}

	if err := tx.Commit().Error; err != nil {
		return user, response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return user, response.Details{Code: 200, Message: "Success get user", Error: nil}
}

func (us *UserService) UpdateAccountNumber(ctx *gin.Context, requests model.RequestUpdateAccountNumber) response.Details {
	account := new(entity.AccountNumberType)

	tx := us.db.Begin()
	defer tx.Rollback()

	respDetails := us.ar.Find(tx, account, model.ParamForFind{
		ID: requests.ID,
	})
	if respDetails.Error != nil {
		return respDetails
	}

	user, err := us.jwtAuth.GetLoginUser(ctx)
	if err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to get login user", Error: err}
	}

	user.AccountNumber = requests.AccountNumber
	user.AccountNumberID = requests.ID
	if respDetails = us.ur.Update(tx, &user); respDetails.Error != nil {
		log.Println(respDetails.Error)

		return respDetails
	}

	if err = tx.Commit().Error; err != nil {
		log.Println(err)

		return response.Details{Code: 500, Message: "Failed to commit transaction", Error: err}
	}

	return response.Details{Code: 200, Message: "Success update account number", Error: nil}
}
