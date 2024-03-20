package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"intern-bcc-2024/entity"
)

type Interface interface {
	CreateAccessToken(userId uuid.UUID) (string, error)
	CreateRefreshToken(userId uuid.UUID) (string, error)
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
	GetLoginUser(ctx *gin.Context) (entity.User, error)
}

type jsonWebToken struct {
	ASecretKey   string
	AExpiredTime time.Duration
	RSecretKey   string
}

type Claims struct {
	UserId uuid.UUID
	jwt.RegisteredClaims
}

var PKG *jsonWebToken

func Init() {
	aSecretKey := os.Getenv("SECRET_KEY_ACCESS_TOKEN")
	aExpiredTime, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TIME"))
	rSecretKey := os.Getenv("SECRET_KEY_REFRESH_TOKEN")
	if err != nil {
		log.Fatalf("failed set expired time for jwt : %v", err.Error())
	}

	PKG = &jsonWebToken{
		ASecretKey:   aSecretKey,
		AExpiredTime: time.Duration(aExpiredTime) * time.Minute,
		RSecretKey:   rSecretKey,
	}
}

func (j *jsonWebToken) CreateAccessToken(userId uuid.UUID) (string, error) {
	claim := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AExpiredTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(j.ASecretKey))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

func (j *jsonWebToken) CreateRefreshToken(userId uuid.UUID) (string, error) {
	claim := &Claims{
		UserId: userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(j.RSecretKey))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

func (j *jsonWebToken) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	var (
		claims Claims
		userId uuid.UUID
	)

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.ASecretKey), nil
	})
	if err != nil {
		return userId, err
	}

	if !token.Valid {
		return userId, err
	}

	userId = claims.UserId

	return userId, nil
}

func (j *jsonWebToken) GetLoginUser(ctx *gin.Context) (entity.User, error) {
	user, ok := ctx.Get("user")
	if !ok {
		return entity.User{}, errors.New("failed to get user")
	}

	return user.(entity.User), nil
}
