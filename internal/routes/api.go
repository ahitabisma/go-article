package routes

import (
	"go-article/internal/handler"
	"go-article/internal/middleware"
	"go-article/internal/repository"
	"go-article/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Dependency injections
	userRepository := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepository)
	authHandler := handler.NewAuthHandler(authService)

	// Auth Routes (Public)
	auth := r.Group("/auth")
	{
		auth.POST("/register", middleware.RateLimitByIP(), authHandler.Register)
		auth.POST("/login", middleware.RateLimitByIP(), authHandler.Login)
	}

	return r
}
