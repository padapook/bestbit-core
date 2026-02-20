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

	userRoutes := router.Group("")
	{
		userRoutes.POST("/user/register", userCtrl.Register)
		userRoutes.POST("/login", userCtrl.Login)
		userRoutes.POST("/login/share-token", userCtrl.LoginByShareToken)
	}

	protectedUserRoutes := router.Group("/user")
	protectedUserRoutes.Use(middleware.AuthMiddleware())
	{
		protectedUserRoutes.POST("/logout", userCtrl.Logout)
		protectedUserRoutes.POST("/share-token", userCtrl.GenerateShareToken)
		protectedUserRoutes.GET("/:username", userCtrl.GetProfile)
	}
}
