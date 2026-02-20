package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/padapook/bestbit-core/internal/middleware"
	"github.com/padapook/bestbit-core/internal/wallet/controller"
	"github.com/padapook/bestbit-core/internal/wallet/repository"
	"github.com/padapook/bestbit-core/internal/wallet/service"
	"gorm.io/gorm"
)

func RegisterWalletRoutes(router *gin.RouterGroup, db *gorm.DB) {
	walletRepo := repository.NewWalletRepository(db)
	walletSvc := service.NewWalletService(walletRepo)
	walletCtrl := controller.NewWalletController(walletSvc)

	walletRoutes := router.Group("/wallet")
	walletRoutes.Use(middleware.AuthMiddleware())
	{
		walletRoutes.GET("/", walletCtrl.GetWallets)
		walletRoutes.GET("/:currency", walletCtrl.GetWalletByCurrency)
		walletRoutes.POST("/deposit", walletCtrl.Deposit)
		walletRoutes.POST("/withdraw", walletCtrl.Withdraw)
		walletRoutes.POST("/transfer", walletCtrl.Transfer)
	}
}
