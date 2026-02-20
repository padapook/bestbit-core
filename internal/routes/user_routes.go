package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/padapook/bestbit-core/internal/account/controller"
	"github.com/padapook/bestbit-core/internal/account/repository"
	"github.com/padapook/bestbit-core/internal/account/service"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, db)
	userCtrl := controller.NewUserController(userSvc)

	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/register", userCtrl.Register)
		userRoutes.POST("/login", userCtrl.Login)
		userRoutes.POST("/logout", userCtrl.Logout)
		userRoutes.POST("/login/share-token", userCtrl.LoginByShareToken)
		userRoutes.POST("/share-token", userCtrl.GenerateShareToken)
		userRoutes.GET("/:username", userCtrl.GetProfile)
	}
}
