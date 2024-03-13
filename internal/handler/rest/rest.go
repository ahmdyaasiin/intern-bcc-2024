package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"intern-bcc-2024/internal/service"
	"intern-bcc-2024/pkg/middleware"
	"log"
	"os"
)

type Rest struct {
	router     *gin.Engine
	service    *service.Service
	middleware middleware.Interface
}

func NewRest(service *service.Service, middleware middleware.Interface) *Rest {
	return &Rest{
		router:     gin.Default(),
		service:    service,
		middleware: middleware,
	}
}

func (r *Rest) MountEndpoint() {
	r.router.Use(r.middleware.Cors())
	r.router.Use(r.middleware.Timeout())

	routerGroup := r.router.Group("/api/v1")
	routerGroup.GET("/status", func(c *gin.Context) {
		c.JSON(200, "Server OK!")
	})

	auth := routerGroup.Group("/auth")
	auth.POST("/register", r.RegisterAccount)
	auth.PATCH("/register/verify", r.VerifyAccount)
	auth.PATCH("/register/resend", r.ResendOtp)
	auth.POST("/reset", r.ResetPassword)
	auth.GET("/reset/:token", r.CheckToken)
	auth.PATCH("/reset/:token", r.ChangePassword)
	auth.POST("/login", r.LoginAccount)
	auth.POST("/renew-access-token", r.RenewSession)
	auth.DELETE("/logout", r.middleware.AuthenticateUser, r.LogoutAccount)
	auth.GET("/my-data", r.middleware.AuthenticateUser, r.MyData)

	//////
	product := routerGroup.Group("/product")
	product.GET("/", r.HomePage)
	product.GET("/:id", r.middleware.AuthenticateUser)
	product.POST("/:id/buy", r.middleware.AuthenticateUser)
	product.GET("/search", r.middleware.AuthenticateUser)

}

func (r *Rest) Run() {
	addr := os.Getenv("APP_ADDRESS")
	port := os.Getenv("APP_PORT")

	err := r.router.Run(fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
		log.Fatalf("Error while serving: %v", err)
	}
}
