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
	jwt.Init()
	bcrypt.Init()
	mysql.ConnectDatabase()
	mysql.Migrate(mysql.Connection)
	mysql.SeedData(mysql.Connection)

	repo := repository.NewRepository(mysql.Connection)
	srvc := service.NewService(service.InitParam{Repository: repo, Bcrypt: bcrypt.PKG, JwtAuth: jwt.PKG})

	mw := middleware.Init(jwt.PKG, srvc)

	r := rest.NewRest(srvc, mw)
	r.MountEndpoint()
	r.Run()
}
