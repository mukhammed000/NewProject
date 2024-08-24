package api

import (
	"api/api/handler"
	"api/api/middleware"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth service
// @version 1.0
// @description Auth service
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewGin(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:8080/swagger/doc.json"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	ca, err := casbin.NewEnforcer("config/model.conf", "config/policy.csv")
	if err != nil {
		log.Fatal("Casbin error creating enforcer: ", err)
	}
	err = ca.LoadPolicy()
	if err != nil {
		log.Fatal("Casbin error loading policy: ", err)
	}

	auth := r.Group("/auth")
	{
		auth.POST("/signup", h.SignUp)
		auth.POST("/login", h.LogIn)
		auth.POST("/enter-email", h.EnterEmail)
	}

	protected := r.Group("/auth")
	protected.Use(middleware.NewAuth(ca))
	{
		protected.POST("/change-password", h.ChangePassword)
		protected.POST("/forget-password", h.ForgetPassword)
		protected.POST("/reset-password", h.ResetPassword)
		protected.POST("/change-email", h.ChangeEmail)
		protected.POST("/verify-email", h.VerifyEmail)
		protected.POST("/validateToken", h.ValidateToken)
		protected.POST("/refreshToken", h.RefreshToken)
	}

	swaggerURL := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, swaggerURL))

	return r
}
