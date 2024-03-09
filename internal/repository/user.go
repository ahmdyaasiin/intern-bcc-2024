package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/model"
	"time"
)

type IUserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	FirstOTP(otp *entity.OtpCode) error
	GetOTP(param model.OtpParam) (*entity.OtpCode, error)
	VerifyUser(otp *entity.OtpCode) error
	UpdateOTP(otp *entity.OtpCode) error
	GetSession(param model.SessionParam) (entity.RefreshToken, error)
	ClearSession(session *entity.RefreshToken) error
	SetSession(rToken *entity.RefreshToken) error
	GetToken(param model.TokenParam) (entity.ResetToken, error)
	DeleteToken(token *entity.ResetToken) error
	SetToken(token *entity.ResetToken) error
	UpdatePassword(user *entity.User) error

	CreateOTP(user *entity.OtpCode) error
	VerifyAccount(id uuid.UUID, verification_code string) (entity.OtpCode, error)
	ChangeStatusUser(id uuid.UUID) error
	ChangeStatusOTP(id uuid.UUID) error
	GetUserById(id uuid.UUID) (entity.User, error)
	GetOTPById(id uuid.UUID) (entity.OtpCode, error)
	GetUser(param model.UserParam) (entity.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) GetUser(param model.UserParam) (entity.User, error) {
	user := entity.User{}
	err := ur.db.Debug().Where(&param).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (ur *UserRepository) GetSession(param model.SessionParam) (entity.RefreshToken, error) {
	token := entity.RefreshToken{}
	err := ur.db.Debug().Where(&param).First(&token).Error
	if err != nil {
		return token, err
	}

	return token, nil
}

func (ur *UserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	if err := ur.db.Debug().Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) FirstOTP(otp *entity.OtpCode) error {
	if err := ur.db.Debug().Create(otp).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetOTP(param model.OtpParam) (*entity.OtpCode, error) {
	var otp entity.OtpCode
	if err := ur.db.Debug().Where(param).First(&otp).Error; err != nil {
		return nil, err
	}

	return &otp, nil
}

func (ur *UserRepository) VerifyUser(otp *entity.OtpCode) error {
	tx := ur.db.Begin()

	var user entity.User
	if err := ur.db.Debug().Where("id = ?", otp.UserID).First(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	user.StatusAccount = "active"
	if err := ur.db.Debug().Where("id = ?", otp.UserID).Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := ur.db.Debug().Delete(&otp).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateOTP(otp *entity.OtpCode) error {
	if err := ur.db.Debug().Save(otp).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdatePassword(user *entity.User) error {
	if err := ur.db.Debug().Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) ClearSession(session *entity.RefreshToken) error {
	if err := ur.db.Debug().Delete(&session).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) SetSession(rToken *entity.RefreshToken) error {
	if err := ur.db.Debug().Create(rToken).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetToken(param model.TokenParam) (entity.ResetToken, error) {
	token := entity.ResetToken{}
	err := ur.db.Debug().Where(&param).First(&token).Error
	if err != nil {
		return token, err
	}

	return token, nil
}

func (ur *UserRepository) SetToken(token *entity.ResetToken) error {
	if err := ur.db.Debug().Create(token).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) DeleteToken(token *entity.ResetToken) error {
	if err := ur.db.Debug().Delete(&token).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetUserById(id uuid.UUID) (entity.User, error) {
	user := entity.User{}
	err := ur.db.Debug().Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (ur *UserRepository) GetOTPById(id uuid.UUID) (entity.OtpCode, error) {
	otp := entity.OtpCode{}
	err := ur.db.Debug().Where("user_id = ?", id).First(&otp).Error
	if err != nil {
		return otp, err
	}

	return otp, nil
}

func (ur *UserRepository) CreateOTP(otp *entity.OtpCode) error {
	if err := ur.db.Debug().Create(otp).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) VerifyAccount(id uuid.UUID, verification_code string) (entity.OtpCode, error) {
	otp := entity.OtpCode{}
	err := ur.db.Debug().Where("user_id = ? AND verification_code = ?", id, verification_code).First(&otp).Error
	if err != nil {
		return otp, err
	}

	return otp, nil
}

func (ur *UserRepository) ChangeStatusUser(id uuid.UUID) error {
	user := entity.User{}
	if err := ur.db.Debug().Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	user.StatusAccount = "active"
	if err := ur.db.Debug().Where("id = ?", id).Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) ChangeStatusOTP(id uuid.UUID) error {
	otp := entity.OtpCode{}
	if err := ur.db.Debug().Where("user_id = ?", id).First(&otp).Error; err != nil {
		return err
	}

	otp.ExpiredAt = time.Now().Local().UnixMilli()
	if err := ur.db.Debug().Where("user_id = ?", id).Save(&otp).Error; err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) Login(login *model.Login) (entity.User, error) {
	user := entity.User{}
	err := ur.db.Debug().Where(&login).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
