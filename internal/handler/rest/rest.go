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
	auth.POST("/reset", r.ResetPassword)
	auth.GET("/reset/:token", r.CheckResetToken)
	auth.PATCH("/reset/:token", r.ChangePasswordFromReset)
	auth.POST("/login", r.Login)
	auth.POST("/renew-access-token", r.RenewSession)
	auth.DELETE("/logout", r.middleware.Authentication, r.Logout)
	auth.GET("/my-data", r.middleware.Authentication, r.MyData)

	// PRODUCT ROUTE
	product := routerGroup.Group("/product")
	product.GET("/homepage", r.HomePage)
	product.GET("/search", r.middleware.Authentication, r.SearchProducts)
	product.GET("/detail/:id", r.middleware.Authentication, r.DetailProduct)
	product.GET("/:id", r.middleware.Authorization, r.DetailProductOwner) // add middleware authorization
	product.POST("/:id", r.middleware.Authorization, r.BuyProduct)        // add middleware authorization
	product.PATCH("/:id", r.middleware.Authorization, r.UpdateProduct)    // add middleware authorization
	product.DELETE("/:id", r.middleware.Authorization, r.DeleteProduct)   // add middleware authorization
	product.POST("/:id/callback", r.CheckPayment)

	// TRANSACTION ROUTE
	transaction := routerGroup.Group("/transaction")
	transaction.GET("/buy-list", r.middleware.Authentication, r.AllMyTransaction)                // add middleware authentication
	transaction.GET("/sell-list", r.middleware.Authentication, r.AllMyProduct)                   // add middleware authentication
	transaction.DELETE("/:id", r.middleware.Authorization, r.CancelTransaction)                  // add middleware authentication and middleware authorization
	transaction.PATCH("/:id/cash-on-delivery", r.middleware.Authorization, r.RefuseTransaction)  // add middleware authentication and middleware authorization
	transaction.DELETE("/:id/cash-on-delivery", r.middleware.Authorization, r.AcceptTransaction) // add middleware authentication and middleware authorization

}

func CheckStatus(c *gin.Context) {
	c.JSON(200, "Server OK!")
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
