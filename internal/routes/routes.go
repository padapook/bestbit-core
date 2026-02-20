package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/padapook/bestbit-core/internal/account/controller"
	"github.com/padapook/bestbit-core/internal/account/repository"
	"github.com/padapook/bestbit-core/internal/account/service"
	"gorm.io/gorm"
)

func Routes(r *gin.Engine, db *gorm.DB) {
	v1 := r.Group("/api/v1")
	{
		userRepo := repository.NewUserRepository(db)
		userSvc := service.NewUserService(userRepo, db)
		userCtrl := controller.NewUserController(userSvc)

		userRoutes := v1.Group("/user")
		{
			userRoutes.POST("/register", userCtrl.Register)
			userRoutes.GET("/:username", userCtrl.GetProfile)
		}
	}
}
