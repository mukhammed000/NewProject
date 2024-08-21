package api

import (
	"api/api/handler"
	// "api/api/middleware"
	// "log"

	// "github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Voting service
// @version 1.0
// @description Voting service
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewGin(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080/swagger/doc.json"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ca, err := casbin.NewEnforcer("config/model.conf", "config/policy.csv")
	// if err != nil {
	// 	panic(err)
	// }

	// err = ca.LoadPolicy()
	// if err != nil {
	// 	log.Fatal("casbin error load policy: ", err)
	// 	panic(err)
	// }

	auth := r.Group("/auth")
	// auth.Use(middleware.NewAuth(ca))
	{
		auth.POST("/signup", h.SignUp)
		auth.POST("/login", h.LogIn)
		auth.POST("/change-password", h.ChangePassword)
		auth.POST("/forget-password", h.ForgetPassword)
		auth.POST("/reset-password", h.ResetPassword)
		auth.POST("/change-email", h.ChangeEmail)
		auth.POST("/verify-email", h.VerifyEmail)
		auth.POST("/enter-email", h.EnterEmail)
		auth.POST("/validateToken", h.ValidateToken)
		auth.POST("/refreshToken", h.RefreshToken)
	}

	swaggerURL := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, swaggerURL))

	return r
}
