package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/padapook/bestbit-core/internal/wallet/service"
)

type WalletController interface {
	GetWallets(c *gin.Context)
	GetWalletByCurrency(c *gin.Context)
}

type walletController struct {
	walletService service.WalletService
}

func NewWalletController(walletService service.WalletService) WalletController {
	return &walletController{walletService: walletService}
}

func (ctrl *walletController) GetWallets(c *gin.Context) {
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	wallets, err := ctrl.walletService.GetUserWallets(accountID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch wallets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": wallets,
	})
}

func (ctrl *walletController) GetWalletByCurrency(c *gin.Context) {
	accountID, exists := c.Get("account_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currency := c.Param("currency")
	if currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currency is required"})
		return
	}

	wallet, err := ctrl.walletService.GetWalletBalance(accountID.(string), currency)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": wallet,
	})
}
