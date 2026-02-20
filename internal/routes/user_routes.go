package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/padapook/bestbit-core/internal/account/controller"
	"github.com/padapook/bestbit-core/internal/account/repository"
	"github.com/padapook/bestbit-core/internal/account/service"
	"github.com/padapook/bestbit-core/internal/middleware"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, db)
	userCtrl := controller.NewUserController(userSvc)

	publicUserRoutes := router.Group("")
	{
		publicUserRoutes.POST("/user/register", userCtrl.Register)
		publicUserRoutes.POST("/login", userCtrl.Login)
		publicUserRoutes.POST("/login/share-token", userCtrl.LoginByShareToken)
	}

	userRoutes := router.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.POST("/logout", userCtrl.Logout)
		userRoutes.POST("/share-token", userCtrl.GenerateShareToken)
		userRoutes.GET("/:username", userCtrl.GetProfile)
	}
}
