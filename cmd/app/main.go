package main

import (
	"intern-bcc-2024/internal/handler/rest"
	"intern-bcc-2024/internal/repository"
	"intern-bcc-2024/internal/service"
	"intern-bcc-2024/pkg/bcrypt"
	"intern-bcc-2024/pkg/config"
	"intern-bcc-2024/pkg/database/mysql"
	"intern-bcc-2024/pkg/jwt"
	"intern-bcc-2024/pkg/middleware"
	"intern-bcc-2024/pkg/validation"
)

func main() {
	config.LoadEnv()
	validation.AddValidator()
	jwtAuth := jwt.Init()
	h_bcrypt := bcrypt.Init()

	db := mysql.ConnectDatabase()
	mysql.Migrate(db)

	repo := repository.NewRepository(db)
	srvc := service.NewService(service.InitParam{Repository: repo, Bcrypt: h_bcrypt, JwtAuth: jwtAuth})

	mw := middleware.Init(jwtAuth, srvc)

	r := rest.NewRest(srvc, mw)
	r.MountEndpoint()
	r.Run()
}
