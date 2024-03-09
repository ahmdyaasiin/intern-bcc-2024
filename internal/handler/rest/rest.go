package rest

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"intern-bcc-2024/internal/service"
	"log"
	"os"
	"time"
)

type Rest struct {
	router  *gin.Engine
	service *service.Service
}

func NewRest(service *service.Service) *Rest {
	return &Rest{
		router:  gin.Default(),
		service: service,
	}
}

func (r *Rest) MountEndpoint() {
	routerGroup := r.router.Group("/api/v1")
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	routerGroup.Use(cors.New(config))
	auth := routerGroup.Group("/auth")
	auth.POST("/register", r.Register)
	auth.PATCH("/register/verify", r.Verify)
	auth.PATCH("/register/resend", r.Resend)
	////
	auth.POST("/login", r.Login)
	//auth.PATCH("", func(context *gin.Context) {})
	//auth.PATCH("/logout", func(context *gin.Context) {})
	auth.POST("/renew-access-token", r.Renew)
	//
	auth.POST("/reset", r.Reset)
	auth.GET("/reset/:token", r.ResetGet)
	auth.PATCH("/reset/:token", r.ResetPost)
}

func (r *Rest) Run() {
	addr := os.Getenv("APP_ADDRESS")
	port := os.Getenv("APP_PORT")

	err := r.router.Run(fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
		log.Fatalf("Error while serving: %v", err)
	}
}
