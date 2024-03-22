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
	routerGroup.GET("/status", CheckStatus)

	// AUTH ROUTE
	auth := routerGroup.Group("/auth")
	auth.POST("/register", r.Register)
	auth.PATCH("/register", r.VerifyAfterRegister)
	auth.PATCH("/register/resend", r.ResendOtp)
	auth.POST("/login", r.Login)
	auth.GET("/my-data", r.middleware.Auth, r.MyData)
	auth.DELETE("/logout", r.middleware.Auth, r.Logout)
	auth.POST("/renew-access-token", r.RenewSession)
	auth.POST("/reset", r.ResetPassword)
	auth.GET("/reset/:token", r.CheckResetToken)
	auth.PATCH("/reset/:token", r.ChangePasswordFromReset)

	// PROFILE ROUTE
	profile := routerGroup.Group("/profile")
	profile.GET("/account_number", r.middleware.Auth, r.AllAccountNumber)
	profile.PATCH("/account_number", r.middleware.Auth, r.UpdateAccountNumber)

	// CHAT ROUTE
	chat := routerGroup.Group("/chat")
	chat.GET("/who/:id", r.middleware.Auth, r.GetName)

	// PRODUCT ROUTE
	product := routerGroup.Group("/product")
	product.GET("/homepage", r.HomePage)
	product.GET("/search", r.middleware.Auth, r.SearchProducts)
	product.GET("/detail/:id", r.middleware.Auth, r.DetailProduct)
	product.POST("", r.middleware.Auth, r.AddProduct)
	product.GET("/:id", r.middleware.Auth, r.DetailProductOwner)
	product.POST("/:id", r.middleware.Auth, r.BuyProduct)
	product.PATCH("/:id", r.middleware.Auth, r.UpdateProduct) // later
	product.DELETE("/:id", r.middleware.Auth, r.DeleteProduct)
	product.POST("/payment/callback", r.CheckPayment)

	// TRANSACTION ROUTE
	transaction := routerGroup.Group("/transaction")
	transaction.GET("/buy-list", r.middleware.Auth, r.FindActiveTransactions)
	transaction.GET("/sell-list", r.middleware.Auth, r.FindActiveProducts)
	transaction.DELETE("/:id", r.middleware.Auth, r.CancelTransaction)
	transaction.PATCH("/:id/cash-on-delivery", r.middleware.Auth, r.AcceptTransaction)
	transaction.DELETE("/:id/cash-on-delivery", r.middleware.Auth, r.RefuseTransaction)

}

func CheckStatus(ctx *gin.Context) {
	ctx.JSON(200, "Server OK!")
}

func (r *Rest) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	err := r.router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error while serving: %v", err)
	}
}
